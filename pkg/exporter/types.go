package exporter

// Constantes fundamentais da serialização binária para MSX 2 (V9938 / MSX-DOS 2)
const (
	// Tamanho de um segmento do Memory Mapper na Página 2 (8000h - BFFFh)
	SegmentSize = 16384 // 16 KB (0x4000)

	// Magic Number do cabeçalho mestre (4 bytes: 'Z', 'R', '0', '1')
	HeaderMagic = "ZR01"

	// Versão atual do formato binário Z-Realm
	FormatVersion = 1

	// Tipos de Recursos no Diretório de Alocação
	ResTypeTileset   uint8 = 1
	ResTypeSprite    uint8 = 2
	ResTypeRoom      uint8 = 3
	ResTypeStrings   uint8 = 4
	ResTypeGameData  uint8 = 5
	ResTypeScript    uint8 = 6

	// Limites da arquitetura MSX
	MaxEntitiesPerRoom = 8
	MaxTilesPerTileset = 256
	MaxSprites         = 32
)

// MasterHeader representa o cabeçalho inicial de 32 bytes gravado em HEADER.BIN.
type MasterHeader struct {
	Magic           [4]byte  // 'ZR01'
	Version         uint16   // Versão do formato (1)
	SegmentCount    uint16   // Total de segmentos de 16 KB contidos em GAME.DAT
	ResourceCount   uint16   // Total de recursos catalogados no diretório
	InitialRoomID   uint16   // ID da sala onde o jogador inicia
	InitialHeroX    uint8    // Posição X inicial do herói (0..31)
	InitialHeroY    uint8    // Posição Y inicial do herói (0..17)
	InitialTileset  uint8    // Tileset da sala inicial
	Reserved        [17]byte // Preenchimento para alinhamento a 32 bytes
}

// ResourceEntry representa um registro no diretório de recursos da tabela mestra (8 bytes por entrada).
type ResourceEntry struct {
	Type     uint8  // ResType (Tileset, Sprite, Room, Strings, etc.)
	Segment  uint8  // Índice do segmento de 16 KB (0..N)
	ID       uint16 // ID lógico do recurso
	Offset   uint16 // Offset relativo dentro da Página 2 (0x0000 a 0x3FFF)
	Size     uint16 // Tamanho do recurso em bytes
}

// BinaryRoomHeader define o cabeçalho binário compacto de uma sala (16 bytes).
type BinaryRoomHeader struct {
	NorthRoom   uint16  // ID da sala ao Norte (0xFFFF se bloqueado/inexistente)
	SouthRoom   uint16  // ID da sala ao Sul
	EastRoom    uint16  // ID da sala a Leste
	WestRoom    uint16  // ID da sala a Oeste
	TilesetID   uint8   // ID do tileset usado nesta sala
	EntityCount uint8   // Quantidade de entidades presentes (0..8)
	Reserved    [6]byte // Alinhamento a 16 bytes
}

// BinaryEntity representa a estrutura binária de uma entidade posicionada na sala (8 bytes).
type BinaryEntity struct {
	PosX         uint8   // Coordenada X no grid (0..31)
	PosY         uint8   // Coordenada Y no grid (0..17)
	SpriteID     uint8   // ID do sprite (0..31, 0xFF se invisível)
	BehaviorType uint8   // Tipo de comportamento / IA
	EventScriptID uint16 // ID do script acionado
	Padding      [2]byte // Alinhamento a 8 bytes
}

// BinaryTileset representa um tileset compilado pronto para VRAM do V9938 SCREEN 4 (4608 bytes).
// - PatternTable: 256 tiles * 8 bytes = 2048 bytes
// - ColorTable: 256 tiles * 8 bytes = 2048 bytes
// - CollisionTable: 256 tiles * 1 byte = 256 bytes
// - AnimationTable: 256 tiles * 1 byte = 256 bytes
type BinaryTileset struct {
	PatternTable   [2048]byte
	ColorTable     [2048]byte
	CollisionTable [256]byte
	AnimationTable [256]byte
}

// ExportResult contém informações do resultado do processo de exportação.
type ExportResult struct {
	HeaderPath     string
	DataPath       string
	TotalSegments  int
	ResourceCount  int
	TotalDataBytes int
}
