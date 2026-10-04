// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Entidades e Atores da Sala (MSX 2)
// Subfase 4.2: Spawn Dinâmico, IA Simples, Colisão e Interação
//─────────────────────────────────────────────────────────────────────────────
#include "entity.h"
#include "world.h"
#include "hero.h"
#include "vm.h"
#include "ui.h"

EntityInstance g_Entities[MAX_ACTIVE_ENTITIES];
u8 g_ActiveEntityCount = 0;

// Gerador de números pseudo-aleatórios leve para Z80
static u8 s_RndSeed = 0x5A;

static u8 RandomByte(void)
{
	s_RndSeed = (s_RndSeed * 37) + 11;
	return s_RndSeed;
}

void ENTITY_Init(void)
{
	u8 i;
	g_ActiveEntityCount = 0;

	for (i = 0; i < MAX_ACTIVE_ENTITIES; i++)
	{
		g_Entities[i].SlotID = i;
		g_Entities[i].VdpSpriteSlot = i + 1; // Sprites 1..8 no VDP
		g_Entities[i].Active = FALSE;
		g_Entities[i].Interacted = FALSE;
		g_Entities[i].State = 0;
	}
}

void ENTITY_Cleanup(void)
{
	u8 i;
	for (i = 0; i < MAX_ACTIVE_ENTITIES; i++)
	{
		g_Entities[i].Active = FALSE;
		g_Entities[i].Interacted = FALSE;
		g_Entities[i].State = 0;
		VDP_Screen4_HideSprite(g_Entities[i].VdpSpriteSlot);
	}
	g_ActiveEntityCount = 0;
}

void ENTITY_LoadRoomEntities(u8 count, const BinaryEntity* entities)
{
	u8 i;

	// 1. Limpa todas as entidades ativas da sala anterior
	ENTITY_Cleanup();

	if (count > MAX_ACTIVE_ENTITIES)
	{
		count = MAX_ACTIVE_ENTITIES;
	}

	// 2. Carrega cada entidade válida presente no banco da sala
	for (i = 0; i < count; i++)
	{
		const BinaryEntity* src = &entities[i];
		EntityInstance* ent = &g_Entities[i];

		// Slot vazio se não possui sprite e behavior for 0
		if (src->SpriteID == 0xFF && src->BehaviorType == 0)
		{
			continue;
		}

		ent->SlotID = i;
		ent->TileX = src->PosX;
		ent->TileY = src->PosY;
		ent->PixelX = src->PosX * 8;
		ent->PixelY = src->PosY * 8;
		ent->SpriteID = src->SpriteID;
		ent->VdpSpriteSlot = i + 1; // Slot VDP 1 a 8
		ent->BehaviorType = src->BehaviorType;
		ent->EventScriptID = src->EventScriptID;
		ent->Direction = 2; // DOWN por padrão
		ent->Timer = 30 + (i * 15); // Desfasamento de temporização entre entidades
		ent->Active = TRUE;
		ent->Interacted = FALSE;
		ent->State = 0;

		// 3. Se possuir sprite, carrega o padrão/cor no VDP e posiciona na tela
		if (src->SpriteID != 0xFF)
		{
			BinarySprite* spr = LOADER_GetSprite(src->SpriteID);
			if (spr)
			{
				VDP_LoadSprite(ent->VdpSpriteSlot, spr->Pattern, spr->Color);
				VDP_SetSpritePos(ent->VdpSpriteSlot, ent->PixelX, ent->PixelY);
			}
			else
			{
				VDP_Screen4_HideSprite(ent->VdpSpriteSlot);
			}
		}
		else
		{
			// Entidade sem sprite visual (gatilho invisível, etc.)
			VDP_Screen4_HideSprite(ent->VdpSpriteSlot);
		}
	}

	g_ActiveEntityCount = count;
}

EntityInstance* ENTITY_GetAt(u8 tileX, u8 tileY)
{
	u8 i;
	for (i = 0; i < MAX_ACTIVE_ENTITIES; i++)
	{
		if (g_Entities[i].Active && g_Entities[i].TileX == tileX && g_Entities[i].TileY == tileY)
		{
			return &g_Entities[i];
		}
	}
	return NULL;
}

bool ENTITY_IsSolidAt(u8 tileX, u8 tileY)
{
	EntityInstance* ent = ENTITY_GetAt(tileX, tileY);
	if (!ent)
	{
		return FALSE;
	}

	// Gatilhos invisíveis são passáveis
	if (ent->BehaviorType == BEHAVIOR_TRIGGER)
	{
		return FALSE;
	}

	// Baús, NPCs, portas e monstros bloqueiam o deslocamento
	return TRUE;
}

bool ENTITY_InteractAt(u8 tileX, u8 tileY)
{
	EntityInstance* ent = ENTITY_GetAt(tileX, tileY);
	if (!ent)
	{
		return FALSE;
	}

	ent->Interacted = TRUE;

	switch (ent->BehaviorType)
	{
		case BEHAVIOR_STATIC_NPC:
		case BEHAVIOR_WANDERING_NPC:
		case BEHAVIOR_PATROL_NPC:
			// NPC pausa a caminhada temporariamente ao ser abordado pelo herói
			ent->Timer = 120; // 2 segundos pausado
			break;

		case BEHAVIOR_CHEST:
			// Alterna estado do baú para aberto
			if (ent->State == 0)
			{
				ent->State = 1;
			}
			break;

		case BEHAVIOR_DOOR:
			ent->State = 1;
			break;

		case BEHAVIOR_HOSTILE:
			// Herói ataca e derrota a criatura hostil com a tecla de Ação
			ent->Active = FALSE;
			VDP_Screen4_HideSprite(ent->VdpSpriteSlot);
			VM_PlaySFX(1); // Som de vitória / abate
			break;

		default:
			break;
	}

	// Executa o script de evento vinculado à entidade
	if (ent->EventScriptID != 0 && ent->EventScriptID != 0xFFFF)
	{
		VM_ExecuteScript(ent->EventScriptID);
	}

	return TRUE;
}

void ENTITY_CheckStepTrigger(u8 tileX, u8 tileY)
{
	EntityInstance* ent = ENTITY_GetAt(tileX, tileY);
	if (ent && ent->BehaviorType == BEHAVIOR_TRIGGER)
	{
		ent->Interacted = TRUE;
		ent->State = 1;

		if (ent->EventScriptID != 0 && ent->EventScriptID != 0xFFFF)
		{
			VM_ExecuteScript(ent->EventScriptID);
		}
	}
}

void ENTITY_Update(void)
{
	u8 i;

	for (i = 0; i < MAX_ACTIVE_ENTITIES; i++)
	{
		EntityInstance* ent = &g_Entities[i];

		if (!ent->Active)
		{
			continue;
		}

		switch (ent->BehaviorType)
		{
			case BEHAVIOR_WANDERING_NPC:
			{
				if (ent->Timer > 0)
				{
					ent->Timer--;
				}
				else
				{
					// Recarrega temporizador (~0.75 a 1.75 segundos a 60 Hz)
					ent->Timer = 45 + (RandomByte() & 0x3F);

					// Sorteia direção de movimento (0: Cima, 1: Direita, 2: Baixo, 3: Esquerda)
					{
						u8 dir = RandomByte() & 0x03;
						i8 dx = 0;
						i8 dy = 0;
						i8 targetX, targetY;

						if (dir == 0)      dy = -1; // Cima
						else if (dir == 1) dx = 1;  // Direita
						else if (dir == 2) dy = 1;  // Baixo
						else if (dir == 3) dx = -1; // Esquerda

						targetX = (i8)ent->TileX + dx;
						targetY = (i8)ent->TileY + dy;

						// Valida se o passo permanece dentro da área jogável
						if (targetX >= 1 && targetX < (VIEWPORT_WIDTH - 1) && targetY >= 1 && targetY < (VIEWPORT_HEIGHT - 1))
						{
							// Checa se o tile do mundo é passável
							if (WORLD_GetCollision((u8)targetX, (u8)targetY) == COLLISION_PASSABLE)
							{
								// Não colide com a posição do herói
								if (!(targetX == (i8)HERO_GetTileX() && targetY == (i8)HERO_GetTileY()))
								{
									// Não colide com outras entidades ativas
									bool blocked = FALSE;
									u8 j;
									for (j = 0; j < MAX_ACTIVE_ENTITIES; j++)
									{
										if (j != i && g_Entities[j].Active &&
											g_Entities[j].TileX == (u8)targetX &&
											g_Entities[j].TileY == (u8)targetY &&
											g_Entities[j].BehaviorType != BEHAVIOR_TRIGGER)
										{
											blocked = TRUE;
											break;
										}
									}

									if (!blocked)
									{
										ent->TileX = (u8)targetX;
										ent->TileY = (u8)targetY;
										ent->PixelX = ent->TileX * 8;
										ent->PixelY = ent->TileY * 8;
										ent->Direction = dir;

										if (ent->SpriteID != 0xFF)
										{
											VDP_SetSpritePos(ent->VdpSpriteSlot, ent->PixelX, ent->PixelY);
										}
									}
								}
							}
						}
					}
				}
				break;
			}

			case BEHAVIOR_PATROL_NPC:
			{
				if (ent->Timer > 0)
				{
					ent->Timer--;
				}
				else
				{
					ent->Timer = 35; // Passo cadenciado de patrulha

					// Movimento simples de vaivém na direção atual
					{
						i8 dx = (ent->Direction == 1) ? 1 : ((ent->Direction == 3) ? -1 : 0);
						i8 dy = (ent->Direction == 2) ? 1 : ((ent->Direction == 0) ? -1 : 0);
						i8 targetX = (i8)ent->TileX + dx;
						i8 targetY = (i8)ent->TileY + dy;
						bool canStep = FALSE;

						if (targetX >= 1 && targetX < (VIEWPORT_WIDTH - 1) && targetY >= 1 && targetY < (VIEWPORT_HEIGHT - 1))
						{
							if (WORLD_GetCollision((u8)targetX, (u8)targetY) == COLLISION_PASSABLE &&
								!(targetX == (i8)HERO_GetTileX() && targetY == (i8)HERO_GetTileY()) &&
								!ENTITY_IsSolidAt((u8)targetX, (u8)targetY))
							{
								canStep = TRUE;
							}
						}

						if (canStep)
						{
							ent->TileX = (u8)targetX;
							ent->TileY = (u8)targetY;
							ent->PixelX = ent->TileX * 8;
							ent->PixelY = ent->TileY * 8;
							if (ent->SpriteID != 0xFF)
							{
								VDP_SetSpritePos(ent->VdpSpriteSlot, ent->PixelX, ent->PixelY);
							}
						}
						else
						{
							// Inverte direção se atingiu obstáculo
							if (ent->Direction == 1) ent->Direction = 3;
							else if (ent->Direction == 3) ent->Direction = 1;
							else if (ent->Direction == 0) ent->Direction = 2;
							else ent->Direction = 0;
						}
					}
				}
				break;
			}

			case BEHAVIOR_HOSTILE:
			{
				if (ent->Timer > 0)
				{
					ent->Timer--;
				}
				else
				{
					i8 hx = (i8)HERO_GetTileX();
					i8 hy = (i8)HERO_GetTileY();
					i8 ex = (i8)ent->TileX;
					i8 ey = (i8)ent->TileY;
					i8 dx = (hx > ex) ? (hx - ex) : (ex - hx);
					i8 dy = (hy > ey) ? (hy - ey) : (ey - hy);
					i8 dist = dx + dy;

					ent->Timer = 35; // Cadência de ação do monstro

					// Ataque do monstro: Se o herói estiver adjacente (distância Manhattan == 1)
					if (dist == 1)
					{
						if (g_HeroStats.HP > 8)
						{
							g_HeroStats.HP -= 8;
						}
						else
						{
							g_HeroStats.HP = 1;
						}
						UI_UpdateHUD();
						VM_PlaySFX(2); // Dano de combate
					}
					// Perseguição: Se o herói estiver dentro do raio de visão (distância <= 7)
					else if (dist <= 7)
					{
						i8 stepX = 0;
						i8 stepY = 0;

						if (dx >= dy)
						{
							stepX = (hx > ex) ? 1 : -1;
						}
						else
						{
							stepY = (hy > ey) ? 1 : -1;
						}

						i8 targetX = ex + stepX;
						i8 targetY = ey + stepY;

						if (targetX >= 1 && targetX < (VIEWPORT_WIDTH - 1) &&
						    targetY >= 1 && targetY < (VIEWPORT_HEIGHT - 1))
						{
							if (WORLD_GetCollision((u8)targetX, (u8)targetY) == COLLISION_PASSABLE &&
							    !(targetX == hx && targetY == hy) &&
							    !ENTITY_IsSolidAt((u8)targetX, (u8)targetY))
							{
								ent->TileX = (u8)targetX;
								ent->TileY = (u8)targetY;
								ent->PixelX = ent->TileX * 8;
								ent->PixelY = ent->TileY * 8;
								if (ent->SpriteID != 0xFF)
								{
									VDP_SetSpritePos(ent->VdpSpriteSlot, ent->PixelX, ent->PixelY);
								}
							}
						}
					}
				}
				break;
			}

			default:
				break;
		}
	}
}
