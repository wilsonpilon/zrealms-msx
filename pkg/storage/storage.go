package storage

import (
	"database/sql"
)

// Storage é o agregador central de repositórios de dados do projeto Z-Realm.
type Storage struct {
	db       *sql.DB
	Tilesets *TilesetRepository
	Sprites  *SpriteRepository
	Rooms    *RoomRepository
	GameData *GameDataRepository
}

// New cria e inicializa todos os repositórios vinculados à conexão SQLite ativa.
func New(db *sql.DB) *Storage {
	return &Storage{
		db:       db,
		Tilesets: &TilesetRepository{db: db},
		Sprites:  &SpriteRepository{db: db},
		Rooms:    &RoomRepository{db: db},
		GameData: &GameDataRepository{db: db},
	}
}

// DB retorna a conexão direta com o banco para transações customizadas.
func (s *Storage) DB() *sql.DB {
	return s.db
}
