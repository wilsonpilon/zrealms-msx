package project

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zrealm-msx/zrealm/pkg/project/migrations"
	"github.com/zrealm-msx/zrealm/pkg/storage"
	_ "modernc.org/sqlite"
)

// Project representa um projeto Z-Realm aberto e conectado ao seu arquivo SQLite (.rpgproj).
type Project struct {
	db      *sql.DB
	path    string
	storage *storage.Storage
}

// DefaultSettings contém as chaves e valores padrão na criação de um novo projeto.
var DefaultSettings = map[string]string{
	"version":         "0.1.0",
	"target_platform": "MSX2_MSXDOS2",
	"screen_mode":     "SCREEN4",
	"mapper_type":     "MSX_DOS2_MAPPER",
	"viewport_width":  "32",
	"viewport_height": "18",
}

// Create cria um novo arquivo .rpgproj, aplica as migrações DDL e inicializa os metadados.
func Create(filePath string, projectName string) (*Project, error) {
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("o arquivo de projeto já existe em: %s", filePath)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório do projeto: %w", err)
	}

	db, err := openDB(filePath)
	if err != nil {
		return nil, err
	}

	p := &Project{
		db:      db,
		path:    filePath,
		storage: storage.New(db),
	}

	if err := p.migrate(); err != nil {
		_ = db.Close()
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("falha ao aplicar migrações na criação: %w", err)
	}

	if err := p.initDefaults(projectName); err != nil {
		_ = db.Close()
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("falha ao inicializar configurações padrão: %w", err)
	}

	return p, nil
}

// Open abre um arquivo .rpgproj existente, verifica a integridade e aplica migrações pendentes.
func Open(filePath string) (*Project, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de projeto não encontrado: %s", filePath)
	}

	db, err := openDB(filePath)
	if err != nil {
		return nil, err
	}

	p := &Project{
		db:      db,
		path:    filePath,
		storage: storage.New(db),
	}

	if err := p.migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("falha ao atualizar schema de migrações: %w", err)
	}

	if err := p.ValidateIntegrity(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("falha na validação de integridade do projeto: %w", err)
	}

	return p, nil
}

func openDB(filePath string) (*sql.DB, error) {
	// Parâmetros de conexão otimizados para SQLite com foreign keys ativadas
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", filePath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco de dados SQLite: %w", err)
	}

	// Garante que o pool mantenha concorrência controlada para SQLite WAL
	db.SetMaxOpenConns(1)

	// Validação explícita de pragmas essenciais
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("falha ao ativar PRAGMA foreign_keys: %w", err)
	}

	return db, nil
}

// Close encerra a conexão com o banco de dados.
func (p *Project) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// DB retorna a conexão subjacente *sql.DB para operações de baixo nível.
func (p *Project) DB() *sql.DB {
	return p.db
}

// Storage retorna os repositórios tipados do projeto Z-Realm.
func (p *Project) Storage() *storage.Storage {
	return p.storage
}

// Path retorna o caminho absoluto ou relativo do arquivo de projeto.
func (p *Project) Path() string {
	return p.path
}

// migrate executa todas as migrações SQL embutidas que ainda não foram aplicadas.
func (p *Project) migrate() error {
	// Garante a existência da tabela de controle de migrações
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("falha ao criar tabela schema_migrations: %w", err)
	}

	allMigrations, err := migrations.LoadMigrations()
	if err != nil {
		return err
	}

	for _, m := range allMigrations {
		var exists bool
		err := p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", m.Version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("falha ao checar status da migração %d: %w", m.Version, err)
		}

		if !exists {
			tx, err := p.db.Begin()
			if err != nil {
				return fmt.Errorf("falha ao abrir transação para migração %d: %w", m.Version, err)
			}

			if _, err := tx.Exec(m.SQL); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("falha ao executar migração %s (v%d): %w", m.Name, m.Version, err)
			}

			if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.Version); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("falha ao registrar migração %d: %w", m.Version, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("falha ao efetivar transação da migração %d: %w", m.Version, err)
			}
		}
	}

	return nil
}

func (p *Project) initDefaults(projectName string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO project_settings (key, value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec("name", projectName); err != nil {
		return err
	}
	if _, err := stmt.Exec("created_at", time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}

	for k, v := range DefaultSettings {
		if _, err := stmt.Exec(k, v); err != nil {
			return err
		}
	}

	// Cria um tileset padrão inicial (Overworld)
	res, err := tx.Exec("INSERT INTO tilesets (name, description) VALUES (?, ?)", "Overworld", "Tileset principal de exploração")
	if err != nil {
		return fmt.Errorf("falha ao criar tileset inicial: %w", err)
	}

	tilesetID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Cria um tile 0 padrão (Passável, padrão monocromático 0x00, cor padrão MSX2 Fg=15 Bg=1)
	emptyPattern := make([]byte, 8)
	defaultColor := []byte{0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1}
	_, err = tx.Exec(`
		INSERT INTO tiles (tileset_id, tile_index, pattern_bytes, color_bytes, collision_type)
		VALUES (?, 0, ?, ?, 0)
	`, tilesetID, emptyPattern, defaultColor)
	if err != nil {
		return fmt.Errorf("falha ao criar tile base: %w", err)
	}

	return tx.Commit()
}

// ValidateIntegrity executa checagens de integridade física e integridade de chaves estrangeiras.
func (p *Project) ValidateIntegrity() error {
	// Checagem de integridade estrutural do SQLite
	var integrityResult string
	if err := p.db.QueryRow("PRAGMA quick_check;").Scan(&integrityResult); err != nil {
		return fmt.Errorf("erro ao executar quick_check: %w", err)
	}
	if integrityResult != "ok" {
		return fmt.Errorf("corrupção de banco detectada: %s", integrityResult)
	}

	// Checagem de integridade referencial (FKs)
	rows, err := p.db.Query("PRAGMA foreign_key_check;")
	if err != nil {
		return fmt.Errorf("erro ao executar foreign_key_check: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var table, fkid, target string
		var rowid int64
		_ = rows.Scan(&table, &rowid, &target, &fkid)
		return fmt.Errorf("violação de chave estrangeira encontrada na tabela '%s' (rowid: %d, alvo: %s)", table, rowid, target)
	}

	// Verifica se a tabela de configurações existe e contém o nome do projeto
	var name string
	err = p.db.QueryRow("SELECT value FROM project_settings WHERE key = 'name'").Scan(&name)
	if err != nil {
		return fmt.Errorf("projeto inválido: chave 'name' ausente em project_settings: %w", err)
	}

	return nil
}

// GetSetting retorna o valor de uma chave em project_settings.
func (p *Project) GetSetting(key string) (string, error) {
	var val string
	err := p.db.QueryRow("SELECT value FROM project_settings WHERE key = ?", key).Scan(&val)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("chave de configuração '%s' não encontrada", key)
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// SetSetting define ou atualiza uma chave em project_settings.
func (p *Project) SetSetting(key string, val string) error {
	_, err := p.db.Exec(`
		INSERT INTO project_settings (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, val)
	return err
}
