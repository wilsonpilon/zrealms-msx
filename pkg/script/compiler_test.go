package script_test

import (
	"bytes"
	"testing"

	"github.com/zrealm-msx/zrealm/pkg/script"
)

func TestCompileAllOpcodes(t *testing.T) {
	src := `
; Script de teste completo
NOP
MSG 42
GIVE_ITEM 100
TAKE_ITEM 50
SET_FLAG 5 1
CHECK_FLAG 5 fim
HEAL 25
DAMAGE 10
PLAY_SFX 2
TELEPORT 3 16 9
fim:
END
`
	bc, err := script.CompileScript(src)
	if err != nil {
		t.Fatalf("Erro ao compilar script valido: %v", err)
	}

	if len(bc) == 0 {
		t.Fatalf("Bytecode compilado esta vazio!")
	}

	// Verifica opcodes esperados
	if bc[0] != script.OpNop {
		t.Errorf("Esperado OpNop (0x00) na posicao 0, obtido: 0x%02X", bc[0])
	}

	// MSG 42 -> 0x01, 42, 0
	if bc[1] != script.OpMsg || bc[2] != 42 || bc[3] != 0 {
		t.Errorf("MSG 42 codificado incorretamente: %v", bc[1:4])
	}

	// Último byte deve ser OpEnd (0xFF)
	if bc[len(bc)-1] != script.OpEnd {
		t.Errorf("Esperado OpEnd (0xFF) no fim, obtido: 0x%02X", bc[len(bc)-1])
	}
}

func TestCompileLabelBranching(t *testing.T) {
	src := `
CHECK_FLAG 1 ja_falou
MSG 1
SET_FLAG 1 1
END
ja_falou:
MSG 2
END
`
	bc, err := script.CompileScript(src)
	if err != nil {
		t.Fatalf("Falha na compilacao com labels: %v", err)
	}

	// CHECK_FLAG ocupa 4 bytes: [OpCheckFlag, FlagID, OffsetLow, OffsetHigh]
	if bc[0] != script.OpCheckFlag || bc[1] != 1 {
		t.Errorf("Cabecalho de CHECK_FLAG incorreto: %v", bc[0:2])
	}

	// Offset deve desviar sobre MSG 1 (3 bytes), SET_FLAG 1 1 (3 bytes), END (1 byte) = 7 bytes
	offset := int16(bc[2]) | (int16(bc[3]) << 8)
	if offset != 7 {
		t.Errorf("Offset esperado para label ja_falou = 7, obtido = %d", offset)
	}
}

func TestCompileErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"MnemonicInvalido", "FOOBAR 123"},
		{"FaltandoArgumentos", "SET_FLAG 1"},
		{"LabelInexistente", "CHECK_FLAG 1 label_fantasma\nEND"},
		{"ArgumentoNaoNumerico", "MSG banana"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := script.CompileScript(tt.src)
			if err == nil {
				t.Errorf("Esperava erro ao compilar [%s], mas passou", tt.src)
			}
		})
	}
}

func TestEmptyScript(t *testing.T) {
	bc, err := script.CompileScript("")
	if err != nil {
		t.Fatalf("Script vazio nao deveria retornar erro: %v", err)
	}
	if !bytes.Equal(bc, []byte{script.OpEnd}) {
		t.Errorf("Script vazio deveria gerar apenas OpEnd (0xFF), obtido: %v", bc)
	}
}
