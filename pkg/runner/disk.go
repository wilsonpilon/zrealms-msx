package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PackOptions define as opções de geração e empacotamento da imagem de disquete (.DSK).
type PackOptions struct {
	DskPath         string // Caminho final para o arquivo .dsk
	HeaderBinPath   string // Caminho do HEADER.BIN gerado pelo exporter
	GameDatPath     string // Caminho do GAME.DAT gerado pelo exporter
	DskToolPath     string // Caminho customizado para msxtar.exe (opcional)
	CustomAutoExec  string // Conteúdo opcional para substituir autoexec.bat
	CustomZrealmCom string // Caminho para zrealm.com customizado (opcional)
}

// FindDskTool localiza o utilitário msxtar no sistema.
func FindDskTool(customPath string) (string, error) {
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			abs, _ := filepath.Abs(customPath)
			return abs, nil
		}
		return "", fmt.Errorf("utilitário dsk especificado não encontrado: %s", customPath)
	}

	if env := os.Getenv("MSXTAR_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			abs, _ := filepath.Abs(env)
			return abs, nil
		}
	}

	// 1. Caminho relativo ao repositório local
	candidates := []string{
		"MSXgl/tools/build/msxtar/msxtar.exe",
		"MSXgl/tools/build/msxtar/msxtar",
		"../MSXgl/tools/build/msxtar/msxtar.exe",
		"../../MSXgl/tools/build/msxtar/msxtar.exe",
	}

	// 2. Caminho relativo ao executável atual
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "tools", "msxtar.exe"),
			filepath.Join(exeDir, "msxtar.exe"),
			filepath.Join(exeDir, "..", "MSXgl", "tools", "build", "msxtar", "msxtar.exe"),
		)
	}

	for _, cand := range candidates {
		if _, err := os.Stat(cand); err == nil {
			abs, err := filepath.Abs(cand)
			if err == nil {
				return abs, nil
			}
			return cand, nil
		}
	}

	// 3. Procura no PATH
	if path, err := exec.LookPath("msxtar"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("msxtar.exe"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("utilitário msxtar não encontrado. Certifique-se de que MSXgl/tools/build/msxtar/msxtar.exe está acessível ou adicione ao PATH")
}

// PackDsk empacota o disquete virtual de 720 KB (.DSK) inicializável em MSX-DOS 2 contendo a engine e os dados do jogo.
func PackDsk(opts PackOptions) (string, error) {
	if opts.DskPath == "" {
		return "", fmt.Errorf("caminho de saída do DSK não especificado")
	}
	if opts.HeaderBinPath == "" {
		return "", fmt.Errorf("caminho de HEADER.BIN não especificado")
	}
	if opts.GameDatPath == "" {
		return "", fmt.Errorf("caminho de GAME.DAT não especificado")
	}

	if _, err := os.Stat(opts.HeaderBinPath); err != nil {
		return "", fmt.Errorf("arquivo HEADER.BIN não encontrado: %w", err)
	}
	if _, err := os.Stat(opts.GameDatPath); err != nil {
		return "", fmt.Errorf("arquivo GAME.DAT não encontrado: %w", err)
	}

	dskTool, err := FindDskTool(opts.DskToolPath)
	if err != nil {
		return "", err
	}

	// Cria o diretório de destino se necessário
	dskDir := filepath.Dir(opts.DskPath)
	if err := os.MkdirAll(dskDir, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar pasta do DSK: %w", err)
	}

	absDskPath, err := filepath.Abs(opts.DskPath)
	if err != nil {
		return "", fmt.Errorf("falha ao resolver caminho absoluto do DSK: %w", err)
	}

	// Cria pasta temporária de montagem (staging)
	stagingDir, err := os.MkdirTemp("", "zrealm_dsk_pack_*")
	if err != nil {
		return "", fmt.Errorf("falha ao criar pasta temporária para DSK: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	// Extrai arquivos de sistema padrão
	if err := ExtractSystemFiles(stagingDir); err != nil {
		return "", fmt.Errorf("falha ao extrair arquivos de sistema MSX-DOS 2: %w", err)
	}

	// Sobrescreve autoexec se customizado
	if opts.CustomAutoExec != "" {
		if err := os.WriteFile(filepath.Join(stagingDir, "autoexec.bat"), []byte(opts.CustomAutoExec), 0644); err != nil {
			return "", fmt.Errorf("falha ao salvar autoexec customizado: %w", err)
		}
	}

	// Sobrescreve zrealm.com se customizado
	if opts.CustomZrealmCom != "" {
		zbytes, err := os.ReadFile(opts.CustomZrealmCom)
		if err != nil {
			return "", fmt.Errorf("falha ao ler zrealm.com customizado: %w", err)
		}
		if err := os.WriteFile(filepath.Join(stagingDir, "zrealm.com"), zbytes, 0644); err != nil {
			return "", fmt.Errorf("falha ao copiar zrealm.com customizado: %w", err)
		}
	}

	// Copia HEADER.BIN
	headerData, err := os.ReadFile(opts.HeaderBinPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler HEADER.BIN: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stagingDir, "HEADER.BIN"), headerData, 0644); err != nil {
		return "", fmt.Errorf("falha ao preparar HEADER.BIN no staging: %w", err)
	}

	// Copia GAME.DAT
	gameData, err := os.ReadFile(opts.GameDatPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler GAME.DAT: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stagingDir, "GAME.DAT"), gameData, 0644); err != nil {
		return "", fmt.Errorf("falha ao preparar GAME.DAT no staging: %w", err)
	}

	// Monta o arquivo DSK com msxtar
	// msxtar.exe -cf <dsk> --dos2 --verbose --size=double autoexec.bat zrealm.com COMMAND2.COM MSXDOS2.SYS HEADER.BIN GAME.DAT
	cmdArgs := []string{
		"-cf",
		absDskPath,
		"--dos2",
		"--verbose",
		"--size=double",
		"autoexec.bat",
		"zrealm.com",
		"COMMAND2.COM",
		"MSXDOS2.SYS",
		"HEADER.BIN",
		"GAME.DAT",
	}

	cmd := exec.Command(dskTool, cmdArgs...)
	cmd.Dir = stagingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao executar msxtar (%w):\nSTDOUT: %s\nSTDERR: %s", err, stdout.String(), stderr.String())
	}

	// Validação final de integridade do arquivo .DSK gerado
	fi, err := os.Stat(absDskPath)
	if err != nil {
		return "", fmt.Errorf("arquivo DSK não foi criado em %s: %w", absDskPath, err)
	}

	const expectedSize = 737280 // 720 KB
	if fi.Size() != expectedSize {
		return "", fmt.Errorf("tamanho do DSK inválido: esperado %d bytes, obtido %d bytes", expectedSize, fi.Size())
	}

	return absDskPath, nil
}

// VerifyDskContents executa msxtar -tf para inspecionar os arquivos gravados no DSK.
func VerifyDskContents(dskPath string, customDskTool string) ([]string, error) {
	dskTool, err := FindDskTool(customDskTool)
	if err != nil {
		return nil, err
	}

	absDskPath, err := filepath.Abs(dskPath)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(dskTool, "-tf", absDskPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("erro ao inspecionar DSK: %w (output: %s)", err, string(out))
	}

	lines := strings.Split(string(out), "\n")
	var result []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			result = append(result, l)
		}
	}
	return result, nil
}
