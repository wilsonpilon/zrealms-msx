package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/zrealm-msx/zrealm/pkg/models"
)

func TestFloodFillRoom(t *testing.T) {
	matrix := make([]byte, models.RoomMatrixSize) // 576 bytes todos 0

	// 1. Cria uma barreira fechada de tiles com valor 1 ao redor de uma área 4x4 no canto superior-esquerdo
	// Área interna: x: 1..2, y: 1..2 (4 tiles)
	// Borda: x=0, x=3, y=0, y=3
	for x := 0; x <= 3; x++ {
		matrix[0*models.RoomWidth+x] = 1
		matrix[3*models.RoomWidth+x] = 1
	}
	for y := 0; y <= 3; y++ {
		matrix[y*models.RoomWidth+0] = 1
		matrix[y*models.RoomWidth+3] = 1
	}

	// 2. Preenche o interior (1, 1) com tile 5
	FloodFillRoom(matrix, 1, 1, 5)

	// Verifica se a área interna foi preenchida
	for y := 1; y <= 2; y++ {
		for x := 1; x <= 2; x++ {
			val := matrix[y*models.RoomWidth+x]
			if val != 5 {
				t.Fatalf("Pixel interno (%d, %d) deveria ser 5, mas é %d", x, y, val)
			}
		}
	}

	// Verifica se a barreira com valor 1 permaneceu intacta
	for x := 0; x <= 3; x++ {
		if matrix[0*models.RoomWidth+x] != 1 || matrix[3*models.RoomWidth+x] != 1 {
			t.Fatal("Barreira horizontal foi corrompida pelo FloodFill")
		}
	}

	// Verifica se fora da barreira continua 0
	if matrix[4*models.RoomWidth+4] != 0 {
		t.Fatal("FloodFill vazou para fora da barreira")
	}

	// 3. Testa FloodFill quando o destino já possui a mesma cor (deve retornar imediatamente)
	FloodFillRoom(matrix, 1, 1, 5)
	if matrix[1*models.RoomWidth+1] != 5 {
		t.Fatal("Matriz alterada indevidamente em preenchimento idêntico")
	}
}

func TestFillBorderRoom(t *testing.T) {
	matrix := make([]byte, models.RoomMatrixSize)

	FillBorderRoom(matrix, 99)

	// Verifica linha superior (y=0) e inferior (y=17)
	for x := 0; x < models.RoomWidth; x++ {
		if matrix[0*models.RoomWidth+x] != 99 {
			t.Fatalf("Borda superior falhou em x=%d", x)
		}
		if matrix[(models.RoomHeight-1)*models.RoomWidth+x] != 99 {
			t.Fatalf("Borda inferior falhou em x=%d", x)
		}
	}

	// Verifica coluna esquerda (x=0) e direita (x=31)
	for y := 0; y < models.RoomHeight; y++ {
		if matrix[y*models.RoomWidth+0] != 99 {
			t.Fatalf("Borda esquerda falhou em y=%d", y)
		}
		if matrix[y*models.RoomWidth+(models.RoomWidth-1)] != 99 {
			t.Fatalf("Borda direita falhou em y=%d", y)
		}
	}

	// Verifica centro (1, 1) e (16, 9) que devem continuar 0
	if matrix[1*models.RoomWidth+1] != 0 || matrix[9*models.RoomWidth+16] != 0 {
		t.Fatal("Tiles internos não deveriam ser afetados pelo FillBorderRoom")
	}
}

func TestClearRoomMatrix(t *testing.T) {
	matrix := make([]byte, models.RoomMatrixSize)
	ClearRoomMatrix(matrix, 42)

	for i, b := range matrix {
		if b != 42 {
			t.Fatalf("Byte na posição %d deveria ser 42, encontrado %d", i, b)
		}
	}
}

func TestAutoConnectRooms(t *testing.T) {
	// Cria grade 2x2 de salas:
	// A(0,0) - B(1,0)
	// |        |
	// C(0,1) - D(1,1)
	roomA := models.NewRoom(0, 0, "Sala A", 1)
	roomA.ID = 101

	roomB := models.NewRoom(1, 0, "Sala B", 1)
	roomB.ID = 102

	roomC := models.NewRoom(0, 1, "Sala C", 1)
	roomC.ID = 103

	roomD := models.NewRoom(1, 1, "Sala D", 1)
	roomD.ID = 104

	rooms := []*models.Room{roomA, roomB, roomC, roomD}
	changes := AutoConnectRooms(rooms)

	if changes != 8 { // 4 salas com 2 conexões cada = 8 conexões estabelecidas
		t.Fatalf("Esperado 8 conexões criadas, obtido %d", changes)
	}

	// Validações da Sala A
	if roomA.EastRoomID == nil || *roomA.EastRoomID != 102 {
		t.Fatalf("Sala A deveria conectar ao Leste com Sala B (102)")
	}
	if roomA.SouthRoomID == nil || *roomA.SouthRoomID != 103 {
		t.Fatalf("Sala A deveria conectar ao Sul com Sala C (103)")
	}
	if roomA.NorthRoomID != nil || roomA.WestRoomID != nil {
		t.Fatalf("Sala A não deveria ter conexão Norte ou Oeste")
	}

	// Validações da Sala B
	if roomB.WestRoomID == nil || *roomB.WestRoomID != 101 {
		t.Fatalf("Sala B deveria conectar ao Oeste com Sala A (101)")
	}
	if roomB.SouthRoomID == nil || *roomB.SouthRoomID != 104 {
		t.Fatalf("Sala B deveria conectar ao Sul com Sala D (104)")
	}

	// Validações da Sala C
	if roomC.NorthRoomID == nil || *roomC.NorthRoomID != 101 {
		t.Fatalf("Sala C deveria conectar ao Norte com Sala A (101)")
	}
	if roomC.EastRoomID == nil || *roomC.EastRoomID != 104 {
		t.Fatalf("Sala C deveria conectar ao Leste com Sala D (104)")
	}

	// Validações da Sala D
	if roomD.NorthRoomID == nil || *roomD.NorthRoomID != 102 {
		t.Fatalf("Sala D deveria conectar ao Norte com Sala B (102)")
	}
	if roomD.WestRoomID == nil || *roomD.WestRoomID != 103 {
		t.Fatalf("Sala D deveria conectar ao Oeste com Sala C (103)")
	}

	// Testa remoção da Sala B: A.East e D.North devem se desconectar
	roomsWithoutB := []*models.Room{roomA, roomC, roomD}
	changes = AutoConnectRooms(roomsWithoutB)
	if changes != 2 {
		t.Fatalf("Esperado 2 conexões desfeitas com a remoção da sala B, obtido %d", changes)
	}
	if roomA.EastRoomID != nil {
		t.Fatal("Sala A deveria ter desconectado do Leste após remoção de B")
	}
	if roomD.NorthRoomID != nil {
		t.Fatal("Sala D deveria ter desconectado do Norte após remoção de B")
	}
}

func TestRoomCanvasPainting(t *testing.T) {
	test.NewApp()

	room := models.NewRoom(0, 0, "Test Room", 1)
	room.ID = 1

	entities := []*models.Entity{
		{
			ID:           1,
			RoomID:       1,
			Name:         "NPC Guard",
			PosX:         5,
			PosY:         5,
			BehaviorType: models.BehaviorStaticNPC,
		},
	}

	var changed bool
	rc := NewRoomCanvas(room, nil, entities)
	rc.OnChanged(func() {
		changed = true
	})

	// 1. Testa Pincel (Carimbo)
	rc.SetTool(ToolBrush)
	rc.SetSelectedTile(7)
	rc.applyTool(10, 5)

	if !changed {
		t.Fatal("Esperado callback onChanged após pintura de tile")
	}
	val, err := room.GetTile(10, 5)
	if err != nil || val != 7 {
		t.Fatalf("Esperado tile 7 na posição (10, 5), obtido %d, erro: %v", val, err)
	}

	// 2. Testa Borracha (Tile 0)
	rc.SetTool(ToolEraser)
	rc.applyTool(10, 5)
	val, _ = room.GetTile(10, 5)
	if val != 0 {
		t.Fatalf("Esperado tile 0 após borracha, obtido %d", val)
	}

	// 3. Testa Conta-Gotas
	var pickedTile byte
	rc.OnTilePicked(func(tileIdx byte) {
		pickedTile = tileIdx
	})
	_ = room.SetTile(20, 10, 15)
	rc.SetTool(ToolEyedropper)
	rc.applyTool(20, 10)

	if pickedTile != 15 {
		t.Fatalf("Conta-Gotas deveria ter capturado tile 15, obtido %d", pickedTile)
	}
	if rc.GetSelectedTile() != 15 {
		t.Fatalf("Tile selecionado no canvas deveria ser 15, obtido %d", rc.GetSelectedTile())
	}

	// 4. Testa Clear e FillBorder
	rc.Clear(3)
	val, _ = room.GetTile(0, 0)
	if val != 3 {
		t.Fatalf("Clear falhou: esperado tile 3, obtido %d", val)
	}

	rc.FillBorder(9)
	valTop, _ := room.GetTile(0, 0)
	valInner, _ := room.GetTile(1, 1)
	if valTop != 9 || valInner != 3 {
		t.Fatalf("FillBorder falhou: borda=%d (esperado 9), centro=%d (esperado 3)", valTop, valInner)
	}
}

func TestRenderTileThumbnail(t *testing.T) {
	tile := &models.Tile{
		TileIndex:    1,
		PatternBytes: []byte{0xFF, 0x81, 0x81, 0x81, 0x81, 0x81, 0x81, 0xFF},
		ColorBytes:   []byte{0x71, 0x71, 0x71, 0x71, 0x71, 0x71, 0x71, 0x71},
	}

	img := RenderTileThumbnail(tile)
	if img == nil {
		t.Fatal("RenderTileThumbnail retornou imagem nula")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 8 || bounds.Dy() != 8 {
		t.Fatalf("Dimensões do thumbnail devem ser 8x8, obtido %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Thumbnail para tile nulo (não deve panicar)
	nilImg := RenderTileThumbnail(nil)
	if nilImg == nil {
		t.Fatal("RenderTileThumbnail com nil retornou imagem nula")
	}
}
