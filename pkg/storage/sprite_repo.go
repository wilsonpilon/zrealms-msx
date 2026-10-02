package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// SpriteRepository gerencia a persistência de sprites de 16x16 (Modo 2 do V9938).
type SpriteRepository struct {
	db *sql.DB
}

// CreateSprite insere um novo sprite no banco e preenche seu ID gerado.
func (r *SpriteRepository) CreateSprite(ctx context.Context, s *models.Sprite) error {
	if len(s.PatternBytes) != models.SpritePatternSize {
		return fmt.Errorf("tamanho inválido de PatternBytes para sprite: %d (esperado %d)", len(s.PatternBytes), models.SpritePatternSize)
	}
	if len(s.ColorBytes) != models.SpriteColorSize {
		return fmt.Errorf("tamanho inválido de ColorBytes para sprite: %d (esperado %d)", len(s.ColorBytes), models.SpriteColorSize)
	}

	query := `INSERT INTO sprites (name, pattern_bytes, color_bytes) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.Name, s.PatternBytes, s.ColorBytes)
	if err != nil {
		return fmt.Errorf("falha ao criar sprite: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

// GetSprite busca um sprite pelo ID.
func (r *SpriteRepository) GetSprite(ctx context.Context, id int64) (*models.Sprite, error) {
	query := `SELECT id, name, pattern_bytes, color_bytes FROM sprites WHERE id = ?`
	s := &models.Sprite{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.PatternBytes, &s.ColorBytes)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sprite com ID %d não encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar sprite: %w", err)
	}
	return s, nil
}

// ListSprites lista todos os sprites cadastrados no projeto.
func (r *SpriteRepository) ListSprites(ctx context.Context) ([]*models.Sprite, error) {
	query := `SELECT id, name, pattern_bytes, color_bytes FROM sprites ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar sprites: %w", err)
	}
	defer rows.Close()

	var list []*models.Sprite
	for rows.Next() {
		s := &models.Sprite{}
		if err := rows.Scan(&s.ID, &s.Name, &s.PatternBytes, &s.ColorBytes); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// UpdateSprite atualiza nome, padrão e cores de um sprite existente.
func (r *SpriteRepository) UpdateSprite(ctx context.Context, s *models.Sprite) error {
	if len(s.PatternBytes) != models.SpritePatternSize {
		return fmt.Errorf("tamanho inválido de PatternBytes para sprite: %d", len(s.PatternBytes))
	}
	if len(s.ColorBytes) != models.SpriteColorSize {
		return fmt.Errorf("tamanho inválido de ColorBytes para sprite: %d", len(s.ColorBytes))
	}

	query := `UPDATE sprites SET name = ?, pattern_bytes = ?, color_bytes = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, s.Name, s.PatternBytes, s.ColorBytes, s.ID)
	if err != nil {
		return fmt.Errorf("falha ao atualizar sprite %d: %w", s.ID, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("sprite com ID %d não encontrado para atualização", s.ID)
	}
	return nil
}

// DeleteSprite remove um sprite do banco.
func (r *SpriteRepository) DeleteSprite(ctx context.Context, id int64) error {
	query := `DELETE FROM sprites WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar sprite %d: %w", id, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("sprite com ID %d não encontrado para exclusão", id)
	}
	return nil
}
