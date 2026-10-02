package exporter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/storage"
)

// ExportOptions define opções configuráveis para o processo de exportação.
type ExportOptions struct {
	OutputDir        string
	ExportBankFiles  bool // Se verdadeiro, exporta também arquivos individuais SEGxx.BNK (16KB cada)
}

// Export compila todos os recursos do banco SQLite do projeto em binários estáticos para o MSX 2.
func Export(p *project.Project, opts ExportOptions) (*ExportResult, error) {
	if p == nil {
		return nil, fmt.Errorf("projeto não pode ser nulo")
	}
	if opts.OutputDir == "" {
		opts.OutputDir = filepath.Join(filepath.Dir(p.Path()), "build_msx")
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de saída: %w", err)
	}

	ctx := context.Background()
	st := p.Storage()
	packer := NewPacker()

	// 1. Exporta Tilesets (4.608 bytes por tileset)
	tilesets, err := st.Tilesets.ListTilesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar tilesets: %w", err)
	}

	for _, ts := range tilesets {
		tsData, err := exportTileset(ctx, st, ts.ID)
		if err != nil {
			return nil, fmt.Errorf("falha ao compilar tileset %d: %w", ts.ID, err)
		}
		if _, err := packer.WriteResource(ResTypeTileset, uint16(ts.ID), tsData); err != nil {
			return nil, err
		}
	}

	// 2. Exporta Sprites (48 bytes por sprite)
	sprites, err := st.Sprites.ListSprites(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar sprites: %w", err)
	}

	for _, spr := range sprites {
		sprData := exportSprite(spr)
		if _, err := packer.WriteResource(ResTypeSprite, uint16(spr.ID), sprData); err != nil {
			return nil, err
		}
	}

	// 3. Exporta Salas e suas Entidades (656 bytes por sala)
	rooms, err := st.Rooms.ListRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar salas: %w", err)
	}

	for _, room := range rooms {
		roomData, err := exportRoom(ctx, st, room)
		if err != nil {
			return nil, fmt.Errorf("falha ao compilar sala %d: %w", room.ID, err)
		}
		if _, err := packer.WriteResource(ResTypeRoom, uint16(room.ID), roomData); err != nil {
			return nil, err
		}
	}

	// 4. Exporta Strings / Textos de Diálogo
	stringsList, err := st.GameData.ListStrings(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar strings: %w", err)
	}

	if len(stringsList) > 0 {
		strData := exportStrings(stringsList)
		if _, err := packer.WriteResource(ResTypeStrings, 1, strData); err != nil {
			return nil, err
		}
	}

	// 5. Determina parâmetros iniciais do herói e da sala de partida
	initialRoomID := uint16(1)
	if len(rooms) > 0 {
		initialRoomID = uint16(rooms[0].ID)
	}
	if val, err := p.GetSetting("initial_room_id"); err == nil && val != "" {
		if id, parseErr := strconv.Atoi(val); parseErr == nil && id > 0 {
			initialRoomID = uint16(id)
		}
	}

	initialHeroX := uint8(16)
	initialHeroY := uint8(9)
	initialTileset := uint8(1)
	if len(tilesets) > 0 {
		initialTileset = uint8(tilesets[0].ID)
	}

	// 6. Constrói HEADER.BIN e grava no disco
	headerBytes, err := packer.BuildMasterHeader(initialRoomID, initialHeroX, initialHeroY, initialTileset)
	if err != nil {
		return nil, fmt.Errorf("falha ao construir HEADER.BIN: %w", err)
	}

	headerFile := filepath.Join(opts.OutputDir, "HEADER.BIN")
	if err := os.WriteFile(headerFile, headerBytes, 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar HEADER.BIN: %w", err)
	}

	// 7. Constrói GAME.DAT (todos os segmentos contíguos de 16 KB)
	flatData := packer.GetFlatData()
	dataFile := filepath.Join(opts.OutputDir, "GAME.DAT")
	if err := os.WriteFile(dataFile, flatData, 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar GAME.DAT: %w", err)
	}

	// 8. Opcional: grava segmentos individuais SEGxx.BNK
	if opts.ExportBankFiles {
		for i, seg := range packer.GetSegments() {
			segName := fmt.Sprintf("SEG%02d.BNK", i)
			segPath := filepath.Join(opts.OutputDir, segName)
			if err := os.WriteFile(segPath, seg, 0644); err != nil {
				return nil, fmt.Errorf("falha ao gravar arquivo de banco %s: %w", segName, err)
			}
		}
	}

	return &ExportResult{
		HeaderPath:     headerFile,
		DataPath:       dataFile,
		TotalSegments:  packer.SegmentCount(),
		ResourceCount:  packer.ResourceCount(),
		TotalDataBytes: len(flatData),
	}, nil
}

func exportTileset(ctx context.Context, st *storage.Storage, tilesetID int64) ([]byte, error) {
	tiles, err := st.Tilesets.ListTiles(ctx, tilesetID)
	if err != nil {
		return nil, err
	}

	binTs := &BinaryTileset{}
	// Preenche animações com auto-referência por padrão
	for i := range binTs.AnimationTable {
		binTs.AnimationTable[i] = byte(i)
	}

	for _, tile := range tiles {
		if tile.TileIndex < 0 || tile.TileIndex >= MaxTilesPerTileset {
			continue
		}
		offset := tile.TileIndex * 8

		// Padrão (8 bytes)
		if len(tile.PatternBytes) == 8 {
			copy(binTs.PatternTable[offset:offset+8], tile.PatternBytes)
		}
		// Cores (8 bytes)
		if len(tile.ColorBytes) == 8 {
			copy(binTs.ColorTable[offset:offset+8], tile.ColorBytes)
		}
		// Colisão
		binTs.CollisionTable[tile.TileIndex] = byte(tile.CollisionType)

		// Animação encadeada
		if tile.AnimationNextTileID != nil {
			// Procura índice do próximo tile
			for _, other := range tiles {
				if other.ID == *tile.AnimationNextTileID {
					binTs.AnimationTable[tile.TileIndex] = byte(other.TileIndex)
					break
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	buf.Write(binTs.PatternTable[:])
	buf.Write(binTs.ColorTable[:])
	buf.Write(binTs.CollisionTable[:])
	buf.Write(binTs.AnimationTable[:])

	return buf.Bytes(), nil
}

func exportSprite(spr *models.Sprite) []byte {
	buf := make([]byte, 48) // 32 bytes pattern + 16 bytes color
	if len(spr.PatternBytes) == 32 {
		copy(buf[0:32], spr.PatternBytes)
	}
	if len(spr.ColorBytes) == 16 {
		copy(buf[32:48], spr.ColorBytes)
	}
	return buf
}

func exportRoom(ctx context.Context, st *storage.Storage, room *models.Room) ([]byte, error) {
	buf := new(bytes.Buffer)

	// 1. Conexões cardeais (0xFFFF para bloqueado / sem sala vizinha)
	north := uint16(0xFFFF)
	if room.NorthRoomID != nil {
		north = uint16(*room.NorthRoomID)
	}
	south := uint16(0xFFFF)
	if room.SouthRoomID != nil {
		south = uint16(*room.SouthRoomID)
	}
	east := uint16(0xFFFF)
	if room.EastRoomID != nil {
		east = uint16(*room.EastRoomID)
	}
	west := uint16(0xFFFF)
	if room.WestRoomID != nil {
		west = uint16(*room.WestRoomID)
	}

	// 2. Busca entidades da sala (máximo 8)
	entities, err := st.Rooms.ListEntitiesByRoom(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	entCount := len(entities)
	if entCount > MaxEntitiesPerRoom {
		entCount = MaxEntitiesPerRoom
	}

	// Grava Cabeçalho da Sala (16 bytes)
	_ = binary.Write(buf, binary.LittleEndian, north)
	_ = binary.Write(buf, binary.LittleEndian, south)
	_ = binary.Write(buf, binary.LittleEndian, east)
	_ = binary.Write(buf, binary.LittleEndian, west)
	_ = buf.WriteByte(byte(room.TilesetID))
	_ = buf.WriteByte(byte(entCount))

	roomReserved := make([]byte, 6)
	_, _ = buf.Write(roomReserved)

	// Grava Matriz de Tiles (576 bytes)
	if len(room.TileMatrix) == models.RoomMatrixSize {
		_, _ = buf.Write(room.TileMatrix)
	} else {
		emptyMatrix := make([]byte, models.RoomMatrixSize)
		_, _ = buf.Write(emptyMatrix)
	}

	// Grava Entidades (8 slots fixos de 8 bytes = 64 bytes)
	for i := 0; i < MaxEntitiesPerRoom; i++ {
		if i < entCount {
			e := entities[i]
			sprID := byte(0xFF)
			if e.SpriteID != nil {
				sprID = byte(*e.SpriteID)
			}
			scriptID := uint16(0xFFFF)
			if e.EventScriptID != nil {
				scriptID = uint16(*e.EventScriptID)
			}

			_ = buf.WriteByte(byte(e.PosX))
			_ = buf.WriteByte(byte(e.PosY))
			_ = buf.WriteByte(sprID)
			_ = buf.WriteByte(byte(e.BehaviorType))
			_ = binary.Write(buf, binary.LittleEndian, scriptID)
			_ = buf.WriteByte(0)
			_ = buf.WriteByte(0)
		} else {
			// Slot vazio preenchido com zeros e sprite 0xFF
			emptyEnt := []byte{0, 0, 0xFF, 0, 0xFF, 0xFF, 0, 0}
			_, _ = buf.Write(emptyEnt)
		}
	}

	return buf.Bytes(), nil
}

func exportStrings(list []*models.StringEntry) []byte {
	buf := new(bytes.Buffer)
	count := uint16(len(list))

	// Número de strings
	_ = binary.Write(buf, binary.LittleEndian, count)

	// Tabela temporária de strings em bytes para calcular os offsets relativos
	var stringBytes [][]byte
	for _, entry := range list {
		// Null-terminated string para C no SDCC
		b := append([]byte(entry.TextContent), 0)
		stringBytes = append(stringBytes, b)
	}

	// Calcula os offsets relativos ao início do bloco de strings
	// O bloco de strings começa logo após: 2 bytes (count) + (count * 2 bytes de offset)
	currentOffset := uint16(0)
	for _, sb := range stringBytes {
		_ = binary.Write(buf, binary.LittleEndian, currentOffset)
		currentOffset += uint16(len(sb))
	}

	// Anexa o conteúdo de todas as strings
	for _, sb := range stringBytes {
		buf.Write(sb)
	}

	return buf.Bytes()
}
