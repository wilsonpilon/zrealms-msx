package models

// Constantes fundamentais do hardware MSX 2 (V9938) e da geometria de tela do Z-Realm
const (
	// Dimensões de Tiles (SCREEN 4 / Graphic 3)
	TileWidth       = 8
	TileHeight      = 8
	TilePatternSize = 8 // 8 bytes: 1 bit por pixel (8 linhas de 8 pixels)
	TileColorSize   = 8 // 8 bytes: nibble alto = cor de primeiro plano (Fg), nibble baixo = cor de fundo (Bg)

	// Dimensões de Sprites (Modo 2 do V9938: 16x16 pixels)
	SpriteWidth       = 16
	SpriteHeight      = 16
	SpritePatternSize = 32 // 32 bytes (4 blocos de 8x8 pixels dispostos em 16x16)
	SpriteColorSize   = 16 // 16 bytes (1 byte por scanline para cada uma das 16 linhas)

	// Dimensões da Grade de Salas (Viewport Z-Realm)
	RoomWidth      = 32
	RoomHeight     = 18
	RoomMatrixSize = RoomWidth * RoomHeight // 576 bytes

	// Dimensões da Tela Inteira do V9938 (SCREEN 4)
	ScreenWidth  = 32
	ScreenHeight = 24

	// Áreas da Tela (32x24 tiles)
	ViewportRows = 18 // Linhas 0 a 17: Área jogável da sala
	HUDRows      = 2  // Linhas 18 a 19: Barra de Status / HUD
	DialogueRows = 4  // Linhas 20 a 23: Janela de Diálogos / Mensagens
)

// CollisionType define a física e o comportamento de colisão de um tile.
type CollisionType int

const (
	CollisionPassable CollisionType = 0 // Terreno livre para travessia
	CollisionSolid    CollisionType = 1 // Obstáculo intransponível (paredes, rochas)
	CollisionWater    CollisionType = 2 // Água / Requer barco ou afoga
	CollisionDamage   CollisionType = 3 // Lava / Espinhos / Causa dano
	CollisionTrigger  CollisionType = 4 // Gatilho de evento ao pisar
)

// BehaviorType define a inteligência artificial ou a mecânica de uma entidade.
type BehaviorType int

const (
	BehaviorStaticNPC    BehaviorType = 0 // NPC estacionário (interage por tecla de ação)
	BehaviorWanderingNPC BehaviorType = 1 // NPC que caminha aleatoriamente pelo grid
	BehaviorPatrolNPC    BehaviorType = 2 // NPC com rota de patrulha predefinida
	BehaviorChest        BehaviorType = 3 // Baú com itens
	BehaviorDoor         BehaviorType = 4 // Porta ou passagem com fechadura
	BehaviorTrigger      BehaviorType = 5 // Gatilho invisível acionado por aproximação
	BehaviorHostile      BehaviorType = 6 // Inimigo simples que persegue o herói
)

// ItemType define a categoria de um item do RPG.
type ItemType int

const (
	ItemWeapon     ItemType = 0
	ItemArmor      ItemType = 1
	ItemConsumable ItemType = 2
	ItemKey        ItemType = 3
	ItemQuest      ItemType = 4
)
