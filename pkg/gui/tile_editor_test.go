package gui

import (
	"bytes"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/zrealm-msx/zrealm/pkg/models"
)

func TestMSXPalette(t *testing.T) {
	if len(MSXPalette) != 16 {
		t.Fatalf("Esperado 16 cores na paleta MSX, obtido %d", len(MSXPalette))
	}

	// Testa cor 7 (Ciano MSX)
	cyan := GetMSXColor(7)
	if cyan.R != 0x42 || cyan.G != 0xEB || cyan.B != 0xF5 {
		t.Fatalf("Cor 7 (Ciano) incorreta: %+v", cyan)
	}

	// Testa cor 15 (Branco)
	white := GetMSXColor(15)
	if white.R != 0xFF || white.G != 0xFF || white.B != 0xFF {
		t.Fatalf("Cor 15 (Branco) incorreta: %+v", white)
	}
}

func TestPatternTransformations(t *testing.T) {
	// Padrão de teste: linha 0 totalmente preenchida (0xFF), outras vazias
	orig := []byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// 1. Flip Horizontal
	// Linha 0 com apenas o pixel da esquerda (0x80)
	singlePixel := []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	flippedH := flipHorizontal(singlePixel)
	if flippedH[0] != 0x01 {
		t.Fatalf("Flip horizontal falhou: esperado 0x01, obtido 0x%02X", flippedH[0])
	}
	doubleFlippedH := flipHorizontal(flippedH)
	if !bytes.Equal(doubleFlippedH, singlePixel) {
		t.Fatalf("Duplo flip horizontal deve restaurar o original: %+v != %+v", doubleFlippedH, singlePixel)
	}

	// 2. Flip Vertical
	flippedV := flipVertical(orig)
	if flippedV[7] != 0xFF || flippedV[0] != 0x00 {
		t.Fatalf("Flip vertical falhou: linha 7 deveria ser 0xFF, obtido %+v", flippedV)
	}
	doubleFlippedV := flipVertical(flippedV)
	if !bytes.Equal(doubleFlippedV, orig) {
		t.Fatalf("Duplo flip vertical deve restaurar o original: %+v != %+v", doubleFlippedV, orig)
	}

	// 3. Rotação em 90° (4 rotações = 360° = original)
	rot1 := rotatePattern90(orig)
	// A linha horizontal 0 deve ter se tornado uma coluna vertical à direita (0x01 em cada linha)
	for y := 0; y < 8; y++ {
		if rot1[y] != 0x01 {
			t.Fatalf("Rot 90 falhou na linha %d: esperado 0x01, obtido 0x%02X", y, rot1[y])
		}
	}
	rot2 := rotatePattern90(rot1)
	rot3 := rotatePattern90(rot2)
	rot4 := rotatePattern90(rot3)
	if !bytes.Equal(rot4, orig) {
		t.Fatalf("Quatro rotações de 90° devem restaurar o padrão original: %+v != %+v", rot4, orig)
	}
}

func TestCollisionTypeParsing(t *testing.T) {
	if collisionTypeFromString("Passável (0)") != models.CollisionPassable {
		t.Fatal("Esperado CollisionPassable")
	}
	if collisionTypeFromString("Sólido (1)") != models.CollisionSolid {
		t.Fatal("Esperado CollisionSolid")
	}
	if collisionTypeFromString("Água (2)") != models.CollisionWater {
		t.Fatal("Esperado CollisionWater")
	}
	if collisionTypeFromString("Dano (3)") != models.CollisionDamage {
		t.Fatal("Esperado CollisionDamage")
	}
	if collisionTypeFromString("Gatilho (4)") != models.CollisionTrigger {
		t.Fatal("Esperado CollisionTrigger")
	}
}

func TestTileCanvasPixelPainting(t *testing.T) {
	test.NewApp()
	var changed bool
	tc := NewTileCanvas(nil, nil, func() {
		changed = true
	})

	// Inicialmente todos os pixels são 0
	pat := tc.PatternBytes()
	for _, b := range pat {
		if b != 0 {
			t.Fatalf("Esperado byte inicial 0, obtido: %02X", b)
		}
	}

	// 1. Alterna pixel (0, 0) -> deve setar bit 7 da linha 0 (0x80)
	tc.applyPixel(0, 0)
	if !changed {
		t.Fatal("Callback onChanged deveria ter sido disparado")
	}
	pat = tc.PatternBytes()
	if pat[0] != 0x80 {
		t.Fatalf("Esperado 0x80 na linha 0, obtido: 0x%02X", pat[0])
	}

	// 2. Alterna novamente (0, 0) -> deve voltar para 0x00
	tc.applyPixel(0, 0)
	pat = tc.PatternBytes()
	if pat[0] != 0x00 {
		t.Fatalf("Esperado 0x00 na linha 0 após desmarcar, obtido: 0x%02X", pat[0])
	}

	// 3. Ferramenta Lápis (Draw) e Borracha (Erase)
	tc.SetTool(PenDraw)
	tc.applyPixel(7, 3) // Linha 3, pixel da direita (bit 0 -> 0x01)
	pat = tc.PatternBytes()
	if pat[3] != 0x01 {
		t.Fatalf("Esperado 0x01 na linha 3, obtido: 0x%02X", pat[3])
	}

	tc.SetTool(PenErase)
	tc.applyPixel(7, 3)
	pat = tc.PatternBytes()
	if pat[3] != 0x00 {
		t.Fatalf("Esperado 0x00 na linha 3 após borracha, obtido: 0x%02X", pat[3])
	}
}
