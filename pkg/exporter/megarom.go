package exporter

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zrealm-msx/zrealm/pkg/project"
)

//go:embed zrealm_base.rom
var defaultBaseROM []byte

const (
	// ROMBankSize representa o tamanho de um banco de 16 KB no padrão ASCII-16
	ROMBankSize = 16 * 1024

	// ROMHeaderOffset é o deslocamento no Banco 0 (ROM) para o cabeçalho mestre (CPU 0x7C00 - 0x4000 = 0x3C00)
	ROMHeaderOffset = 0x3C00

	// ROMHeaderMaxLen é o espaço máximo alocado para MasterHeader + Resource Directory no Banco 0
	ROMHeaderMaxLen = 0x0400 // 1024 bytes (0x3C00 a 0x3FFF)

	// MinMegaROMSize é o menor tamanho padrão de cartucho MegaROM (128 KB = 8 bancos de 16 KB)
	MinMegaROMSize = 128 * 1024
)

// ExportROMOptions define opções para a geração de cartuchos .ROM.
type ExportROMOptions struct {
	OutputDir   string
	ROMFilename string // Nome do arquivo gerado (padrão: "<nome_projeto>.rom")
	BaseROMPath string // Caminho opcional para binário base customizado (se vazio, usa o embedded)
	PadSizeKB   int    // Tamanho final em KB (128, 256, 512, 1024, etc.); se 0, ajusta automaticamente
}

// ROMExportResult contém metadados sobre o cartucho ROM gerado.
type ROMExportResult struct {
	ROMPath        string
	TotalROMBytes  int
	TotalBanks     int
	TotalSegments  int
	ResourceCount  int
	HeaderOffset   int
	MapperType     string
}

// CalculateStandardROMSize calcula o menor tamanho padrão de cartucho MSX (128KB, 256KB, 512KB, etc.) para os dados.
func CalculateStandardROMSize(totalBytes int) int {
	standardSizes := []int{
		128 * 1024,
		256 * 1024,
		512 * 1024,
		1024 * 1024,
		2048 * 1024,
		4096 * 1024,
	}

	for _, sz := range standardSizes {
		if totalBytes <= sz {
			return sz
		}
	}
	// Se ultrapassar 4MB, alinha ao próximo múltiplo de 16 KB
	rem := totalBytes % ROMBankSize
	if rem != 0 {
		return totalBytes + (ROMBankSize - rem)
	}
	return totalBytes
}

// BuildMegaROM monta a imagem binária de um cartucho MegaROM ASCII-16 completo:
// - Banco 0 (16 KB): Engine Core + MasterHeader e Diretório em 0x3C00 (CPU 0x7C00)
// - Bancos 1..N (16 KB cada): Segmentos de dados do GAME.DAT
// - Preenchimento (padding) com 0xFF até o tamanho alvo padrão.
func BuildMegaROM(baseROM []byte, headerBytes []byte, gameData []byte, targetSizeBytes int) ([]byte, error) {
	if len(baseROM) < ROMBankSize {
		return nil, fmt.Errorf("binário base do motor deve ter ao menos %d bytes (tem %d)", ROMBankSize, len(baseROM))
	}

	if len(headerBytes) > ROMHeaderMaxLen {
		return nil, fmt.Errorf("cabeçalho mestre (%d bytes) excede o limite máximo reservado no Banco 0 (%d bytes)", len(headerBytes), ROMHeaderMaxLen)
	}

	// 1. Prepara o Banco 0 (16 KB)
	bank0 := make([]byte, ROMBankSize)
	copy(bank0, baseROM[:ROMBankSize])

	// 2. Injeta o MasterHeader e o Resource Directory no offset 0x3C00 (CPU 0x7C00)
	copy(bank0[ROMHeaderOffset:ROMHeaderOffset+len(headerBytes)], headerBytes)

	// 3. Calcula o tamanho mínimo necessário (Banco 0 + GameData)
	minSize := ROMBankSize + len(gameData)

	finalSize := targetSizeBytes
	if finalSize <= 0 {
		finalSize = CalculateStandardROMSize(minSize)
	} else if finalSize < minSize {
		return nil, fmt.Errorf("tamanho solicitado (%d bytes) é menor que os dados mínimos necessários (%d bytes)", finalSize, minSize)
	}

	// 4. Cria a imagem ROM final preenchida com 0xFF
	rom := make([]byte, finalSize)
	for i := range rom {
		rom[i] = 0xFF
	}

	// Copia Banco 0 (offset 0x0000)
	copy(rom[0:ROMBankSize], bank0)

	// Copia Game Data a partir do Banco 1 (offset 0x4000)
	if len(gameData) > 0 {
		copy(rom[ROMBankSize:ROMBankSize+len(gameData)], gameData)
	}

	return rom, nil
}

// ExportROM compila os recursos do projeto e gera um arquivo de cartucho .ROM (MegaROM ASCII-16).
func ExportROM(p *project.Project, opts ExportROMOptions) (*ROMExportResult, error) {
	if p == nil {
		return nil, fmt.Errorf("projeto não pode ser nulo")
	}

	if opts.OutputDir == "" {
		opts.OutputDir = filepath.Join(filepath.Dir(p.Path()), "build_msx")
	}

	// Executa a exportação lógica para gerar cabeçalho e dados de jogo
	expResult, err := Export(p, ExportOptions{
		OutputDir:       opts.OutputDir,
		ExportBankFiles: false,
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao empacotar recursos para ROM: %w", err)
	}

	// Lê o HEADER.BIN gerado
	headerBytes, err := os.ReadFile(expResult.HeaderPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler header gerado: %w", err)
	}

	// Lê o GAME.DAT gerado
	gameData, err := os.ReadFile(expResult.DataPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler game data gerado: %w", err)
	}

	// Obtém o binário base do motor
	var baseROM []byte
	if opts.BaseROMPath != "" {
		baseROM, err = os.ReadFile(opts.BaseROMPath)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler binário base customizado (%s): %w", opts.BaseROMPath, err)
		}
	} else if len(defaultBaseROM) >= ROMBankSize {
		baseROM = defaultBaseROM
	} else {
		// Fallback: tenta localizar em engine_msx/out/zrealm.rom
		altPath := filepath.Join(filepath.Dir(p.Path()), "engine_msx", "out", "zrealm.rom")
		baseROM, err = os.ReadFile(altPath)
		if err != nil {
			return nil, fmt.Errorf("binário base do motor não encontrado nem embutido nem em %s: %w", altPath, err)
		}
	}

	targetSizeBytes := 0
	if opts.PadSizeKB > 0 {
		targetSizeBytes = opts.PadSizeKB * 1024
	}

	romBytes, err := BuildMegaROM(baseROM, headerBytes, gameData, targetSizeBytes)
	if err != nil {
		return nil, fmt.Errorf("falha ao montar MegaROM: %w", err)
	}

	// Determina o nome do arquivo .ROM
	romFilename := opts.ROMFilename
	if romFilename == "" {
		baseName := filepath.Base(p.Path())
		baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		if baseName == "" || baseName == "." {
			baseName = "game"
		}
		romFilename = baseName + ".rom"
	}
	if !strings.HasSuffix(strings.ToLower(romFilename), ".rom") {
		romFilename += ".rom"
	}

	romPath := filepath.Join(opts.OutputDir, romFilename)
	if err := os.WriteFile(romPath, romBytes, 0644); err != nil {
		return nil, fmt.Errorf("falha ao salvar arquivo ROM (%s): %w", romPath, err)
	}

	return &ROMExportResult{
		ROMPath:       romPath,
		TotalROMBytes: len(romBytes),
		TotalBanks:    len(romBytes) / ROMBankSize,
		TotalSegments: expResult.TotalSegments,
		ResourceCount: expResult.ResourceCount,
		HeaderOffset:  ROMHeaderOffset,
		MapperType:    "ASCII16",
	}, nil
}
