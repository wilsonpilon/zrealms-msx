package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project/migrations"
	_ "modernc.org/sqlite"
)

func setupTestStorage(t *testing.T) (*Storage, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "zrealm_storage_test_*")
	if err != nil {
		t.Fatalf("falha ao criar temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.rpgproj")
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("falha ao abrir banco de teste: %v", err)
	}
	db.SetMaxOpenConns(1)

	// Aplica migrações
	allMigrations, err := migrations.LoadMigrations()
	if err != nil {
		_ = db.Close()
		_ = os.RemoveAll(tempDir)
		t.Fatalf("falha ao carregar migrações: %v", err)
	}

	for _, m := range allMigrations {
		if _, err := db.Exec(m.SQL); err != nil {
			_ = db.Close()
			_ = os.RemoveAll(tempDir)
			t.Fatalf("falha ao executar migração %s: %v", m.Name, err)
		}
	}

	st := New(db)
	cleanup := func() {
		_ = db.Close()
		_ = os.RemoveAll(tempDir)
	}

	return st, cleanup
}

func TestTilesetAndTileCRUD(t *testing.T) {
	st, cleanup := setupTestStorage(t)
	defer cleanup()
	ctx := context.Background()

	// 1. Criar novo Tileset
	ts := &models.Tileset{
		Name:        "Catacumbas",
		Description: "Masmorra subterrânea escura",
	}
	if err := st.Tilesets.CreateTileset(ctx, ts); err != nil {
		t.Fatalf("falha ao criar tileset: %v", err)
	}
	if ts.ID <= 0 {
		t.Fatalf("ID inválido gerado para tileset: %d", ts.ID)
	}

	// 2. Buscar Tileset criado
	fetched, err := st.Tilesets.GetTileset(ctx, ts.ID)
	if err != nil {
		t.Fatalf("falha ao buscar tileset %d: %v", ts.ID, err)
	}
	if fetched.Name != "Catacumbas" {
		t.Errorf("esperado 'Catacumbas', obtido '%s'", fetched.Name)
	}

	// 3. Salvar Tile (índice 1: parede de pedra sólida)
	tile := models.NewTile(ts.ID, 1)
	tile.CollisionType = models.CollisionSolid
	// Preenche padrão com listras
	for y := 0; y < 8; y++ {
		tile.PatternBytes[y] = 0xAA
		tile.ColorBytes[y] = 0xE1 // Fg: 14 (Cinza), Bg: 1 (Preto)
	}

	if err := st.Tilesets.SaveTile(ctx, tile); err != nil {
		t.Fatalf("falha ao salvar tile: %v", err)
	}

	// 4. Buscar Tile
	savedTile, err := st.Tilesets.GetTile(ctx, ts.ID, 1)
	if err != nil {
		t.Fatalf("falha ao buscar tile: %v", err)
	}
	if savedTile.CollisionType != models.CollisionSolid {
		t.Errorf("esperado colisão sólida, obtido %d", savedTile.CollisionType)
	}
	if savedTile.PatternBytes[0] != 0xAA {
		t.Errorf("esperado byte 0xAA, obtido 0x%X", savedTile.PatternBytes[0])
	}

	// 5. Atualizar Tile (Upsert)
	savedTile.CollisionType = models.CollisionDamage
	if err := st.Tilesets.SaveTile(ctx, savedTile); err != nil {
		t.Fatalf("falha no upsert do tile: %v", err)
	}
	updatedTile, err := st.Tilesets.GetTile(ctx, ts.ID, 1)
	if err != nil {
		t.Fatalf("falha ao recuperar tile atualizado: %v", err)
	}
	if updatedTile.CollisionType != models.CollisionDamage {
		t.Errorf("esperado colisão de dano, obtido %d", updatedTile.CollisionType)
	}
}

func TestSpriteCRUD(t *testing.T) {
	st, cleanup := setupTestStorage(t)
	defer cleanup()
	ctx := context.Background()

	sprite := models.NewSprite("Cavaleiro")
	_ = sprite.SetPixel(8, 8, true)
	_ = sprite.SetLineColor(8, 7) // Ciano

	if err := st.Sprites.CreateSprite(ctx, sprite); err != nil {
		t.Fatalf("falha ao criar sprite: %v", err)
	}

	fetched, err := st.Sprites.GetSprite(ctx, sprite.ID)
	if err != nil {
		t.Fatalf("falha ao buscar sprite: %v", err)
	}
	if fetched.Name != "Cavaleiro" {
		t.Errorf("esperado 'Cavaleiro', obtido '%s'", fetched.Name)
	}

	px, err := fetched.GetPixel(8, 8)
	if err != nil || !px {
		t.Errorf("pixel (8, 8) deveria ser true no sprite recuperado")
	}

	// Update
	fetched.Name = "Paladino"
	if err := st.Sprites.UpdateSprite(ctx, fetched); err != nil {
		t.Fatalf("falha ao atualizar sprite: %v", err)
	}
	reFetched, err := st.Sprites.GetSprite(ctx, sprite.ID)
	if err != nil || reFetched.Name != "Paladino" {
		t.Errorf("esperado nome atualizado 'Paladino'")
	}

	// Delete
	if err := st.Sprites.DeleteSprite(ctx, sprite.ID); err != nil {
		t.Fatalf("falha ao deletar sprite: %v", err)
	}
}

func TestRoomAndEntities(t *testing.T) {
	st, cleanup := setupTestStorage(t)
	defer cleanup()
	ctx := context.Background()

	// Cria tileset
	ts := &models.Tileset{Name: "Floresta"}
	_ = st.Tilesets.CreateTileset(ctx, ts)

	// Cria sala (0, 0)
	room := models.NewRoom(0, 0, "Entrada da Floresta", ts.ID)
	_ = room.SetTile(0, 0, 1) // canto superior esquerdo
	_ = room.SetTile(31, 17, 2) // canto inferior direito

	if err := st.Rooms.CreateRoom(ctx, room); err != nil {
		t.Fatalf("falha ao criar sala: %v", err)
	}

	// Busca por coordenadas
	foundRoom, err := st.Rooms.GetRoomByCoords(ctx, 0, 0)
	if err != nil {
		t.Fatalf("falha ao buscar sala por coordenadas: %v", err)
	}
	if foundRoom.Name != "Entrada da Floresta" {
		t.Errorf("esperado 'Entrada da Floresta', obtido '%s'", foundRoom.Name)
	}

	t1, _ := foundRoom.GetTile(0, 0)
	t2, _ := foundRoom.GetTile(31, 17)
	if t1 != 1 || t2 != 2 {
		t.Errorf("tiles da sala não conferem: t1=%d, t2=%d", t1, t2)
	}

	// Cria entidade na sala
	entity := &models.Entity{
		RoomID:       foundRoom.ID,
		Name:         "Guarda da Entrada",
		PosX:         16,
		PosY:         9,
		BehaviorType: models.BehaviorStaticNPC,
	}
	if err := st.Rooms.CreateEntity(ctx, entity); err != nil {
		t.Fatalf("falha ao criar entidade: %v", err)
	}

	// Lista entidades da sala
	entities, err := st.Rooms.ListEntitiesByRoom(ctx, foundRoom.ID)
	if err != nil {
		t.Fatalf("falha ao listar entidades: %v", err)
	}
	if len(entities) != 1 {
		t.Fatalf("esperada 1 entidade na sala, obtido %d", len(entities))
	}
	if entities[0].Name != "Guarda da Entrada" {
		t.Errorf("esperado 'Guarda da Entrada', obtido '%s'", entities[0].Name)
	}

	// Tenta criar entidade com posição inválida (fora do grid 32x18)
	invalidEntity := &models.Entity{
		RoomID: foundRoom.ID,
		Name:   "Fora da Tela",
		PosX:   32, // limite é 31
		PosY:   0,
	}
	if err := st.Rooms.CreateEntity(ctx, invalidEntity); err == nil {
		t.Error("esperava erro ao tentar criar entidade fora dos limites da sala")
	}
}

func TestGameDataCRUD(t *testing.T) {
	st, cleanup := setupTestStorage(t)
	defer cleanup()
	ctx := context.Background()

	// 1. Script
	script := &models.Script{
		Name:       "OnChestOpen",
		SourceCode: "give_item(GOLDEN_KEY); msg(\"Você encontrou uma chave dourada!\");",
		Bytecode:   []byte{0x01, 0x0A, 0x02, 0x05},
	}
	if err := st.GameData.CreateScript(ctx, script); err != nil {
		t.Fatalf("falha ao criar script: %v", err)
	}
	fetchedScript, err := st.GameData.GetScript(ctx, script.ID)
	if err != nil || fetchedScript.Name != "OnChestOpen" {
		t.Fatalf("falha ao buscar script criado")
	}

	// 2. String
	dialogue := &models.StringEntry{
		ContextTag:  "NPC_GUARD_GREETING",
		TextContent: "Bem-vindo a Valdor, viajante!",
	}
	if err := st.GameData.CreateString(ctx, dialogue); err != nil {
		t.Fatalf("falha ao criar diálogo: %v", err)
	}
	fetchedStr, err := st.GameData.GetString(ctx, dialogue.ID)
	if err != nil || fetchedStr.TextContent != dialogue.TextContent {
		t.Fatalf("falha ao buscar string criada")
	}

	// 3. Hero Class
	heroClass := &models.HeroClass{
		Name:    "Guerreiro",
		BaseHP:  120,
		BaseMP:  20,
		BaseAtk: 15,
		BaseDef: 12,
	}
	if err := st.GameData.CreateHeroClass(ctx, heroClass); err != nil {
		t.Fatalf("falha ao criar classe de herói: %v", err)
	}
	fetchedClass, err := st.GameData.GetHeroClass(ctx, heroClass.ID)
	if err != nil || fetchedClass.BaseHP != 120 {
		t.Fatalf("falha ao buscar classe de herói criada")
	}

	// 4. Item
	item := &models.Item{
		Name:          "Espada de Aço",
		ItemType:      models.ItemWeapon,
		ModifierStat:  1, // Atk
		ModifierValue: 8,
		Price:         150,
	}
	if err := st.GameData.CreateItem(ctx, item); err != nil {
		t.Fatalf("falha ao criar item: %v", err)
	}
	fetchedItem, err := st.GameData.GetItem(ctx, item.ID)
	if err != nil || fetchedItem.Price != 150 {
		t.Fatalf("falha ao buscar item criado")
	}
}
