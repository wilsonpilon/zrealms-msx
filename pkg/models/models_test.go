package models

import (
	"testing"
)

func TestTilePixelAndColors(t *testing.T) {
	tile := NewTile(1, 0)

	// Valida tamanho padrão
	if len(tile.PatternBytes) != TilePatternSize {
		t.Fatalf("tamanho de pattern esperado %d, obtido %d", TilePatternSize, len(tile.PatternBytes))
	}
	if len(tile.ColorBytes) != TileColorSize {
		t.Fatalf("tamanho de color esperado %d, obtido %d", TileColorSize, len(tile.ColorBytes))
	}

	// Inicialmente todos os pixels devem estar apagados (false)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			pixel, err := tile.GetPixel(x, y)
			if err != nil {
				t.Fatalf("erro ao ler pixel (%d, %d): %v", x, y, err)
			}
			if pixel {
				t.Errorf("pixel (%d, %d) deveria ser false, mas é true", x, y)
			}
		}
	}

	// Acende os pixels na diagonal principal
	for i := 0; i < 8; i++ {
		if err := tile.SetPixel(i, i, true); err != nil {
			t.Fatalf("erro ao setar pixel (%d, %d): %v", i, i, err)
		}
	}

	// Verifica se apenas a diagonal está acesa
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			pixel, err := tile.GetPixel(x, y)
			if err != nil {
				t.Fatalf("erro ao ler pixel (%d, %d): %v", x, y, err)
			}
			expected := (x == y)
			if pixel != expected {
				t.Errorf("pixel (%d, %d) esperado %v, obtido %v", x, y, expected, pixel)
			}
		}
	}

	// Testa cores da linha (Fg: 10 Verde Claro, Bg: 4 Azul Escuro)
	if err := tile.SetColors(2, 10, 4); err != nil {
		t.Fatalf("erro ao definir cores na linha 2: %v", err)
	}
	fg, bg, err := tile.GetColors(2)
	if err != nil {
		t.Fatalf("erro ao ler cores da linha 2: %v", err)
	}
	if fg != 10 || bg != 4 {
		t.Errorf("esperado Fg=10, Bg=4, obtido Fg=%d, Bg=%d", fg, bg)
	}
}

func TestSpriteModo2Layout(t *testing.T) {
	sprite := NewSprite("Hero")

	if len(sprite.PatternBytes) != SpritePatternSize {
		t.Fatalf("esperado %d bytes de pattern para sprite 16x16, obtido %d", SpritePatternSize, len(sprite.PatternBytes))
	}
	if len(sprite.ColorBytes) != SpriteColorSize {
		t.Fatalf("esperado %d bytes de cor para sprite 16x16, obtido %d", SpriteColorSize, len(sprite.ColorBytes))
	}

	// Testa pixel no quadrante superior-esquerdo: (0, 0) -> byte 0, bit 7
	if err := sprite.SetPixel(0, 0, true); err != nil {
		t.Fatalf("falha ao setar pixel (0,0): %v", err)
	}
	if (sprite.PatternBytes[0] & 0x80) == 0 {
		t.Errorf("byte 0 bit 7 deveria estar setado para pixel (0,0)")
	}

	// Testa pixel no quadrante inferior-esquerdo: (7, 15) -> byte 15, bit 0
	if err := sprite.SetPixel(7, 15, true); err != nil {
		t.Fatalf("falha ao setar pixel (7,15): %v", err)
	}
	if (sprite.PatternBytes[15] & 0x01) == 0 {
		t.Errorf("byte 15 bit 0 deveria estar setado para pixel (7,15)")
	}

	// Testa pixel no quadrante superior-direito: (8, 0) -> byte 16, bit 7
	if err := sprite.SetPixel(8, 0, true); err != nil {
		t.Fatalf("falha ao setar pixel (8,0): %v", err)
	}
	if (sprite.PatternBytes[16] & 0x80) == 0 {
		t.Errorf("byte 16 bit 7 deveria estar setado para pixel (8,0)")
	}

	// Testa pixel no quadrante inferior-direito: (15, 15) -> byte 31, bit 0
	if err := sprite.SetPixel(15, 15, true); err != nil {
		t.Fatalf("falha ao setar pixel (15,15): %v", err)
	}
	if (sprite.PatternBytes[31] & 0x01) == 0 {
		t.Errorf("byte 31 bit 0 deveria estar setado para pixel (15,15)")
	}

	// Valida GetPixel nos 4 cantos
	testPoints := [][2]int{{0, 0}, {7, 15}, {8, 0}, {15, 15}}
	for _, p := range testPoints {
		val, err := sprite.GetPixel(p[0], p[1])
		if err != nil {
			t.Fatalf("erro ao ler pixel (%d, %d): %v", p[0], p[1], err)
		}
		if !val {
			t.Errorf("pixel (%d, %d) deveria ser true", p[0], p[1])
		}
	}

	// Testa cor de scanline
	if err := sprite.SetLineColor(5, 8); err != nil { // Vermelho médio = 8
		t.Fatalf("erro ao setar cor do scanline 5: %v", err)
	}
	c, err := sprite.GetLineColor(5)
	if err != nil {
		t.Fatalf("erro ao ler cor do scanline 5: %v", err)
	}
	if c != 8 {
		t.Errorf("esperado cor 8, obtido %d", c)
	}
}

func TestRoomMatrixAccess(t *testing.T) {
	room := NewRoom(0, 0, "Sala do Trono", 1)

	if len(room.TileMatrix) != RoomMatrixSize {
		t.Fatalf("esperado %d bytes na matriz da sala, obtido %d", RoomMatrixSize, len(room.TileMatrix))
	}

	// Define um tile na posição central (16, 9)
	if err := room.SetTile(16, 9, 42); err != nil {
		t.Fatalf("erro ao setar tile (16, 9): %v", err)
	}

	tileVal, err := room.GetTile(16, 9)
	if err != nil {
		t.Fatalf("erro ao obter tile (16, 9): %v", err)
	}
	if tileVal != 42 {
		t.Errorf("esperado tile 42, obtido %d", tileVal)
	}

	// Valida limites fora do grid
	if err := room.SetTile(32, 0, 1); err == nil {
		t.Error("esperava erro ao setar x=32 fora da grade de 32 colunas")
	}
	if err := room.SetTile(0, 18, 1); err == nil {
		t.Error("esperava erro ao setar y=18 fora da grade de 18 linhas")
	}
}
