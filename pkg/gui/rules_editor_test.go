package gui

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompileScriptBasic(t *testing.T) {
	src := `
		// Comentário inicial
		MSG 10
		GIVE_ITEM 5
		END
	`
	bytecode, err := CompileScript(src)
	if err != nil {
		t.Fatalf("Erro inesperado ao compilar script básico: %v", err)
	}

	// Esperado:
	// OpMsg (0x01), uint16(10) -> 0x0A, 0x00
	// OpGiveItem (0x02), uint16(5) -> 0x05, 0x00
	// OpEnd (0xFF)
	expected := []byte{0x01, 0x0A, 0x00, 0x02, 0x05, 0x00, 0xFF}
	if !bytes.Equal(bytecode, expected) {
		t.Fatalf("Bytecode incorreto: obtido %v, esperado %v", bytecode, expected)
	}
}

func TestCompileScriptWithLabelsAndBranches(t *testing.T) {
	src := `
		CHECK_FLAG 3 ja_aberto
		MSG 1
		SET_FLAG 3 1
		TELEPORT 2 16 9
		END
	ja_aberto:
		HEAL 20
		DAMAGE 5
		PLAY_SFX 2
		END
	`
	bytecode, err := CompileScript(src)
	if err != nil {
		t.Fatalf("Erro ao compilar script com labels: %v", err)
	}

	if len(bytecode) == 0 {
		t.Fatal("Bytecode compilado não deveria estar vazio")
	}

	// Teste de desassemblagem
	disasm := DisassembleScript(bytecode)
	if !strings.Contains(disasm, "CHECK_FLAG 3") || !strings.Contains(disasm, "HEAL 20") {
		t.Fatalf("Desassemblador não encontrou instruções esperadas: \n%s", disasm)
	}
}

func TestCompileScriptErrors(t *testing.T) {
	// 1. Mnemônico desconhecido
	_, err := CompileScript("INVADIR_SISTEMA 123\nEND")
	if err == nil {
		t.Fatal("Esperado erro para mnemônico desconhecido")
	}

	// 2. Argumentos insuficientes
	_, err = CompileScript("MSG\nEND")
	if err == nil {
		t.Fatal("Esperado erro para MSG sem parâmetros")
	}

	// 3. Coordenadas fora dos limites no TELEPORT
	_, err = CompileScript("TELEPORT 1 99 99\nEND")
	if err == nil {
		t.Fatal("Esperado erro para coordenadas fora da tela 32x18 no TELEPORT")
	}

	// 4. Label não encontrado
	_, err = CompileScript("CHECK_FLAG 1 label_inexistente\nEND")
	if err == nil {
		t.Fatal("Esperado erro para label não declarado")
	}
}

func TestDisassembleAndFormatHex(t *testing.T) {
	raw := []byte{0x01, 0x02, 0x00, 0xFF}
	hexStr := FormatBytecodeHex(raw)
	if hexStr != "01 02 00 FF" {
		t.Fatalf("Formato hexadecimal incorreto: obtido '%s'", hexStr)
	}

	disasm := DisassembleScript(raw)
	expected := "MSG 2\nEND"
	if disasm != expected {
		t.Fatalf("Desassemblagem incorreta: obtido '%s', esperado '%s'", disasm, expected)
	}
}

func TestWrapText(t *testing.T) {
	text := "O heroi entra na masmorra sombria do MSX 2 e encontra uma espada mágica brilhante."
	lines := wrapText(text, 25)

	if len(lines) == 0 {
		t.Fatal("wrapText não gerou linhas")
	}
	for i, l := range lines {
		if len([]rune(l)) > 25 {
			t.Fatalf("Linha %d excedeu limite de 25 caracteres: '%s' (%d)", i, l, len([]rune(l)))
		}
	}
}
