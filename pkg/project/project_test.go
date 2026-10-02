package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestProject(t *testing.T, name string) (*Project, string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "zrealm_test_*")
	if err != nil {
		t.Fatalf("falha ao criar diretório temporário: %v", err)
	}

	projPath := filepath.Join(tempDir, "game.rpgproj")
	proj, err := Create(projPath, name)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("falha ao criar projeto de teste: %v", err)
	}

	cleanup := func() {
		_ = proj.Close()
		_ = os.RemoveAll(tempDir)
	}

	return proj, projPath, cleanup
}

func TestCreateProject(t *testing.T) {
	proj, projPath, cleanup := setupTestProject(t, "A Lenda de Valdor")
	defer cleanup()

	if _, err := os.Stat(projPath); os.IsNotExist(err) {
		t.Fatalf("arquivo do projeto não foi gravado no disco: %s", projPath)
	}

	name, err := proj.GetSetting("name")
	if err != nil {
		t.Fatalf("erro ao obter 'name': %v", err)
	}
	if name != "A Lenda de Valdor" {
		t.Errorf("esperado 'A Lenda de Valdor', obtido '%s'", name)
	}

	platform, err := proj.GetSetting("target_platform")
	if err != nil {
		t.Fatalf("erro ao obter 'target_platform': %v", err)
	}
	if platform != "MSX2_MSXDOS2" {
		t.Errorf("esperado 'MSX2_MSXDOS2', obtido '%s'", platform)
	}

	// Verifica se tileset padrão foi criado
	var tilesetCount int
	err = proj.DB().QueryRow("SELECT COUNT(*) FROM tilesets WHERE name = 'Overworld'").Scan(&tilesetCount)
	if err != nil {
		t.Fatalf("erro ao consultar tilesets: %v", err)
	}
	if tilesetCount != 1 {
		t.Errorf("esperado 1 tileset 'Overworld', obtido %d", tilesetCount)
	}

	// Verifica se tile 0 foi criado
	var tileCount int
	err = proj.DB().QueryRow("SELECT COUNT(*) FROM tiles WHERE tile_index = 0").Scan(&tileCount)
	if err != nil {
		t.Fatalf("erro ao consultar tile 0: %v", err)
	}
	if tileCount != 1 {
		t.Errorf("esperado 1 tile com índice 0, obtido %d", tileCount)
	}
}

func TestCreateExistingFileFails(t *testing.T) {
	_, projPath, cleanup := setupTestProject(t, "Teste Duplicado")
	defer cleanup()

	// Tentativa de criar no mesmo caminho deve falhar
	_, err := Create(projPath, "Outro Jogo")
	if err == nil {
		t.Fatal("esperava erro ao tentar criar arquivo já existente, mas obteve nil")
	}
}

func TestOpenExistingProject(t *testing.T) {
	proj, projPath, cleanup := setupTestProject(t, "Projeto Reaberto")

	// Fecha o projeto atual
	if err := proj.Close(); err != nil {
		cleanup()
		t.Fatalf("falha ao fechar projeto: %v", err)
	}

	// Reabre
	reopened, err := Open(projPath)
	if err != nil {
		cleanup()
		t.Fatalf("falha ao reabrir projeto existente: %v", err)
	}
	defer reopened.Close()
	defer os.RemoveAll(filepath.Dir(projPath))

	name, err := reopened.GetSetting("name")
	if err != nil {
		t.Fatalf("erro ao buscar setting do projeto reaberto: %v", err)
	}
	if name != "Projeto Reaberto" {
		t.Errorf("esperado 'Projeto Reaberto', obtido '%s'", name)
	}

	if err := reopened.ValidateIntegrity(); err != nil {
		t.Fatalf("validação de integridade falhou no projeto reaberto: %v", err)
	}
}

func TestForeignKeyEnforcement(t *testing.T) {
	proj, _, cleanup := setupTestProject(t, "Teste FK")
	defer cleanup()

	// Tentativa de inserir tile com tileset_id inexistente (ex: 9999)
	emptyPattern := make([]byte, 8)
	defaultColor := make([]byte, 8)
	_, err := proj.DB().Exec(`
		INSERT INTO tiles (tileset_id, tile_index, pattern_bytes, color_bytes, collision_type)
		VALUES (9999, 1, ?, ?, 0)
	`, emptyPattern, defaultColor)

	if err == nil {
		t.Fatal("esperava erro de violação de Foreign Key, mas a inserção foi permitida!")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY") {
		t.Errorf("esperava erro contendo 'FOREIGN KEY', recebido: %v", err)
	}
}

func TestCascadeDeleteTileset(t *testing.T) {
	proj, _, cleanup := setupTestProject(t, "Teste Cascade")
	defer cleanup()

	// Obtém o ID do tileset Overworld criado
	var tilesetID int64
	err := proj.DB().QueryRow("SELECT id FROM tilesets LIMIT 1").Scan(&tilesetID)
	if err != nil {
		t.Fatalf("erro ao buscar tileset: %v", err)
	}

	// Remove o tileset
	_, err = proj.DB().Exec("DELETE FROM tilesets WHERE id = ?", tilesetID)
	if err != nil {
		t.Fatalf("erro ao deletar tileset: %v", err)
	}

	// O tile filho deve ter sido deletado em cascata
	var remainingTiles int
	err = proj.DB().QueryRow("SELECT COUNT(*) FROM tiles WHERE tileset_id = ?", tilesetID).Scan(&remainingTiles)
	if err != nil {
		t.Fatalf("erro ao verificar tiles remanescentes: %v", err)
	}
	if remainingTiles != 0 {
		t.Errorf("esperava 0 tiles após cascade delete, encontrados %d", remainingTiles)
	}
}

func TestSettingsCRUD(t *testing.T) {
	proj, _, cleanup := setupTestProject(t, "Teste Settings")
	defer cleanup()

	err := proj.SetSetting("author", "Miyamoto-san")
	if err != nil {
		t.Fatalf("erro ao definir setting: %v", err)
	}

	val, err := proj.GetSetting("author")
	if err != nil {
		t.Fatalf("erro ao buscar setting: %v", err)
	}
	if val != "Miyamoto-san" {
		t.Errorf("esperado 'Miyamoto-san', obtido '%s'", val)
	}

	// Atualização
	err = proj.SetSetting("author", "Hironobu Sakaguchi")
	if err != nil {
		t.Fatalf("erro ao atualizar setting: %v", err)
	}

	val, err = proj.GetSetting("author")
	if err != nil {
		t.Fatalf("erro ao buscar setting atualizado: %v", err)
	}
	if val != "Hironobu Sakaguchi" {
		t.Errorf("esperado 'Hironobu Sakaguchi', obtido '%s'", val)
	}
}

func TestValidateIntegrity(t *testing.T) {
	proj, _, cleanup := setupTestProject(t, "Teste Integridade")
	defer cleanup()

	if err := proj.ValidateIntegrity(); err != nil {
		t.Fatalf("esperava integridade válida, obteve erro: %v", err)
	}
}

func TestProjectStorageIntegration(t *testing.T) {
	proj, _, cleanup := setupTestProject(t, "Teste Storage Integration")
	defer cleanup()

	ctx := t.Context()
	tilesets, err := proj.Storage().Tilesets.ListTilesets(ctx)
	if err != nil {
		t.Fatalf("falha ao listar tilesets via proj.Storage(): %v", err)
	}
	if len(tilesets) != 1 || tilesets[0].Name != "Overworld" {
		t.Fatalf("esperado 1 tileset 'Overworld', obtido: %+v", tilesets)
	}

	tile0, err := proj.Storage().Tilesets.GetTile(ctx, tilesets[0].ID, 0)
	if err != nil {
		t.Fatalf("falha ao buscar tile 0 via proj.Storage(): %v", err)
	}
	if tile0.TileIndex != 0 {
		t.Errorf("esperado índice 0, obtido %d", tile0.TileIndex)
	}
}
