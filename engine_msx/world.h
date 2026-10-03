// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Salas, Colisão e Mundo (MSX 2)
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "loader.h"
#include "vdp_screen4.h"

// Tipos de colisão física de tiles (conforme modelo do Z-Realm)
#define COLLISION_PASSABLE  0 // Terreno livre para travessia
#define COLLISION_SOLID     1 // Obstáculo intransponível (paredes, rochas)
#define COLLISION_WATER     2 // Água / Requer barco ou afoga
#define COLLISION_DAMAGE    3 // Dano (lava, espinhos)
#define COLLISION_TRIGGER   4 // Gatilho de evento ao pisar

// Estado da sala e do ambiente ativo
typedef struct
{
	u16              CurrentRoomID;
	BinaryRoomHeader Header;
	u8               TileMatrix[VIEWPORT_TILE_COUNT]; // 576 bytes em RAM local
	u8               CollisionTable[256];             // 256 bytes da física do tileset
	u8               ActiveTilesetID;
} WorldState;

extern WorldState g_World;

// Inicializa as estruturas globais de mundo
void WORLD_Init(void);

// Carrega uma sala pelo seu ID, atualizando tileset e VRAM se necessário
bool WORLD_LoadRoom(u16 roomID);

// Retorna o tipo de colisão do tile na coordenada de grid (X: 0..31, Y: 0..17)
u8 WORLD_GetCollision(u8 x, u8 y);

// Retorna o índice do tile na coordenada de grid
u8 WORLD_GetTile(u8 x, u8 y);

// Verifica saída pelas bordas da sala e executa transição automática para sala vizinha
// Retorna TRUE se a transição ocorreu com sucesso, atualizando outHeroX e outHeroY
bool WORLD_CheckRoomTransition(i8 targetTileX, i8 targetTileY, u8* outHeroX, u8* outHeroY);
