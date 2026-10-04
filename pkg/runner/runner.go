package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zrealm-msx/zrealm/pkg/exporter"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

// RunOptions parametriza o pipeline de exportação, empacotamento DSK e execução no openMSX.
type RunOptions struct {
	OutputDir   string   // Diretório para onde os binários exportados irão (padrão: <pasta_projeto>/build_msx)
	DskPath     string   // Destino do arquivo .dsk (padrão: <OutputDir>/zrealm.dsk)
	EmulatorExe string   // Caminho customizado para o executável openmsx (opcional)
	Machine     string   // Nome da máquina MSX no openMSX (padrão: "Philips_NMS_8250")
	ExtraArgs   []string // Argumentos adicionais repassados ao openMSX
	TclScript   string   // Script TCL opcional a ser passado com -script
	Async       bool     // Se true, não bloqueia o processo (para uso em GUI)
	AutoRun     bool     // Se true, inicia o openMSX após empacotar o DSK (padrão: true)
	DskToolPath string   // Caminho customizado para msxtar.exe
}

// RunResult sintetiza os artefatos gerados e o processo do emulador disparado.
type RunResult struct {
	ExportResult *exporter.ExportResult
	DskPath      string
	DskSize      int64
	EmulatorExe  string
	CommandLine  []string
	Cmd          *exec.Cmd
	Process      *os.Process
}

// FindOpenMSX localiza o executável openmsx no sistema operacional.
func FindOpenMSX(customPath string) (string, error) {
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			abs, _ := filepath.Abs(customPath)
			return abs, nil
		}
		return "", fmt.Errorf("openmsx não encontrado no caminho especificado: %s", customPath)
	}

	// 1. Variáveis de ambiente
	if env := os.Getenv("OPENMSX_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			abs, _ := filepath.Abs(env)
			return abs, nil
		}
	}
	if env := os.Getenv("OPENMSX_PATH"); env != "" {
		if _, err := os.Stat(env); err == nil {
			abs, _ := filepath.Abs(env)
			return abs, nil
		}
	}

	// 2. PATH do sistema
	if path, err := exec.LookPath("openmsx"); err == nil {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs, nil
		}
		return path, nil
	}
	if path, err := exec.LookPath("openmsx.exe"); err == nil {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs, nil
		}
		return path, nil
	}

	// 3. Caminhos convencionais no Windows
	userProfile := os.Getenv("USERPROFILE")
	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")

	candidates := []string{
		filepath.Join(userProfile, "scoop", "shims", "openmsx.exe"),
		filepath.Join(userProfile, "scoop", "apps", "openmsx", "current", "openmsx.exe"),
		filepath.Join(programFiles, "openMSX", "openmsx.exe"),
		filepath.Join(programFilesX86, "openMSX", "openmsx.exe"),
		"C:\\openmsx\\openmsx.exe",
		"C:\\tools\\openmsx\\openmsx.exe",
	}

	for _, cand := range candidates {
		if cand == "" {
			continue
		}
		if _, err := os.Stat(cand); err == nil {
			abs, _ := filepath.Abs(cand)
			return abs, nil
		}
	}

	return "", fmt.Errorf("emulador openMSX não encontrado. Instale via Scoop ('scoop install openmsx') ou configure o caminho nas opções")
}

// BuildOpenMSXArgs monta a lista de parâmetros padrão de linha de comando para o openMSX.
func BuildOpenMSXArgs(dskPath string, opts RunOptions) []string {
	machine := opts.Machine
	if machine == "" {
		machine = "Philips_NMS_8250"
	}

	args := []string{"-machine", machine}

	if len(opts.ExtraArgs) > 0 {
		args = append(args, opts.ExtraArgs...)
	} else {
		// Se for a máquina padrão Philips NMS 8250, adiciona extensões recomendadas (MSX-DOS 2 + 512KB Mapper)
		if machine == "Philips_NMS_8250" {
			args = append(args, "-ext", "msxdos2", "-ext", "ram512k")
		}
	}

	absDsk, err := filepath.Abs(dskPath)
	if err != nil {
		absDsk = dskPath
	}
	args = append(args, "-diska", absDsk)

	if opts.TclScript != "" {
		absScript, err := filepath.Abs(opts.TclScript)
		if err != nil {
			absScript = opts.TclScript
		}
		args = append(args, "-script", absScript)
	}

	return args
}

// LaunchOpenMSX inicializa o emulador openMSX com a imagem DSK conectada no drive A.
func LaunchOpenMSX(dskPath string, opts RunOptions) (*exec.Cmd, error) {
	emuPath, err := FindOpenMSX(opts.EmulatorExe)
	if err != nil {
		return nil, err
	}

	args := BuildOpenMSXArgs(dskPath, opts)
	cmd := exec.Command(emuPath, args...)

	if opts.Async {
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("falha ao iniciar openmsx: %w", err)
		}
	} else {
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("erro na execução do openmsx: %w", err)
		}
	}

	return cmd, nil
}

// OneClickRun coordena de ponta a ponta o pipeline: Exportação SQLite -> Empacotamento DSK -> Boot no openMSX.
func OneClickRun(proj *project.Project, opts RunOptions) (*RunResult, error) {
	if proj == nil {
		return nil, fmt.Errorf("nenhum projeto fornecido para execução")
	}

	// 1. Determina diretório de saída
	outDir := opts.OutputDir
	if outDir == "" {
		outDir = "./build_msx"
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de saída %s: %w", outDir, err)
	}

	// 2. Exportação Binária do Projeto
	expRes, err := exporter.Export(proj, exporter.ExportOptions{
		OutputDir:       outDir,
		ExportBankFiles: true,
	})
	if err != nil {
		return nil, fmt.Errorf("falha na exportação binária do projeto: %w", err)
	}

	// 3. Determina caminho do arquivo .DSK
	dskPath := opts.DskPath
	if dskPath == "" {
		projName, _ := proj.GetSetting("name")
		sanitized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(projName), " ", "_"))
		if sanitized == "" {
			sanitized = "zrealm"
		}
		dskPath = filepath.Join(outDir, sanitized+".dsk")
	}

	// 4. Empacota o Disquete MSX-DOS 2
	createdDsk, err := PackDsk(PackOptions{
		DskPath:       dskPath,
		HeaderBinPath: expRes.HeaderPath,
		GameDatPath:   expRes.DataPath,
		DskToolPath:   opts.DskToolPath,
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao empacotar disquete .DSK: %w", err)
	}

	fi, err := os.Stat(createdDsk)
	var dskSize int64
	if err == nil {
		dskSize = fi.Size()
	}

	res := &RunResult{
		ExportResult: expRes,
		DskPath:      createdDsk,
		DskSize:      dskSize,
	}

	// 5. Execução no openMSX (se solicitado)
	if opts.AutoRun {
		emuPath, err := FindOpenMSX(opts.EmulatorExe)
		if err != nil {
			return res, fmt.Errorf("disquete gerado com sucesso em %s, mas o emulador não pôde ser iniciado: %w", createdDsk, err)
		}

		res.EmulatorExe = emuPath
		res.CommandLine = append([]string{emuPath}, BuildOpenMSXArgs(createdDsk, opts)...)

		cmd, err := LaunchOpenMSX(createdDsk, opts)
		if err != nil {
			return res, err
		}
		res.Cmd = cmd
		if cmd.Process != nil {
			res.Process = cmd.Process
		}
	}

	return res, nil
}

// RunROMOptions parametriza a exportação e execução de cartuchos .ROM no openMSX.
type RunROMOptions struct {
	OutputDir   string   // Diretório para onde o .ROM exportado irá
	ROMFilename string   // Nome do arquivo .ROM (padrão: <nome_projeto>.rom)
	EmulatorExe string   // Caminho customizado para o executável openmsx (opcional)
	Machine     string   // Nome da máquina MSX no openMSX (padrão: "Philips_NMS_8250")
	ExtraArgs   []string // Argumentos adicionais repassados ao openMSX
	TclScript   string   // Script TCL opcional a ser passado com -script
	Async       bool     // Se true, não bloqueia o processo (para uso em GUI)
	AutoRun     bool     // Se true, inicia o openMSX após gerar o .ROM (padrão: true)
	PadSizeKB   int      // Tamanho padrão do cartucho em KB (128, 256, 512, etc.)
	BaseROMPath string   // Caminho customizado para o binário base do motor
}

// RunROMResult sintetiza o cartucho ROM gerado e o processo do emulador disparado.
type RunROMResult struct {
	ExportResult *exporter.ROMExportResult
	ROMPath      string
	ROMSize      int64
	EmulatorExe  string
	CommandLine  []string
	Cmd          *exec.Cmd
	Process      *os.Process
}

// BuildOpenMSXROMArgs monta a lista de parâmetros de linha de comando para rodar o cartucho no openMSX.
func BuildOpenMSXROMArgs(romPath string, opts RunROMOptions) []string {
	machine := opts.Machine
	if machine == "" {
		machine = "Philips_NMS_8250"
	}

	args := []string{"-machine", machine}

	if len(opts.ExtraArgs) > 0 {
		args = append(args, opts.ExtraArgs...)
	}

	absROM, err := filepath.Abs(romPath)
	if err != nil {
		absROM = romPath
	}
	args = append(args, "-cart", absROM)

	if opts.TclScript != "" {
		absScript, err := filepath.Abs(opts.TclScript)
		if err != nil {
			absScript = opts.TclScript
		}
		args = append(args, "-script", absScript)
	}

	return args
}

// LaunchOpenMSXCart inicializa o emulador openMSX com o cartucho ROM inserido no slot 1 (-cart).
func LaunchOpenMSXCart(romPath string, opts RunROMOptions) (*exec.Cmd, error) {
	emuPath, err := FindOpenMSX(opts.EmulatorExe)
	if err != nil {
		return nil, err
	}

	args := BuildOpenMSXROMArgs(romPath, opts)
	cmd := exec.Command(emuPath, args...)

	if opts.Async {
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("falha ao iniciar openmsx com cartucho: %w", err)
		}
	} else {
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("erro na execução do openmsx com cartucho: %w", err)
		}
	}

	return cmd, nil
}

// OneClickRunROM coordena de ponta a ponta o pipeline: Exportação SQLite -> MegaROM ASCII-16 -> Boot instantâneo no openMSX (-cart).
func OneClickRunROM(proj *project.Project, opts RunROMOptions) (*RunROMResult, error) {
	if proj == nil {
		return nil, fmt.Errorf("nenhum projeto fornecido para execução")
	}

	outDir := opts.OutputDir
	if outDir == "" {
		outDir = "./build_msx"
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de saída %s: %w", outDir, err)
	}

	// 1. Exporta e monta a imagem MegaROM
	expRes, err := exporter.ExportROM(proj, exporter.ExportROMOptions{
		OutputDir:   outDir,
		ROMFilename: opts.ROMFilename,
		BaseROMPath: opts.BaseROMPath,
		PadSizeKB:   opts.PadSizeKB,
	})
	if err != nil {
		return nil, fmt.Errorf("falha na geração do cartucho MegaROM: %w", err)
	}

	fi, err := os.Stat(expRes.ROMPath)
	var romSize int64
	if err == nil {
		romSize = fi.Size()
	}

	res := &RunROMResult{
		ExportResult: expRes,
		ROMPath:      expRes.ROMPath,
		ROMSize:      romSize,
	}

	// 2. Executa no openMSX (se solicitado)
	if opts.AutoRun {
		emuPath, err := FindOpenMSX(opts.EmulatorExe)
		if err != nil {
			return res, fmt.Errorf("cartucho gerado em %s, mas o emulador não pôde ser iniciado: %w", expRes.ROMPath, err)
		}

		res.EmulatorExe = emuPath
		res.CommandLine = append([]string{emuPath}, BuildOpenMSXROMArgs(expRes.ROMPath, opts)...)

		cmd, err := LaunchOpenMSXCart(expRes.ROMPath, opts)
		if err != nil {
			return res, err
		}
		res.Cmd = cmd
		if cmd.Process != nil {
			res.Process = cmd.Process
		}
	}

	return res, nil
}

