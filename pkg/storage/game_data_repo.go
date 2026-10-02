package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// GameDataRepository gerencia scripts, tabelas de textos, classes de heróis e itens.
type GameDataRepository struct {
	db *sql.DB
}

// --- Scripts ---

// CreateScript insere um novo script e popula seu ID.
func (r *GameDataRepository) CreateScript(ctx context.Context, s *models.Script) error {
	query := `INSERT INTO scripts (name, source_code, bytecode) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.Name, s.SourceCode, s.Bytecode)
	if err != nil {
		return fmt.Errorf("falha ao criar script: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

// GetScript busca um script pelo ID.
func (r *GameDataRepository) GetScript(ctx context.Context, id int64) (*models.Script, error) {
	query := `SELECT id, name, source_code, bytecode FROM scripts WHERE id = ?`
	s := &models.Script{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.SourceCode, &s.Bytecode)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("script com ID %d não encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar script: %w", err)
	}
	return s, nil
}

// ListScripts lista todos os scripts cadastrados.
func (r *GameDataRepository) ListScripts(ctx context.Context) ([]*models.Script, error) {
	query := `SELECT id, name, source_code, bytecode FROM scripts ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar scripts: %w", err)
	}
	defer rows.Close()

	var list []*models.Script
	for rows.Next() {
		s := &models.Script{}
		if err := rows.Scan(&s.ID, &s.Name, &s.SourceCode, &s.Bytecode); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// UpdateScript atualiza o código fonte ou bytecode de um script.
func (r *GameDataRepository) UpdateScript(ctx context.Context, s *models.Script) error {
	query := `UPDATE scripts SET name = ?, source_code = ?, bytecode = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, s.Name, s.SourceCode, s.Bytecode, s.ID)
	if err != nil {
		return fmt.Errorf("falha ao atualizar script %d: %w", s.ID, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("script %d não encontrado para atualização", s.ID)
	}
	return nil
}

// DeleteScript remove um script.
func (r *GameDataRepository) DeleteScript(ctx context.Context, id int64) error {
	query := `DELETE FROM scripts WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar script %d: %w", id, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("script %d não encontrado para exclusão", id)
	}
	return nil
}

// --- Textos & Diálogos (string_table) ---

// CreateString insere uma nova entrada na tabela de textos.
func (r *GameDataRepository) CreateString(ctx context.Context, entry *models.StringEntry) error {
	query := `INSERT INTO string_table (context_tag, text_content) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, entry.ContextTag, entry.TextContent)
	if err != nil {
		return fmt.Errorf("falha ao criar string: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	entry.ID = id
	return nil
}

// GetString busca uma string pelo ID.
func (r *GameDataRepository) GetString(ctx context.Context, id int64) (*models.StringEntry, error) {
	query := `SELECT id, context_tag, text_content FROM string_table WHERE id = ?`
	entry := &models.StringEntry{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&entry.ID, &entry.ContextTag, &entry.TextContent)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("string com ID %d não encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar string: %w", err)
	}
	return entry, nil
}

// ListStrings retorna todas as strings cadastradas.
func (r *GameDataRepository) ListStrings(ctx context.Context) ([]*models.StringEntry, error) {
	query := `SELECT id, context_tag, text_content FROM string_table ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar strings: %w", err)
	}
	defer rows.Close()

	var list []*models.StringEntry
	for rows.Next() {
		entry := &models.StringEntry{}
		if err := rows.Scan(&entry.ID, &entry.ContextTag, &entry.TextContent); err != nil {
			return nil, err
		}
		list = append(list, entry)
	}
	return list, rows.Err()
}

// --- Classes de Heróis (hero_classes) ---

// CreateHeroClass cadastra uma nova classe de personagem.
func (r *GameDataRepository) CreateHeroClass(ctx context.Context, h *models.HeroClass) error {
	query := `INSERT INTO hero_classes (name, base_hp, base_mp, base_atk, base_def) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, h.Name, h.BaseHP, h.BaseMP, h.BaseAtk, h.BaseDef)
	if err != nil {
		return fmt.Errorf("falha ao criar classe de herói: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	h.ID = id
	return nil
}

// GetHeroClass busca uma classe de herói pelo ID.
func (r *GameDataRepository) GetHeroClass(ctx context.Context, id int64) (*models.HeroClass, error) {
	query := `SELECT id, name, base_hp, base_mp, base_atk, base_def FROM hero_classes WHERE id = ?`
	h := &models.HeroClass{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&h.ID, &h.Name, &h.BaseHP, &h.BaseMP, &h.BaseAtk, &h.BaseDef)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("classe de herói com ID %d não encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar classe de herói: %w", err)
	}
	return h, nil
}

// ListHeroClasses lista todas as classes de personagens cadastradas.
func (r *GameDataRepository) ListHeroClasses(ctx context.Context) ([]*models.HeroClass, error) {
	query := `SELECT id, name, base_hp, base_mp, base_atk, base_def FROM hero_classes ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar classes de heróis: %w", err)
	}
	defer rows.Close()

	var list []*models.HeroClass
	for rows.Next() {
		h := &models.HeroClass{}
		if err := rows.Scan(&h.ID, &h.Name, &h.BaseHP, &h.BaseMP, &h.BaseAtk, &h.BaseDef); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

// --- Itens (items) ---

// CreateItem insere um novo item no catálogo do jogo.
func (r *GameDataRepository) CreateItem(ctx context.Context, item *models.Item) error {
	query := `INSERT INTO items (name, item_type, modifier_stat, modifier_value, price) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, item.Name, int(item.ItemType), item.ModifierStat, item.ModifierValue, item.Price)
	if err != nil {
		return fmt.Errorf("falha ao criar item: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = id
	return nil
}

// GetItem busca um item pelo ID.
func (r *GameDataRepository) GetItem(ctx context.Context, id int64) (*models.Item, error) {
	query := `SELECT id, name, item_type, modifier_stat, modifier_value, price FROM items WHERE id = ?`
	item := &models.Item{}
	var iType int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&iType,
		&item.ModifierStat,
		&item.ModifierValue,
		&item.Price,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item com ID %d não encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar item: %w", err)
	}
	item.ItemType = models.ItemType(iType)
	return item, nil
}

// ListItems lista todos os itens cadastrados no jogo.
func (r *GameDataRepository) ListItems(ctx context.Context) ([]*models.Item, error) {
	query := `SELECT id, name, item_type, modifier_stat, modifier_value, price FROM items ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar itens: %w", err)
	}
	defer rows.Close()

	var list []*models.Item
	for rows.Next() {
		item := &models.Item{}
		var iType int
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&iType,
			&item.ModifierStat,
			&item.ModifierValue,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		item.ItemType = models.ItemType(iType)
		list = append(list, item)
	}
	return list, rows.Err()
}
