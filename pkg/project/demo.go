package project

import (
	"context"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/script"
)

// CreateDemoProject cria um projeto de demonstração com um tileset completo,
// 2 salas interligadas (Entrada e Câmara dos Pilares), sprite do herói e diálogos.
func CreateDemoProject(filePath string) (*Project, error) {
	proj, err := Create(filePath, "Catacumbas de Cristal")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	st := proj.Storage()

	// 1. Obter ou atualizar o Tileset 1
	tilesets, err := st.Tilesets.ListTilesets(ctx)
	if err != nil || len(tilesets) == 0 {
		return nil, fmt.Errorf("tileset inicial não encontrado: %w", err)
	}
	tilesetID := tilesets[0].ID
	tilesets[0].Name = "Dungeon dos Antigos"
	tilesets[0].Description = "Masmorra de pedra e alvenaria dos reinos perdidos"
	if err := st.Tilesets.UpdateTileset(ctx, tilesets[0]); err != nil {
		return nil, err
	}

	// Tile 0: Chão de masmorra (ponto central sutil)
	// Padrão: 8 bytes, Cor: Fg=9 (Cinza claro), Bg=1 (Preto) -> 0x91
	tile0Pattern := []byte{0x00, 0x00, 0x00, 0x18, 0x18, 0x00, 0x00, 0x00}
	tile0Color := []byte{0x91, 0x91, 0x91, 0x91, 0x91, 0x91, 0x91, 0x91}
	t0 := &models.Tile{
		TilesetID:     tilesetID,
		TileIndex:     0,
		PatternBytes:  tile0Pattern,
		ColorBytes:    tile0Color,
		CollisionType: models.CollisionPassable,
	}
	if err := st.Tilesets.SaveTile(ctx, t0); err != nil {
		return nil, err
	}

	// Tile 1: Parede de tijolos
	// Padrão alternado de tijolos, Cor: Fg=14 (Cinza), Bg=4 (Azul escuro) -> 0xE4
	tile1Pattern := []byte{0xFF, 0x88, 0x88, 0xFF, 0x22, 0x22, 0xFF, 0x88}
	tile1Color := []byte{0xE4, 0xE4, 0xE4, 0xE4, 0xE4, 0xE4, 0xE4, 0xE4}
	t1 := &models.Tile{
		TilesetID:     tilesetID,
		TileIndex:     1,
		PatternBytes:  tile1Pattern,
		ColorBytes:    tile1Color,
		CollisionType: models.CollisionSolid,
	}
	if err := st.Tilesets.SaveTile(ctx, t1); err != nil {
		return nil, err
	}

	// Tile 2: Portal / Passagem Leste
	// Padrão de arco com portal místico, Cor: Fg=10 (Amarelo claro), Bg=1 (Preto) -> 0xA1
	tile2Pattern := []byte{0x7E, 0x81, 0xBD, 0xA5, 0xA5, 0xBD, 0x81, 0xFF}
	tile2Color := []byte{0xA1, 0xA1, 0xA1, 0xA1, 0xA1, 0xA1, 0xA1, 0xA1}
	t2 := &models.Tile{
		TilesetID:     tilesetID,
		TileIndex:     2,
		PatternBytes:  tile2Pattern,
		ColorBytes:    tile2Color,
		CollisionType: models.CollisionTrigger,
	}
	if err := st.Tilesets.SaveTile(ctx, t2); err != nil {
		return nil, err
	}

	// 2. Criação das Matrizes das Salas
	// Sala 1: Entrada das Catacumbas (Perímetro com paredes, saída no Leste)
	room1Matrix := make([]byte, models.RoomMatrixSize)
	for y := 0; y < models.RoomHeight; y++ {
		for x := 0; x < models.RoomWidth; x++ {
			idx := y*models.RoomWidth + x
			if x == 0 || y == 0 || y == models.RoomHeight-1 {
				room1Matrix[idx] = 1 // Parede sólida
			} else if x == models.RoomWidth-1 {
				// Parede Leste com abertura nas linhas 8 e 9
				if y == 8 || y == 9 {
					room1Matrix[idx] = 2 // Portal de saída
				} else {
					room1Matrix[idx] = 1 // Parede
				}
			} else {
				room1Matrix[idx] = 0 // Chão
			}
		}
	}

	// Sala 2: Câmara dos Pilares (Perímetro com abertura no Oeste e 4 pilares centrais)
	room2Matrix := make([]byte, models.RoomMatrixSize)
	for y := 0; y < models.RoomHeight; y++ {
		for x := 0; x < models.RoomWidth; x++ {
			idx := y*models.RoomWidth + x
			if x == models.RoomWidth-1 || y == 0 || y == models.RoomHeight-1 {
				room2Matrix[idx] = 1 // Parede sólida
			} else if x == 0 {
				// Parede Oeste com abertura nas linhas 8 e 9
				if y == 8 || y == 9 {
					room2Matrix[idx] = 0 // Abertura livre
				} else {
					room2Matrix[idx] = 1 // Parede
				}
			} else if (x == 10 || x == 21) && (y == 5 || y == 12) {
				room2Matrix[idx] = 1 // 4 Pilares maciços
			} else if (x == 16 && (y >= 7 && y <= 10)) || (y == 9 && (x >= 14 && x <= 18)) {
				room2Matrix[idx] = 1 // Altar / Cruz central
			} else {
				room2Matrix[idx] = 0 // Chão livre
			}
		}
	}

	// 3. Persistência das Salas no SQLite
	r1 := &models.Room{
		WorldX:     0,
		WorldY:     0,
		Name:       "Entrada das Catacumbas",
		TilesetID:  tilesetID,
		TileMatrix: room1Matrix,
	}
	if err := st.Rooms.CreateRoom(ctx, r1); err != nil {
		return nil, err
	}

	r2 := &models.Room{
		WorldX:     1,
		WorldY:     0,
		Name:       "Camara dos Pilares",
		TilesetID:  tilesetID,
		TileMatrix: room2Matrix,
	}
	if err := st.Rooms.CreateRoom(ctx, r2); err != nil {
		return nil, err
	}

	// Interligação bidirecional (Leste <-> Oeste)
	r1.EastRoomID = &r2.ID
	if err := st.Rooms.UpdateRoom(ctx, r1); err != nil {
		return nil, err
	}

	r2.WestRoomID = &r1.ID
	if err := st.Rooms.UpdateRoom(ctx, r2); err != nil {
		return nil, err
	}

	// 4. Cria Sprite 1: Herói Guerreiro
	sprHero := models.NewSprite("Heroi")
	// Preenche scanlines com cor 15 (Branco)
	for i := range sprHero.ColorBytes {
		sprHero.ColorBytes[i] = 0x0F
	}
	// Desenha uma silhueta simples de herói 16x16
	_ = sprHero.SetPixel(7, 2, true)
	_ = sprHero.SetPixel(8, 2, true)
	_ = sprHero.SetPixel(6, 3, true)
	_ = sprHero.SetPixel(7, 3, true)
	_ = sprHero.SetPixel(8, 3, true)
	_ = sprHero.SetPixel(9, 3, true)
	_ = sprHero.SetPixel(7, 4, true)
	_ = sprHero.SetPixel(8, 4, true)
	_ = sprHero.SetPixel(5, 5, true)
	_ = sprHero.SetPixel(6, 5, true)
	_ = sprHero.SetPixel(7, 5, true)
	_ = sprHero.SetPixel(8, 5, true)
	_ = sprHero.SetPixel(9, 5, true)
	_ = sprHero.SetPixel(10, 5, true)
	_ = sprHero.SetPixel(7, 6, true)
	_ = sprHero.SetPixel(8, 6, true)
	_ = sprHero.SetPixel(6, 7, true)
	_ = sprHero.SetPixel(9, 7, true)
	_ = sprHero.SetPixel(6, 8, true)
	_ = sprHero.SetPixel(9, 8, true)
	if err := st.Sprites.CreateSprite(ctx, sprHero); err != nil {
		return nil, err
	}

	// 5. Cria Sprite 2: Guardião / Sentinela
	sprGuardian := models.NewSprite("Guardiao")
	// Cores Modo 2: topo dourado/amarelo (10), armadura ciano (7), cinto azul (4), botas cinza (14)
	for i := 0; i <= 3; i++ {
		sprGuardian.ColorBytes[i] = 0x0A
	}
	for i := 4; i <= 7; i++ {
		sprGuardian.ColorBytes[i] = 0x07
	}
	for i := 8; i <= 11; i++ {
		sprGuardian.ColorBytes[i] = 0x04
	}
	for i := 12; i <= 15; i++ {
		sprGuardian.ColorBytes[i] = 0x0E
	}
	// Padrão do Guardião com elmo, armadura e lança
	_ = sprGuardian.SetPixel(7, 1, true); _ = sprGuardian.SetPixel(8, 1, true)
	_ = sprGuardian.SetPixel(6, 2, true); _ = sprGuardian.SetPixel(7, 2, true); _ = sprGuardian.SetPixel(8, 2, true); _ = sprGuardian.SetPixel(9, 2, true)
	_ = sprGuardian.SetPixel(6, 3, true); _ = sprGuardian.SetPixel(7, 3, true); _ = sprGuardian.SetPixel(8, 3, true); _ = sprGuardian.SetPixel(9, 3, true)
	for x := 4; x <= 11; x++ { _ = sprGuardian.SetPixel(x, 4, true) }
	_ = sprGuardian.SetPixel(3, 5, true); _ = sprGuardian.SetPixel(6, 5, true); _ = sprGuardian.SetPixel(7, 5, true); _ = sprGuardian.SetPixel(8, 5, true); _ = sprGuardian.SetPixel(9, 5, true); _ = sprGuardian.SetPixel(12, 5, true)
	_ = sprGuardian.SetPixel(3, 6, true); _ = sprGuardian.SetPixel(5, 6, true); _ = sprGuardian.SetPixel(6, 6, true); _ = sprGuardian.SetPixel(7, 6, true); _ = sprGuardian.SetPixel(8, 6, true); _ = sprGuardian.SetPixel(9, 6, true); _ = sprGuardian.SetPixel(10, 6, true); _ = sprGuardian.SetPixel(12, 6, true)
	_ = sprGuardian.SetPixel(3, 7, true); _ = sprGuardian.SetPixel(6, 7, true); _ = sprGuardian.SetPixel(7, 7, true); _ = sprGuardian.SetPixel(8, 7, true); _ = sprGuardian.SetPixel(9, 7, true); _ = sprGuardian.SetPixel(12, 7, true)
	for x := 5; x <= 10; x++ { _ = sprGuardian.SetPixel(x, 8, true) }; _ = sprGuardian.SetPixel(3, 8, true); _ = sprGuardian.SetPixel(12, 8, true)
	_ = sprGuardian.SetPixel(3, 9, true); _ = sprGuardian.SetPixel(6, 9, true); _ = sprGuardian.SetPixel(9, 9, true); _ = sprGuardian.SetPixel(12, 9, true)
	_ = sprGuardian.SetPixel(3, 10, true); _ = sprGuardian.SetPixel(6, 10, true); _ = sprGuardian.SetPixel(9, 10, true); _ = sprGuardian.SetPixel(12, 10, true)
	_ = sprGuardian.SetPixel(3, 11, true); _ = sprGuardian.SetPixel(5, 11, true); _ = sprGuardian.SetPixel(6, 11, true); _ = sprGuardian.SetPixel(9, 11, true); _ = sprGuardian.SetPixel(10, 11, true); _ = sprGuardian.SetPixel(12, 11, true)
	if err := st.Sprites.CreateSprite(ctx, sprGuardian); err != nil {
		return nil, err
	}

	// 6. Cria Sprite 3: Baú de Tesouro
	sprChest := models.NewSprite("Bau")
	// Cores Modo 2: ferragens douradas/amarelas (10) e corpo castanho/madeira (6)
	for i := 0; i <= 7; i++ {
		sprChest.ColorBytes[i] = 0x0A
	}
	for i := 8; i <= 15; i++ {
		sprChest.ColorBytes[i] = 0x06
	}
	// Padrão do baú com tampa arredondada e fechadura
	for x := 4; x <= 11; x++ { _ = sprChest.SetPixel(x, 5, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 6, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 7, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 8, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 9, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 10, true) }
	for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, 11, true) }
	_ = sprChest.SetPixel(3, 12, true); _ = sprChest.SetPixel(4, 12, true); _ = sprChest.SetPixel(11, 12, true); _ = sprChest.SetPixel(12, 12, true)
	if err := st.Sprites.CreateSprite(ctx, sprChest); err != nil {
		return nil, err
	}

	// 7. Itens do Catálogo do RPG
	itemKey := &models.Item{
		Name:     "Chave de Ferro",
		ItemType: models.ItemKey,
		Price:    0,
	}
	if err := st.GameData.CreateItem(ctx, itemKey); err != nil {
		return nil, err
	}

	itemPotion := &models.Item{
		Name:          "Pocao de Vida",
		ItemType:      models.ItemConsumable,
		ModifierStat:  1, // HP
		ModifierValue: 25,
		Price:         10,
	}
	if err := st.GameData.CreateItem(ctx, itemPotion); err != nil {
		return nil, err
	}

	// 8. Tabela de Strings de Diálogo (IDs 1 a 5)
	msgWelcome := &models.StringEntry{
		ContextTag:  "MSG_WELCOME",
		TextContent: "Bem-vindo as Catacumbas de Cristal! Pressione [ESPACO] para interagir.",
	}
	if err := st.GameData.CreateString(ctx, msgWelcome); err != nil {
		return nil, err
	}

	msgGuardianIntro := &models.StringEntry{
		ContextTag:  "MSG_GUARDIAN_1",
		TextContent: "Guardiao: As profundezas sao perigosas. Pegue a Chave de Ferro!",
	}
	if err := st.GameData.CreateString(ctx, msgGuardianIntro); err != nil {
		return nil, err
	}

	msgGuardianRepeat := &models.StringEntry{
		ContextTag:  "MSG_GUARDIAN_2",
		TextContent: "Guardiao: Que a luz guie seus passos nas profundezas.",
	}
	if err := st.GameData.CreateString(ctx, msgGuardianRepeat); err != nil {
		return nil, err
	}

	msgChestOpen := &models.StringEntry{
		ContextTag:  "MSG_CHEST_OPEN",
		TextContent: "Voce abriu o bau e encontrou uma Pocao de Vida! (+25 HP)",
	}
	if err := st.GameData.CreateString(ctx, msgChestOpen); err != nil {
		return nil, err
	}

	msgChestEmpty := &models.StringEntry{
		ContextTag:  "MSG_CHEST_EMPTY",
		TextContent: "O bau esta vazio.",
	}
	if err := st.GameData.CreateString(ctx, msgChestEmpty); err != nil {
		return nil, err
	}

	// 9. Scripts de Eventos da Bytecode VM
	guardianScriptSrc := "CHECK_FLAG 1 ja_falou\nMSG 2\nGIVE_ITEM 1\nSET_FLAG 1 1\nPLAY_SFX 1\nEND\nja_falou:\nMSG 3\nEND"
	guardianBC, err := script.CompileScript(guardianScriptSrc)
	if err != nil {
		return nil, fmt.Errorf("falha ao compilar script do guardiao: %w", err)
	}
	scrGuardian := &models.Script{
		Name:       "Script Guardiao",
		SourceCode: guardianScriptSrc,
		Bytecode:   guardianBC,
	}
	if err := st.GameData.CreateScript(ctx, scrGuardian); err != nil {
		return nil, err
	}

	chestScriptSrc := "CHECK_FLAG 2 bau_aberto\nMSG 4\nGIVE_ITEM 2\nHEAL 25\nSET_FLAG 2 1\nPLAY_SFX 1\nEND\nbau_aberto:\nMSG 5\nEND"
	chestBC, err := script.CompileScript(chestScriptSrc)
	if err != nil {
		return nil, fmt.Errorf("falha ao compilar script do bau: %w", err)
	}
	scrChest := &models.Script{
		Name:       "Script Bau",
		SourceCode: chestScriptSrc,
		Bytecode:   chestBC,
	}
	if err := st.GameData.CreateScript(ctx, scrChest); err != nil {
		return nil, err
	}

	// 10. Entidades da Sala 1: Entrada das Catacumbas
	// - Guardião em (12, 9): NPC Estático vinculado ao Script 1
	guardianEnt1 := &models.Entity{
		RoomID:        r1.ID,
		Name:          "Guardiao",
		PosX:          12,
		PosY:          9,
		SpriteID:      &sprGuardian.ID,
		BehaviorType:  models.BehaviorStaticNPC,
		EventScriptID: &scrGuardian.ID,
	}
	if err := st.Rooms.CreateEntity(ctx, guardianEnt1); err != nil {
		return nil, err
	}

	// - Baú de Ferro em (8, 4): Baú de Tesouro vinculado ao Script 2
	chestEnt1 := &models.Entity{
		RoomID:        r1.ID,
		Name:          "Bau de Ferro",
		PosX:          8,
		PosY:          4,
		SpriteID:      &sprChest.ID,
		BehaviorType:  models.BehaviorChest,
		EventScriptID: &scrChest.ID,
	}
	if err := st.Rooms.CreateEntity(ctx, chestEnt1); err != nil {
		return nil, err
	}

	// 11. Entidades da Sala 2: Câmara dos Pilares
	// - Sentinela em (6, 6): NPC Errante
	sentinelEnt2 := &models.Entity{
		RoomID:       r2.ID,
		Name:         "Sentinela",
		PosX:         6,
		PosY:         6,
		SpriteID:     &sprGuardian.ID,
		BehaviorType: models.BehaviorWanderingNPC,
	}
	if err := st.Rooms.CreateEntity(ctx, sentinelEnt2); err != nil {
		return nil, err
	}

	// - Baú Místico em (21, 9): Baú de Tesouro
	chestEnt2 := &models.Entity{
		RoomID:        r2.ID,
		Name:          "Bau Mistico",
		PosX:          21,
		PosY:          9,
		SpriteID:      &sprChest.ID,
		BehaviorType:  models.BehaviorChest,
		EventScriptID: &scrChest.ID,
	}
	if err := st.Rooms.CreateEntity(ctx, chestEnt2); err != nil {
		return nil, err
	}

	// 12. Atualiza configurações do projeto
	_ = proj.SetSetting("initial_room_id", fmt.Sprintf("%d", r1.ID))
	_ = proj.SetSetting("initial_hero_x", "16")
	_ = proj.SetSetting("initial_hero_y", "9")
	_ = proj.SetSetting("initial_tileset", fmt.Sprintf("%d", tilesetID))

	return proj, nil
}
