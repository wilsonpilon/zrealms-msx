// ____________________________
// Z-Realm (zrealm-msx) - Binary Disk Loader para MSX-DOS 2
// Carrega HEADER.BIN e GAME.DAT diretamente para o Memory Mapper (Página 2)
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "dos.h"
#include "mapper.h"

// Constantes do formato binário Z-Realm
#define LOADER_MAGIC_0      'Z'
#define LOADER_MAGIC_1      'R'
#define LOADER_MAGIC_2      '0'
#define LOADER_MAGIC_3      '1'
#define LOADER_VERSION      1

#define RES_TYPE_TILESET    1
#define RES_TYPE_SPRITE     2
#define RES_TYPE_ROOM       3
#define RES_TYPE_STRINGS    4
#define RES_TYPE_GAMEDATA   5

#define MAX_RESOURCES_DIR   64

// MasterHeader: 32 bytes do HEADER.BIN
typedef struct
{
	c8  Magic[4];          // 'ZR01'
	u16 Version;           // 1
	u16 SegmentCount;      // Total de segmentos de 16 KB no GAME.DAT
	u16 ResourceCount;     // Total de recursos no diretório
	u16 InitialRoomID;     // ID da sala de partida
	u8  InitialHeroX;      // Posição inicial X do herói (0..31)
	u8  InitialHeroY;      // Posição inicial Y do herói (0..17)
	u8  InitialTileset;    // ID do tileset da sala inicial
	u8  Reserved[17];
} MasterHeader;

// ResourceEntry: 8 bytes por registro no diretório
typedef struct
{
	u8  Type;              // ResType (1=Tileset, 2=Sprite, 3=Room, etc.)
	u8  Segment;           // Segmento de 16 KB no GAME.DAT (0..N)
	u16 ID;                // ID lógico do recurso
	u16 Offset;            // Offset relativo na Página 2 (0x0000..0x3FFF)
	u16 Size;              // Tamanho em bytes
} ResourceEntry;

// BinaryRoomHeader: 16 bytes no início de cada sala
typedef struct
{
	u16 NorthRoom;         // ID da sala ao Norte (0xFFFF se sem saída)
	u16 SouthRoom;         // ID da sala ao Sul
	u16 EastRoom;          // ID da sala a Leste
	u16 WestRoom;          // ID da sala a Oeste
	u8  TilesetID;         // ID do tileset
	u8  EntityCount;       // Quantidade de entidades (0..8)
	u8  Reserved[6];
} BinaryRoomHeader;

// BinaryEntity: 8 bytes por entidade na sala
typedef struct
{
	u8  PosX;              // Coordenada X (0..31)
	u8  PosY;              // Coordenada Y (0..17)
	u8  SpriteID;          // ID do sprite (0..31, 0xFF se invisível)
	u8  BehaviorType;      // Comportamento
	u16 EventScriptID;     // Script vinculado
	u8  Padding[2];
} BinaryEntity;

// BinaryRoom: Estrutura completa de uma sala residente no segmento (656 bytes)
typedef struct
{
	BinaryRoomHeader Header;
	u8               TileMatrix[576]; // 32 x 18 tiles
	BinaryEntity     Entities[8];     // Até 8 entidades
} BinaryRoom;

// BinaryTileset: Estrutura de um tileset compilado (4608 bytes)
typedef struct
{
	u8 PatternTable[2048];   // 256 tiles * 8 bytes
	u8 ColorTable[2048];     // 256 tiles * 8 bytes
	u8 CollisionTable[256];  // 256 tiles * 1 byte
	u8 AnimationTable[256];  // 256 tiles * 1 byte
} BinaryTileset;

// BinarySprite: Estrutura de um sprite compilado Modo 2 (48 bytes: 32 bytes pattern + 16 bytes color)
typedef struct
{
	u8 Pattern[32]; // 32 bytes de padrões 16x16 (4 blocos 8x8)
	u8 Color[16];   // 16 bytes de atributos de cor por scanline
} BinarySprite;

// Funções públicas do Loader
bool LOADER_LoadGame(const c8* headerPath, const c8* dataPath);
const MasterHeader* LOADER_GetMasterHeader(void);
const ResourceEntry* LOADER_FindResource(u8 type, u16 id);
void* LOADER_MapResource(u8 type, u16 id);
BinaryRoom* LOADER_GetRoom(u16 roomID);
BinaryTileset* LOADER_GetTileset(u16 tilesetID);
BinarySprite* LOADER_GetSprite(u16 spriteID);
