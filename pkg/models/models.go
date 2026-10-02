package models

import (
	"fmt"
)

// Tileset representa um conjunto de até 256 tiles para um ambiente específico (ex: Overworld, Dungeon, Cidade).
type Tileset struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Tile representa um padrão de 8x8 pixels no padrão SCREEN 4 (Graphic 3) do V9938.
type Tile struct {
	ID                  int64         `json:"id"`
	TilesetID           int64         `json:"tileset_id"`
	TileIndex           int           `json:"tile_index"` // 0 a 255 no tileset local
	PatternBytes        []byte        `json:"pattern_bytes"` // 8 bytes (1 bit por pixel, linha a linha)
	ColorBytes          []byte        `json:"color_bytes"`   // 8 bytes (nibble alto = Fg, nibble baixo = Bg)
	CollisionType       CollisionType `json:"collision_type"`
	AnimationNextTileID *int64        `json:"animation_next_tile_id,omitempty"`
}

// NewTile inicializa um tile vazio com 8 bytes de padrão e 8 bytes de cor padrão (Fg=15 Branco, Bg=1 Preto).
func NewTile(tilesetID int64, tileIndex int) *Tile {
	color := make([]byte, TileColorSize)
	for i := range color {
		color[i] = 0xF1 // Fg: 15 (Branco), Bg: 1 (Preto)
	}
	return &Tile{
		TilesetID:     tilesetID,
		TileIndex:     tileIndex,
		PatternBytes:  make([]byte, TilePatternSize),
		ColorBytes:    color,
		CollisionType: CollisionPassable,
	}
}

// GetPixel retorna se o pixel na coordenada (x, y) do tile de 8x8 está aceso (1) ou apagado (0).
func (t *Tile) GetPixel(x, y int) (bool, error) {
	if x < 0 || x >= TileWidth || y < 0 || y >= TileHeight {
		return false, fmt.Errorf("coordenadas fora do tile 8x8: (%d, %d)", x, y)
	}
	if len(t.PatternBytes) != TilePatternSize {
		return false, fmt.Errorf("tamanho inválido de PatternBytes: %d (esperado %d)", len(t.PatternBytes), TilePatternSize)
	}
	lineByte := t.PatternBytes[y]
	shift := 7 - x
	return (lineByte & (1 << shift)) != 0, nil
}

// SetPixel altera o estado de um pixel no tile de 8x8.
func (t *Tile) SetPixel(x, y int, on bool) error {
	if x < 0 || x >= TileWidth || y < 0 || y >= TileHeight {
		return fmt.Errorf("coordenadas fora do tile 8x8: (%d, %d)", x, y)
	}
	if len(t.PatternBytes) != TilePatternSize {
		t.PatternBytes = make([]byte, TilePatternSize)
	}
	shift := 7 - x
	if on {
		t.PatternBytes[y] |= (1 << shift)
	} else {
		t.PatternBytes[y] &= ^(1 << shift)
	}
	return nil
}

// GetColors retorna as cores de primeiro plano (Fg) e fundo (Bg) para a linha de varredura y (0..7).
func (t *Tile) GetColors(y int) (fg byte, bg byte, err error) {
	if y < 0 || y >= TileHeight {
		return 0, 0, fmt.Errorf("linha fora do limite do tile (0..7): %d", y)
	}
	if len(t.ColorBytes) != TileColorSize {
		return 0, 0, fmt.Errorf("tamanho inválido de ColorBytes: %d", len(t.ColorBytes))
	}
	c := t.ColorBytes[y]
	return (c >> 4) & 0x0F, c & 0x0F, nil
}

// SetColors define as cores de primeiro plano (Fg) e fundo (Bg) para a linha de varredura y (0..7).
func (t *Tile) SetColors(y int, fg, bg byte) error {
	if y < 0 || y >= TileHeight {
		return fmt.Errorf("linha fora do limite do tile (0..7): %d", y)
	}
	if len(t.ColorBytes) != TileColorSize {
		t.ColorBytes = make([]byte, TileColorSize)
	}
	t.ColorBytes[y] = ((fg & 0x0F) << 4) | (bg & 0x0F)
	return nil
}

// Sprite representa um sprite de 16x16 pixels no Modo 2 do V9938 (MSX 2).
// No hardware MSX, os 32 bytes de padrão são organizados em dois blocos verticais de 16 bytes:
// - Bytes 0..15: Coluna esquerda (X: 0..7, Y: 0..15)
// - Bytes 16..31: Coluna direita (X: 8..15, Y: 0..15)
type Sprite struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PatternBytes []byte `json:"pattern_bytes"` // 32 bytes
	ColorBytes   []byte `json:"color_bytes"`   // 16 bytes (1 cor por linha de varredura de 0..15)
}

// NewSprite cria um novo sprite 16x16 inicializado com cor 15 (Branco).
func NewSprite(name string) *Sprite {
	colors := make([]byte, SpriteColorSize)
	for i := range colors {
		colors[i] = 0x0F // Cor 15
	}
	return &Sprite{
		Name:         name,
		PatternBytes: make([]byte, SpritePatternSize),
		ColorBytes:   colors,
	}
}

// GetPixel retorna o valor do pixel (0 ou 1) na grade de 16x16 do sprite.
func (s *Sprite) GetPixel(x, y int) (bool, error) {
	if x < 0 || x >= SpriteWidth || y < 0 || y >= SpriteHeight {
		return false, fmt.Errorf("coordenadas fora do sprite 16x16: (%d, %d)", x, y)
	}
	if len(s.PatternBytes) != SpritePatternSize {
		return false, fmt.Errorf("tamanho inválido de PatternBytes no sprite: %d", len(s.PatternBytes))
	}

	byteIndex := y
	bitCol := x
	if x >= 8 {
		byteIndex += 16
		bitCol = x - 8
	}

	shift := 7 - bitCol
	return (s.PatternBytes[byteIndex] & (1 << shift)) != 0, nil
}

// SetPixel altera o pixel na coordenada (x, y) de um sprite 16x16 seguindo o formato V9938.
func (s *Sprite) SetPixel(x, y int, on bool) error {
	if x < 0 || x >= SpriteWidth || y < 0 || y >= SpriteHeight {
		return fmt.Errorf("coordenadas fora do sprite 16x16: (%d, %d)", x, y)
	}
	if len(s.PatternBytes) != SpritePatternSize {
		s.PatternBytes = make([]byte, SpritePatternSize)
	}

	byteIndex := y
	bitCol := x
	if x >= 8 {
		byteIndex += 16
		bitCol = x - 8
	}

	shift := 7 - bitCol
	if on {
		s.PatternBytes[byteIndex] |= (1 << shift)
	} else {
		s.PatternBytes[byteIndex] &= ^(1 << shift)
	}
	return nil
}

// GetLineColor retorna a cor do scanline y (0..15) do sprite.
func (s *Sprite) GetLineColor(y int) (byte, error) {
	if y < 0 || y >= SpriteHeight {
		return 0, fmt.Errorf("scanline fora do limite do sprite (0..15): %d", y)
	}
	if len(s.ColorBytes) != SpriteColorSize {
		return 0, fmt.Errorf("tamanho inválido de ColorBytes no sprite: %d", len(s.ColorBytes))
	}
	return s.ColorBytes[y] & 0x0F, nil
}

// SetLineColor define a cor do scanline y (0..15) do sprite.
func (s *Sprite) SetLineColor(y int, color byte) error {
	if y < 0 || y >= SpriteHeight {
		return fmt.Errorf("scanline fora do limite do sprite (0..15): %d", y)
	}
	if len(s.ColorBytes) != SpriteColorSize {
		s.ColorBytes = make([]byte, SpriteColorSize)
	}
	s.ColorBytes[y] = color & 0x0F
	return nil
}

// Room representa uma sala ou tela estática do mundo (matriz de 32 x 18 tiles = 576 bytes).
type Room struct {
	ID          int64  `json:"id"`
	WorldX      int    `json:"world_x"`
	WorldY      int    `json:"world_y"`
	Name        string `json:"name"`
	TilesetID   int64  `json:"tileset_id"`
	TileMatrix  []byte `json:"tile_matrix"` // Exatamente 576 bytes
	NorthRoomID *int64 `json:"north_room_id,omitempty"`
	SouthRoomID *int64 `json:"south_room_id,omitempty"`
	EastRoomID  *int64 `json:"east_room_id,omitempty"`
	WestRoomID  *int64 `json:"west_room_id,omitempty"`
}

// NewRoom instancia uma nova sala vazia preenchida com o tile 0.
func NewRoom(worldX, worldY int, name string, tilesetID int64) *Room {
	return &Room{
		WorldX:     worldX,
		WorldY:     worldY,
		Name:       name,
		TilesetID:  tilesetID,
		TileMatrix: make([]byte, RoomMatrixSize),
	}
}

// GetTile retorna o índice do tile posicionado nas coordenadas (x, y) da sala (x: 0..31, y: 0..17).
func (r *Room) GetTile(x, y int) (byte, error) {
	if x < 0 || x >= RoomWidth || y < 0 || y >= RoomHeight {
		return 0, fmt.Errorf("coordenadas fora do grid da sala (32x18): (%d, %d)", x, y)
	}
	if len(r.TileMatrix) != RoomMatrixSize {
		return 0, fmt.Errorf("matriz de tiles com tamanho inválido: %d (esperado %d)", len(r.TileMatrix), RoomMatrixSize)
	}
	idx := y*RoomWidth + x
	return r.TileMatrix[idx], nil
}

// SetTile define o índice do tile posicionado nas coordenadas (x, y) da sala.
func (r *Room) SetTile(x, y int, tileIndex byte) error {
	if x < 0 || x >= RoomWidth || y < 0 || y >= RoomHeight {
		return fmt.Errorf("coordenadas fora do grid da sala (32x18): (%d, %d)", x, y)
	}
	if len(r.TileMatrix) != RoomMatrixSize {
		r.TileMatrix = make([]byte, RoomMatrixSize)
	}
	idx := y*RoomWidth + x
	r.TileMatrix[idx] = tileIndex
	return nil
}

// Entity representa um ator, NPC, baú ou gatilho em uma sala.
type Entity struct {
	ID            int64        `json:"id"`
	RoomID        int64        `json:"room_id"`
	Name          string       `json:"name"`
	PosX          int          `json:"pos_x"` // 0..31
	PosY          int          `json:"pos_y"` // 0..17
	SpriteID      *int64       `json:"sprite_id,omitempty"`
	BehaviorType  BehaviorType `json:"behavior_type"`
	EventScriptID *int64       `json:"event_script_id,omitempty"`
}

// Script representa um script de evento no jogo.
type Script struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	SourceCode string `json:"source_code"`
	Bytecode   []byte `json:"bytecode,omitempty"`
}

// StringEntry representa uma mensagem, diálogo ou linha de texto no jogo.
type StringEntry struct {
	ID          int64  `json:"id"`
	ContextTag  string `json:"context_tag"`
	TextContent string `json:"text_content"`
}

// HeroClass define a classe base de personagem do jogador (ex: Guerreiro, Mago, Ladino).
type HeroClass struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BaseHP  int    `json:"base_hp"`
	BaseMP  int    `json:"base_mp"`
	BaseAtk int    `json:"base_atk"`
	BaseDef int    `json:"base_def"`
}

// Item define um item no jogo (armas, armaduras, consumíveis, etc.).
type Item struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	ItemType      ItemType `json:"item_type"`
	ModifierStat  int      `json:"modifier_stat"`
	ModifierValue int      `json:"modifier_value"`
	Price         int      `json:"price"`
}
