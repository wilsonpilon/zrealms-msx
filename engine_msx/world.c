// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Salas, Colisão e Mundo (MSX 2)
//─────────────────────────────────────────────────────────────────────────────
#include "world.h"
#include "entity.h"

WorldState g_World;

void WORLD_Init(void)
{
	u16 i;
	g_World.CurrentRoomID = 0xFFFF;
	g_World.ActiveTilesetID = 0xFF;

	for (i = 0; i < VIEWPORT_TILE_COUNT; i++)
	{
		g_World.TileMatrix[i] = 0;
	}

	for (i = 0; i < 256; i++)
	{
		g_World.CollisionTable[i] = COLLISION_PASSABLE;
	}
}

bool WORLD_LoadRoom(u16 roomID)
{
	BinaryRoom* room;
	u16 i;

	// 1. Mapeia a sala solicitada do Memory Mapper
	room = LOADER_GetRoom(roomID);
	if (!room)
	{
		return FALSE;
	}

	g_World.CurrentRoomID = roomID;
	g_World.Header = room->Header;

	// 2. Copia a matriz de 576 tiles para RAM local segura
	for (i = 0; i < VIEWPORT_TILE_COUNT; i++)
	{
		g_World.TileMatrix[i] = room->TileMatrix[i];
	}

	// 3. Se a sala utilizar um tileset diferente do ativo, recarrega padrões, cores e colisões
	if (g_World.Header.TilesetID != g_World.ActiveTilesetID || g_World.ActiveTilesetID == 0xFF)
	{
		BinaryTileset* ts = LOADER_GetTileset(g_World.Header.TilesetID);
		if (ts)
		{
			g_World.ActiveTilesetID = g_World.Header.TilesetID;
			for (i = 0; i < 256; i++)
			{
				g_World.CollisionTable[i] = ts->CollisionTable[i];
			}
			VDP_LoadTilesetAllBanks(ts->PatternTable, ts->ColorTable);
		}
	}

	// 4. Desenha imediatamente os 576 tiles da sala no viewport do VDP
	VDP_DrawRoomViewport(g_World.TileMatrix);

	// 5. Carrega e posiciona até 8 entidades ativas da sala
	ENTITY_LoadRoomEntities(g_World.Header.EntityCount, room->Entities);
	return TRUE;
}

u8 WORLD_GetTile(u8 x, u8 y)
{
	u16 offset;
	if (x >= VIEWPORT_WIDTH || y >= VIEWPORT_HEIGHT)
	{
		return 0;
	}

	offset = ((u16)y * VIEWPORT_WIDTH) + x;
	return g_World.TileMatrix[offset];
}

u8 WORLD_GetCollision(u8 x, u8 y)
{
	u8 tileID;
	if (x >= VIEWPORT_WIDTH || y >= VIEWPORT_HEIGHT)
	{
		return COLLISION_SOLID;
	}

	tileID = WORLD_GetTile(x, y);
	return g_World.CollisionTable[tileID];
}

bool WORLD_CheckRoomTransition(i8 targetTileX, i8 targetTileY, u8* outHeroX, u8* outHeroY)
{
	// Saída ao Norte (Y < 0)
	if (targetTileY < 0)
	{
		if (g_World.Header.NorthRoom != 0xFFFF)
		{
			if (WORLD_LoadRoom(g_World.Header.NorthRoom))
			{
				*outHeroX = (u8)targetTileX;
				*outHeroY = VIEWPORT_HEIGHT - 1; // 17 (borda Sul da sala norte)
				return TRUE;
			}
		}
		return FALSE;
	}

	// Saída ao Sul (Y >= 18)
	if (targetTileY >= VIEWPORT_HEIGHT)
	{
		if (g_World.Header.SouthRoom != 0xFFFF)
		{
			if (WORLD_LoadRoom(g_World.Header.SouthRoom))
			{
				*outHeroX = (u8)targetTileX;
				*outHeroY = 0; // borda Norte da sala sul
				return TRUE;
			}
		}
		return FALSE;
	}

	// Saída a Oeste (X < 0)
	if (targetTileX < 0)
	{
		if (g_World.Header.WestRoom != 0xFFFF)
		{
			if (WORLD_LoadRoom(g_World.Header.WestRoom))
			{
				*outHeroX = VIEWPORT_WIDTH - 1; // 31 (borda Leste da sala oeste)
				*outHeroY = (u8)targetTileY;
				return TRUE;
			}
		}
		return FALSE;
	}

	// Saída a Leste (X >= 32)
	if (targetTileX >= VIEWPORT_WIDTH)
	{
		if (g_World.Header.EastRoom != 0xFFFF)
		{
			if (WORLD_LoadRoom(g_World.Header.EastRoom))
			{
				*outHeroX = 0; // borda Oeste da sala leste
				*outHeroY = (u8)targetTileY;
				return TRUE;
			}
		}
		return FALSE;
	}

	return FALSE;
}
