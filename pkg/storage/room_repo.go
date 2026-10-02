package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// RoomRepository gerencia a persistência de salas e das entidades associadas a elas.
type RoomRepository struct {
	db *sql.DB
}

// CreateRoom insere uma nova sala no mundo (validando tamanho da matriz de 576 bytes).
func (r *RoomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	if len(room.TileMatrix) != models.RoomMatrixSize {
		return fmt.Errorf("tamanho inválido de TileMatrix: %d (esperado %d bytes)", len(room.TileMatrix), models.RoomMatrixSize)
	}

	query := `
		INSERT INTO rooms (world_x, world_y, name, tileset_id, tile_matrix, north_room_id, south_room_id, east_room_id, west_room_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query,
		room.WorldX,
		room.WorldY,
		room.Name,
		room.TilesetID,
		room.TileMatrix,
		room.NorthRoomID,
		room.SouthRoomID,
		room.EastRoomID,
		room.WestRoomID,
	)
	if err != nil {
		return fmt.Errorf("falha ao criar sala: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	room.ID = id
	return nil
}

// GetRoom busca uma sala pelo ID.
func (r *RoomRepository) GetRoom(ctx context.Context, id int64) (*models.Room, error) {
	query := `
		SELECT id, world_x, world_y, name, tileset_id, tile_matrix, north_room_id, south_room_id, east_room_id, west_room_id
		FROM rooms
		WHERE id = ?
	`
	room := &models.Room{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.WorldX,
		&room.WorldY,
		&room.Name,
		&room.TilesetID,
		&room.TileMatrix,
		&room.NorthRoomID,
		&room.SouthRoomID,
		&room.EastRoomID,
		&room.WestRoomID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sala com ID %d não encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar sala: %w", err)
	}
	return room, nil
}

// GetRoomByCoords busca a sala localizada em uma coordenada específica do mapa global.
func (r *RoomRepository) GetRoomByCoords(ctx context.Context, worldX, worldY int) (*models.Room, error) {
	query := `
		SELECT id, world_x, world_y, name, tileset_id, tile_matrix, north_room_id, south_room_id, east_room_id, west_room_id
		FROM rooms
		WHERE world_x = ? AND world_y = ?
	`
	room := &models.Room{}
	err := r.db.QueryRowContext(ctx, query, worldX, worldY).Scan(
		&room.ID,
		&room.WorldX,
		&room.WorldY,
		&room.Name,
		&room.TilesetID,
		&room.TileMatrix,
		&room.NorthRoomID,
		&room.SouthRoomID,
		&room.EastRoomID,
		&room.WestRoomID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sala na coordenada (%d, %d) não encontrada", worldX, worldY)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar sala por coordenadas: %w", err)
	}
	return room, nil
}

// ListRooms retorna todas as salas cadastradas.
func (r *RoomRepository) ListRooms(ctx context.Context) ([]*models.Room, error) {
	query := `
		SELECT id, world_x, world_y, name, tileset_id, tile_matrix, north_room_id, south_room_id, east_room_id, west_room_id
		FROM rooms
		ORDER BY world_y ASC, world_x ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar salas: %w", err)
	}
	defer rows.Close()

	var list []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID,
			&room.WorldX,
			&room.WorldY,
			&room.Name,
			&room.TilesetID,
			&room.TileMatrix,
			&room.NorthRoomID,
			&room.SouthRoomID,
			&room.EastRoomID,
			&room.WestRoomID,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, room)
	}
	return list, rows.Err()
}

// UpdateRoom atualiza os dados da sala e sua matriz de tiles.
func (r *RoomRepository) UpdateRoom(ctx context.Context, room *models.Room) error {
	if len(room.TileMatrix) != models.RoomMatrixSize {
		return fmt.Errorf("tamanho inválido de TileMatrix: %d", len(room.TileMatrix))
	}

	query := `
		UPDATE rooms SET
			world_x = ?,
			world_y = ?,
			name = ?,
			tileset_id = ?,
			tile_matrix = ?,
			north_room_id = ?,
			south_room_id = ?,
			east_room_id = ?,
			west_room_id = ?
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query,
		room.WorldX,
		room.WorldY,
		room.Name,
		room.TilesetID,
		room.TileMatrix,
		room.NorthRoomID,
		room.SouthRoomID,
		room.EastRoomID,
		room.WestRoomID,
		room.ID,
	)
	if err != nil {
		return fmt.Errorf("falha ao atualizar sala %d: %w", room.ID, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("sala %d não encontrada para atualização", room.ID)
	}
	return nil
}

// DeleteRoom remove uma sala e suas entidades (em cascata).
func (r *RoomRepository) DeleteRoom(ctx context.Context, id int64) error {
	query := `DELETE FROM rooms WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar sala %d: %w", id, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("sala %d não encontrada para exclusão", id)
	}
	return nil
}

// --- Entidades da Sala ---

// CreateEntity insere uma nova entidade associada a uma sala.
func (r *RoomRepository) CreateEntity(ctx context.Context, e *models.Entity) error {
	if e.PosX < 0 || e.PosX >= models.RoomWidth || e.PosY < 0 || e.PosY >= models.RoomHeight {
		return fmt.Errorf("posição (%d, %d) fora dos limites da sala (32x18)", e.PosX, e.PosY)
	}

	query := `
		INSERT INTO entities (room_id, name, pos_x, pos_y, sprite_id, behavior_type, event_script_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query,
		e.RoomID,
		e.Name,
		e.PosX,
		e.PosY,
		e.SpriteID,
		int(e.BehaviorType),
		e.EventScriptID,
	)
	if err != nil {
		return fmt.Errorf("falha ao criar entidade: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

// GetEntity busca uma entidade pelo ID.
func (r *RoomRepository) GetEntity(ctx context.Context, id int64) (*models.Entity, error) {
	query := `
		SELECT id, room_id, name, pos_x, pos_y, sprite_id, behavior_type, event_script_id
		FROM entities
		WHERE id = ?
	`
	e := &models.Entity{}
	var bType int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.RoomID,
		&e.Name,
		&e.PosX,
		&e.PosY,
		&e.SpriteID,
		&bType,
		&e.EventScriptID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entidade com ID %d não encontrada", id)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar entidade: %w", err)
	}
	e.BehaviorType = models.BehaviorType(bType)
	return e, nil
}

// ListEntitiesByRoom retorna todas as entidades presentes em uma sala.
func (r *RoomRepository) ListEntitiesByRoom(ctx context.Context, roomID int64) ([]*models.Entity, error) {
	query := `
		SELECT id, room_id, name, pos_x, pos_y, sprite_id, behavior_type, event_script_id
		FROM entities
		WHERE room_id = ?
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar entidades da sala %d: %w", roomID, err)
	}
	defer rows.Close()

	var list []*models.Entity
	for rows.Next() {
		e := &models.Entity{}
		var bType int
		err := rows.Scan(
			&e.ID,
			&e.RoomID,
			&e.Name,
			&e.PosX,
			&e.PosY,
			&e.SpriteID,
			&bType,
			&e.EventScriptID,
		)
		if err != nil {
			return nil, err
		}
		e.BehaviorType = models.BehaviorType(bType)
		list = append(list, e)
	}
	return list, rows.Err()
}

// UpdateEntity atualiza propriedades ou posição de uma entidade.
func (r *RoomRepository) UpdateEntity(ctx context.Context, e *models.Entity) error {
	if e.PosX < 0 || e.PosX >= models.RoomWidth || e.PosY < 0 || e.PosY >= models.RoomHeight {
		return fmt.Errorf("posição (%d, %d) fora dos limites da sala (32x18)", e.PosX, e.PosY)
	}

	query := `
		UPDATE entities SET
			room_id = ?,
			name = ?,
			pos_x = ?,
			pos_y = ?,
			sprite_id = ?,
			behavior_type = ?,
			event_script_id = ?
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query,
		e.RoomID,
		e.Name,
		e.PosX,
		e.PosY,
		e.SpriteID,
		int(e.BehaviorType),
		e.EventScriptID,
		e.ID,
	)
	if err != nil {
		return fmt.Errorf("falha ao atualizar entidade %d: %w", e.ID, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("entidade %d não encontrada para atualização", e.ID)
	}
	return nil
}

// DeleteEntity remove uma entidade.
func (r *RoomRepository) DeleteEntity(ctx context.Context, id int64) error {
	query := `DELETE FROM entities WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao deletar entidade %d: %w", id, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("entidade %d não encontrada para exclusão", id)
	}
	return nil
}
