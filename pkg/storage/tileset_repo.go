package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// TilesetRepository gerencia a persistência de tilesets e seus respectivos tiles 8x8.
type TilesetRepository struct {
	db *sql.DB
}

// CreateTileset insere um novo tileset no banco e popula seu ID gerado.
func (r *TilesetRepository) CreateTileset(ctx context.Context, ts *models.Tileset) error {
	query := `INSERT INTO tilesets (name, description) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, ts.Name, ts.Description)
	if err != nil {
		return fmt.Errorf("falha ao criar tileset: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ts.ID = id
	return nil
}

// GetTileset busca um tileset pelo ID.
func (r *TilesetRepository) GetTileset(ctx context.Context, id int64) (*models.Tileset, error) {
	query := `SELECT id, name, description FROM tilesets WHERE id = ?`
	ts := &models.Tileset{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&ts.ID, &ts.Name, &ts.Description)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tileset com ID %d não encontrado", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar tileset: %w", err)
	}
	return ts, nil
}

// ListTilesets retorna todos os tilesets cadastrados no projeto.
func (r *TilesetRepository) ListTilesets(ctx context.Context) ([]*models.Tileset, error) {
	query := `SELECT id, name, description FROM tilesets ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar tilesets: %w", err)
	}
	defer rows.Close()

	var list []*models.Tileset
	for rows.Next() {
		ts := &models.Tileset{}
		if err := rows.Scan(&ts.ID, &ts.Name, &ts.Description); err != nil {
			return nil, err
		}
		list = append(list, ts)
	}
	return list, rows.Err()
}

// UpdateTileset atualiza o nome e descrição de um tileset.
func (r *TilesetRepository) UpdateTileset(ctx context.Context, ts *models.Tileset) error {
	query := `UPDATE tilesets SET name = ?, description = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, ts.Name, ts.Description, ts.ID)
	if err != nil {
		return fmt.Errorf("falha ao atualizar tileset %d: %w", ts.ID, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tileset %d não encontrado para atualização", ts.ID)
	}
	return nil
}

// DeleteTileset remove um tileset (os tiles são deletados em cascata via FK).
func (r *TilesetRepository) DeleteTileset(ctx context.Context, id int64) error {
	query := `DELETE FROM tilesets WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar tileset %d: %w", id, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("tileset %d não encontrado para exclusão", id)
	}
	return nil
}

// SaveTile insere ou atualiza um tile (Upsert baseado em UNIQUE(tileset_id, tile_index)).
func (r *TilesetRepository) SaveTile(ctx context.Context, tile *models.Tile) error {
	if len(tile.PatternBytes) != models.TilePatternSize {
		return fmt.Errorf("tamanho inválido de PatternBytes: %d", len(tile.PatternBytes))
	}
	if len(tile.ColorBytes) != models.TileColorSize {
		return fmt.Errorf("tamanho inválido de ColorBytes: %d", len(tile.ColorBytes))
	}

	query := `
		INSERT INTO tiles (tileset_id, tile_index, pattern_bytes, color_bytes, collision_type, animation_next_tile_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(tileset_id, tile_index) DO UPDATE SET
			pattern_bytes = excluded.pattern_bytes,
			color_bytes = excluded.color_bytes,
			collision_type = excluded.collision_type,
			animation_next_tile_id = excluded.animation_next_tile_id
		RETURNING id;
	`
	err := r.db.QueryRowContext(ctx, query,
		tile.TilesetID,
		tile.TileIndex,
		tile.PatternBytes,
		tile.ColorBytes,
		int(tile.CollisionType),
		tile.AnimationNextTileID,
	).Scan(&tile.ID)

	if err != nil {
		return fmt.Errorf("falha ao salvar tile (tileset=%d, index=%d): %w", tile.TilesetID, tile.TileIndex, err)
	}
	return nil
}

// GetTile busca um tile pelo ID do tileset e seu índice de 0 a 255.
func (r *TilesetRepository) GetTile(ctx context.Context, tilesetID int64, tileIndex int) (*models.Tile, error) {
	query := `
		SELECT id, tileset_id, tile_index, pattern_bytes, color_bytes, collision_type, animation_next_tile_id
		FROM tiles
		WHERE tileset_id = ? AND tile_index = ?
	`
	tile := &models.Tile{}
	var colType int
	err := r.db.QueryRowContext(ctx, query, tilesetID, tileIndex).Scan(
		&tile.ID,
		&tile.TilesetID,
		&tile.TileIndex,
		&tile.PatternBytes,
		&tile.ColorBytes,
		&colType,
		&tile.AnimationNextTileID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tile não encontrado (tileset=%d, index=%d)", tilesetID, tileIndex)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar tile: %w", err)
	}
	tile.CollisionType = models.CollisionType(colType)
	return tile, nil
}

// ListTiles retorna todos os tiles cadastrados para um determinado tileset.
func (r *TilesetRepository) ListTiles(ctx context.Context, tilesetID int64) ([]*models.Tile, error) {
	query := `
		SELECT id, tileset_id, tile_index, pattern_bytes, color_bytes, collision_type, animation_next_tile_id
		FROM tiles
		WHERE tileset_id = ?
		ORDER BY tile_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, tilesetID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar tiles: %w", err)
	}
	defer rows.Close()

	var list []*models.Tile
	for rows.Next() {
		tile := &models.Tile{}
		var colType int
		err := rows.Scan(
			&tile.ID,
			&tile.TilesetID,
			&tile.TileIndex,
			&tile.PatternBytes,
			&tile.ColorBytes,
			&colType,
			&tile.AnimationNextTileID,
		)
		if err != nil {
			return nil, err
		}
		tile.CollisionType = models.CollisionType(colType)
		list = append(list, tile)
	}
	return list, rows.Err()
}
