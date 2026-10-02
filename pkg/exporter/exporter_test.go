package exporter

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

func TestPackerSegmentsAndHeader(t *testing.T) {
	packer := NewPacker()

	// 1. Grava recurso pequeno (100 bytes)
	data1 := make([]byte, 100)
	data1[0] = 0xAA
	entry1, err := packer.WriteResource(ResTypeTileset, 1, data1)
	if err != nil {
		t.Fatalf("falha ao gravar recurso 1: %v", err)
	}
	if entry1.Segment != 0 || entry1.Offset != 0 || entry1.Size != 100 {
		t.Errorf("entry 1 inconsistente: %+v", entry1)
	}

	// 2. Grava recurso grande que ainda cabe no segmento 0 (16000 bytes)
	data2 := make([]byte, 16000)
	entry2, err := packer.WriteResource(ResTypeRoom, 1, data2)
	if err != nil {
		t.Fatalf("falha ao gravar recurso 2: %v", err)
	}
	if entry2.Segment != 0 || entry2.Offset != 100 || entry2.Size != 16000 {
		t.Errorf("entry 2 inconsistente: %+v", entry2)
	}

	// 3. Grava recurso de 500 bytes (deve forçar abertura do segmento 1, pois 16100 + 500 > 16384)
	data3 := make([]byte, 500)
	data3[0] = 0xBB
	entry3, err := packer.WriteResource(ResTypeRoom, 2, data3)
	if err != nil {
		t.Fatalf("falha ao gravar recurso 3: %v", err)
	}
	if entry3.Segment != 1 || entry3.Offset != 0 || entry3.Size != 500 {
		t.Errorf("entry 3 deveria estar no segmento 1 no offset 0, obtido: %+v", entry3)
	}

	if packer.SegmentCount() != 2 {
		t.Errorf("esperado 2 segmentos, obtido %d", packer.SegmentCount())
	}

	// 4. Monta MasterHeader e valida bytes
	headerBytes, err := packer.BuildMasterHeader(1, 16, 9, 1)
	if err != nil {
		t.Fatalf("falha ao montar cabeçalho: %v", err)
	}

	expectedHeaderSize := 32 + (3 * 8) // 32 bytes do MasterHeader + 3 entradas de 8 bytes
	if len(headerBytes) != expectedHeaderSize {
		t.Fatalf("tamanho de HEADER.BIN esperado %d, obtido %d", expectedHeaderSize, len(headerBytes))
	}

	// Valida Magic 'ZR01'
	if string(headerBytes[0:4]) != "ZR01" {
		t.Errorf("magic incorreto: %s", string(headerBytes[0:4]))
	}

	// Valida Versão (uint16 LE)
	ver := binary.LittleEndian.Uint16(headerBytes[4:6])
	if ver != 1 {
		t.Errorf("versão esperada 1, obtido %d", ver)
	}

	// Valida contagem de segmentos
	segCount := binary.LittleEndian.Uint16(headerBytes[6:8])
	if segCount != 2 {
		t.Errorf("segmentos esperados 2, obtido %d", segCount)
	}
}

func TestFullExportPipeline(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zrealm_export_test_*")
	if err != nil {
		t.Fatalf("falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projPath := filepath.Join(tempDir, "game.rpgproj")
	proj, err := project.Create(projPath, "Export Test RPG")
	if err != nil {
		t.Fatalf("falha ao criar projeto: %v", err)
	}
	defer proj.Close()

	ctx := context.Background()
	st := proj.Storage()

	// 1. Cadastra Tileset e Tile
	ts := &models.Tileset{Name: "Caverna", Description: "Rochas escuras"}
	_ = st.Tilesets.CreateTileset(ctx, ts)

	tile1 := models.NewTile(ts.ID, 1)
	tile1.CollisionType = models.CollisionSolid
	for i := 0; i < 8; i++ {
		tile1.PatternBytes[i] = 0xFF
		tile1.ColorBytes[i] = 0xF4
	}
	_ = st.Tilesets.SaveTile(ctx, tile1)

	// 2. Cadastra Sprite de 16x16
	spr := models.NewSprite("Guerreiro")
	_ = spr.SetPixel(0, 0, true)
	_ = st.Sprites.CreateSprite(ctx, spr)

	// 3. Cadastra duas salas conectadas
	room1 := models.NewRoom(0, 0, "Entrada da Caverna", ts.ID)
	_ = room1.SetTile(10, 5, 1)
	_ = st.Rooms.CreateRoom(ctx, room1)

	room2 := models.NewRoom(0, 1, "Profundezas", ts.ID)
	room2.NorthRoomID = &room1.ID
	_ = st.Rooms.CreateRoom(ctx, room2)

	// Conecta o sul da sala 1 à sala 2
	room1.SouthRoomID = &room2.ID
	_ = st.Rooms.UpdateRoom(ctx, room1)

	// 4. Cadastra entidade na sala 1
	entity := &models.Entity{
		RoomID:       room1.ID,
		Name:         "Guardião",
		PosX:         16,
		PosY:         9,
		SpriteID:     &spr.ID,
		BehaviorType: models.BehaviorStaticNPC,
	}
	_ = st.Rooms.CreateEntity(ctx, entity)

	// 5. Cadastra diálogo
	dialogue := &models.StringEntry{
		ContextTag:  "MSG_WELCOME",
		TextContent: "Quem ousa entrar na caverna sombria?",
	}
	_ = st.GameData.CreateString(ctx, dialogue)

	// 6. Executa a Exportação
	outputDir := filepath.Join(tempDir, "msx_bin")
	result, err := Export(proj, ExportOptions{
		OutputDir:       outputDir,
		ExportBankFiles: true,
	})
	if err != nil {
		t.Fatalf("Export falhou com erro: %v", err)
	}

	// 7. Validações dos arquivos gerados
	if result.TotalSegments < 1 {
		t.Errorf("esperado pelo menos 1 segmento, obtido %d", result.TotalSegments)
	}

	// Verifica se HEADER.BIN foi gerado
	headerData, err := os.ReadFile(result.HeaderPath)
	if err != nil {
		t.Fatalf("erro ao ler HEADER.BIN: %v", err)
	}
	if string(headerData[0:4]) != "ZR01" {
		t.Errorf("Magic incorreto em HEADER.BIN: %s", string(headerData[0:4]))
	}

	// Verifica se GAME.DAT tem tamanho múltiplo de 16.384
	gameData, err := os.ReadFile(result.DataPath)
	if err != nil {
		t.Fatalf("erro ao ler GAME.DAT: %v", err)
	}
	if len(gameData)%SegmentSize != 0 {
		t.Errorf("tamanho de GAME.DAT (%d) não é múltiplo exato de 16.384 bytes", len(gameData))
	}

	// Verifica se SEG00.BNK foi gerado com exatamente 16.384 bytes
	seg0File := filepath.Join(outputDir, "SEG00.BNK")
	seg0Data, err := os.ReadFile(seg0File)
	if err != nil {
		t.Fatalf("erro ao ler SEG00.BNK: %v", err)
	}
	if len(seg0Data) != SegmentSize {
		t.Errorf("tamanho de SEG00.BNK esperado %d, obtido %d", SegmentSize, len(seg0Data))
	}
}
