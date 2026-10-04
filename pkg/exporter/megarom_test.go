package exporter

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/zrealm-msx/zrealm/pkg/project"
)

func TestCalculateStandardROMSize(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 128 * 1024},
		{16 * 1024, 128 * 1024},
		{128 * 1024, 128 * 1024},
		{128*1024 + 1, 256 * 1024},
		{256 * 1024, 256 * 1024},
		{256*1024 + 1, 512 * 1024},
		{512 * 1024, 512 * 1024},
		{512*1024 + 1, 1024 * 1024},
		{1024 * 1024, 1024 * 1024},
		{2048 * 1024, 2048 * 1024},
		{4096 * 1024, 4096 * 1024},
		{4096*1024 + 10, 4096*1024 + ROMBankSize},
	}

	for _, tt := range tests {
		got := CalculateStandardROMSize(tt.input)
		if got != tt.expected {
			t.Errorf("CalculateStandardROMSize(%d) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestBuildMegaROM(t *testing.T) {
	baseROM := make([]byte, ROMBankSize)
	// Coloca assinatura MSX no início
	baseROM[0] = 0x41
	baseROM[1] = 0x42
	baseROM[2] = 0x10
	baseROM[3] = 0x40

	// Simula header ZR01 (32 bytes)
	headerBytes := []byte("ZR01\x01\x00\x02\x00\x05\x00\x01\x00\x10\x09\x01")
	headerBytes = append(headerBytes, make([]byte, 32-len(headerBytes))...)

	// Simula 2 segmentos de 16 KB de dados (32 KB total)
	gameData := make([]byte, 32*1024)
	gameData[0] = 0xDE
	gameData[1] = 0xAD
	gameData[16*1024] = 0xBE
	gameData[16*1024+1] = 0xEF

	rom, err := BuildMegaROM(baseROM, headerBytes, gameData, 0)
	if err != nil {
		t.Fatalf("BuildMegaROM falhou: %v", err)
	}

	// 1 banco base + 2 bancos de dados = 48 KB -> deve ajustar para padrão 128 KB
	expectedSize := 128 * 1024
	if len(rom) != expectedSize {
		t.Fatalf("tamanho do ROM = %d; esperado = %d", len(rom), expectedSize)
	}

	// Verifica assinatura MSX no início
	if rom[0] != 0x41 || rom[1] != 0x42 {
		t.Errorf("assinatura MSX inválida: %02X %02X; esperado 41 42", rom[0], rom[1])
	}

	// Verifica injeção do header em ROMHeaderOffset (0x3C00)
	if !bytes.Equal(rom[ROMHeaderOffset:ROMHeaderOffset+4], []byte("ZR01")) {
		t.Errorf("Header não foi injetado em 0x3C00: obtido %q", rom[ROMHeaderOffset:ROMHeaderOffset+4])
	}

	// Verifica Bank 1 (início em 0x4000)
	if rom[0x4000] != 0xDE || rom[0x4001] != 0xAD {
		t.Errorf("Banco 1 incorreto em 0x4000: %02X %02X", rom[0x4000], rom[0x4001])
	}

	// Verifica Bank 2 (início em 0x8000)
	if rom[0x8000] != 0xBE || rom[0x8001] != 0xEF {
		t.Errorf("Banco 2 incorreto em 0x8000: %02X %02X", rom[0x8000], rom[0x8001])
	}

	// Verifica padding 0xFF nos bancos vazios (ex: offset 0xC000)
	for i := 48 * 1024; i < expectedSize; i++ {
		if rom[i] != 0xFF {
			t.Fatalf("byte de padding em %d não é 0xFF (obtido %02X)", i, rom[i])
		}
	}
}

func TestBuildMegaROMErrors(t *testing.T) {
	// Base ROM muito curto
	_, err := BuildMegaROM([]byte{1, 2, 3}, []byte("header"), nil, 0)
	if err == nil {
		t.Error("esperava erro para baseROM menor que 16KB")
	}

	// Header muito grande (> 1024 bytes)
	baseROM := make([]byte, ROMBankSize)
	bigHeader := make([]byte, 1025)
	_, err = BuildMegaROM(baseROM, bigHeader, nil, 0)
	if err == nil {
		t.Error("esperava erro para headerBytes maior que 1024 bytes")
	}

	// Tamanho alvo menor que o necessário
	_, err = BuildMegaROM(baseROM, []byte("ok"), make([]byte, 16*1024), 16*1024)
	if err == nil {
		t.Error("esperava erro para targetSizeBytes insuficiente")
	}
}

func TestExportROMIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zrealm_rom_test_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projPath := filepath.Join(tmpDir, "test.rpgproj")
	p, err := project.CreateDemoProject(projPath)
	if err != nil {
		t.Fatalf("falha ao criar projeto demo: %v", err)
	}
	defer p.Close()

	outDir := filepath.Join(tmpDir, "build_rom")
	res, err := ExportROM(p, ExportROMOptions{
		OutputDir:   outDir,
		ROMFilename: "mygame.rom",
		PadSizeKB:   128,
	})
	if err != nil {
		t.Fatalf("ExportROM falhou: %v", err)
	}

	if res.ROMPath != filepath.Join(outDir, "mygame.rom") {
		t.Errorf("caminho da ROM incorreto: %s", res.ROMPath)
	}
	if res.TotalROMBytes != 128*1024 {
		t.Errorf("tamanho da ROM = %d; esperado 128KB", res.TotalROMBytes)
	}
	if res.TotalBanks != 8 {
		t.Errorf("total de bancos = %d; esperado 8", res.TotalBanks)
	}
	if res.MapperType != "ASCII16" {
		t.Errorf("mapper type = %s; esperado ASCII16", res.MapperType)
	}

	// Verifica o arquivo gravado
	data, err := os.ReadFile(res.ROMPath)
	if err != nil {
		t.Fatalf("falha ao ler ROM gerada: %v", err)
	}
	if len(data) != 128*1024 {
		t.Errorf("tamanho do arquivo no disco = %d; esperado 128KB", len(data))
	}
	if data[0] != 0x41 || data[1] != 0x42 {
		t.Errorf("assinatura MSX inválida no arquivo gerado: %02X %02X", data[0], data[1])
	}
	if !bytes.Equal(data[ROMHeaderOffset:ROMHeaderOffset+4], []byte("ZR01")) {
		t.Errorf("Header ZR01 ausente no offset 0x3C00 do arquivo")
	}
}
