package gui

import (
	"bytes"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/zrealm-msx/zrealm/pkg/models"
)

func TestSpriteCanvasPixelPainting(t *testing.T) {
	test.NewApp()
	s := models.NewSprite("TestSprite")
	var changed bool
	sc := NewSpriteCanvas(s, func() {
		changed = true
	})

	// 1. Testa pintura no canto superior-esquerdo (0, 0)
	sc.SetTool(PenDraw)
	sc.applyPixel(0, 0)
	if !changed {
		t.Fatal("Esperado callback onChanged disparado")
	}
	on, err := s.GetPixel(0, 0)
	if err != nil || !on {
		t.Fatalf("Pixel (0,0) deveria estar ligado: %v, %v", on, err)
	}

	// 2. Testa pintura no canto inferior-direito (15, 15)
	sc.applyPixel(15, 15)
	on, err = s.GetPixel(15, 15)
	if err != nil || !on {
		t.Fatalf("Pixel (15,15) deveria estar ligado: %v, %v", on, err)
	}

	// 3. Testa ferramenta de Borracha (Erase)
	sc.SetTool(PenErase)
	sc.applyPixel(0, 0)
	on, err = s.GetPixel(0, 0)
	if err != nil || on {
		t.Fatalf("Pixel (0,0) deveria ter sido apagado: %v, %v", on, err)
	}

	// 4. Testa ferramenta Alternar (Toggle)
	sc.SetTool(PenToggle)
	sc.applyPixel(15, 15) // Estava ligado, deve desligar
	on, _ = s.GetPixel(15, 15)
	if on {
		t.Fatal("Pixel (15,15) deveria ter alternado para desligado")
	}
	sc.applyPixel(15, 15) // Estava desligado, deve ligar
	on, _ = s.GetPixel(15, 15)
	if !on {
		t.Fatal("Pixel (15,15) deveria ter alternado para ligado")
	}
}

func TestSpriteTransformations(t *testing.T) {
	s := models.NewSprite("Hero")
	// Marca pixel assimétrico em (2, 4)
	_ = s.SetPixel(2, 4, true)
	s.ColorBytes[0] = 7  // Ciano
	s.ColorBytes[15] = 2 // Verde

	origPattern := make([]byte, len(s.PatternBytes))
	copy(origPattern, s.PatternBytes)

	// 1. Flip Horizontal
	FlipSpriteHorizontal(s)
	onLeft, _ := s.GetPixel(2, 4)
	onRight, _ := s.GetPixel(13, 4) // 15 - 2 = 13
	if onLeft || !onRight {
		t.Fatalf("Flip horizontal incorreto: Left=%v, Right=%v (esperado Left=false, Right=true)", onLeft, onRight)
	}
	FlipSpriteHorizontal(s)
	if !bytes.Equal(s.PatternBytes, origPattern) {
		t.Fatal("Duplo flip horizontal deve restaurar o padrão original")
	}

	// 2. Flip Vertical
	FlipSpriteVertical(s)
	onTop, _ := s.GetPixel(2, 4)
	onBottom, _ := s.GetPixel(2, 11) // 15 - 4 = 11
	if onTop || !onBottom {
		t.Fatalf("Flip vertical incorreto: Top=%v, Bottom=%v", onTop, onBottom)
	}
	if s.ColorBytes[0] != 2 || s.ColorBytes[15] != 7 {
		t.Fatalf("Cores das scanlines não foram invertidas corretamente: L0=%d, L15=%d", s.ColorBytes[0], s.ColorBytes[15])
	}
	FlipSpriteVertical(s)
	if !bytes.Equal(s.PatternBytes, origPattern) {
		t.Fatal("Duplo flip vertical deve restaurar o padrão original")
	}

	// 3. Rotação em 90° (4x = original)
	RotateSprite90(s)
	RotateSprite90(s)
	RotateSprite90(s)
	RotateSprite90(s)
	if !bytes.Equal(s.PatternBytes, origPattern) {
		t.Fatal("4 rotações de 90° devem restaurar o sprite original")
	}

	// 4. Shift Direcional
	testShift := models.NewSprite("Shift")
	_ = testShift.SetPixel(5, 5, true)
	ShiftSprite(testShift, 2, 3) // Deve ir para (7, 8)
	oldP, _ := testShift.GetPixel(5, 5)
	newP, _ := testShift.GetPixel(7, 8)
	if oldP || !newP {
		t.Fatalf("Shift falhou: (5,5)=%v, (7,8)=%v", oldP, newP)
	}

	// 5. Clear e Fill
	ClearSprite(s)
	for _, b := range s.PatternBytes {
		if b != 0x00 {
			t.Fatalf("Clear falhou, encontrado byte: 0x%02X", b)
		}
	}

	FillSprite(s)
	for _, b := range s.PatternBytes {
		if b != 0xFF {
			t.Fatalf("Fill falhou, encontrado byte: 0x%02X", b)
		}
	}
}
