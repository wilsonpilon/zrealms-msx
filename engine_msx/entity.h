// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Entidades e Atores da Sala (MSX 2)
// Subfase 4.2: Spawn Dinâmico, IA Simples, Colisão e Interação
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "loader.h"
#include "vdp_screen4.h"

// Capacidade máxima de entidades ativas simultâneas por sala
#define MAX_ACTIVE_ENTITIES 8

// Tipos de Comportamento de Entidades (conforme models.BehaviorType no Z-Realm)
#define BEHAVIOR_STATIC_NPC    0 // NPC estacionário (interage por tecla de ação)
#define BEHAVIOR_WANDERING_NPC 1 // NPC que caminha aleatoriamente pelo grid
#define BEHAVIOR_PATROL_NPC    2 // NPC com rota de patrulha simples
#define BEHAVIOR_CHEST         3 // Baú com itens
#define BEHAVIOR_DOOR          4 // Porta ou passagem com fechadura
#define BEHAVIOR_TRIGGER       5 // Gatilho invisível acionado por aproximação
#define BEHAVIOR_HOSTILE       6 // Inimigo simples

// Instância de entidade ativa em tempo de execução
typedef struct
{
	u8   SlotID;        // 0..7
	u8   TileX;         // Coordenada Grid X (0..31)
	u8   TileY;         // Coordenada Grid Y (0..17)
	u8   PixelX;        // Coordenada física horizontal na tela (0..248)
	u8   PixelY;        // Coordenada física vertical na tela (0..136)
	u8   SpriteID;      // ID do sprite de recurso (0..31, 0xFF se sem sprite)
	u8   VdpSpriteSlot; // Slot de hardware no VDP (1..8)
	u8   BehaviorType;  // BEHAVIOR_*
	u16  EventScriptID; // Script vinculado
	u8   Direction;     // Direção atual (0: UP, 1: RIGHT, 2: DOWN, 3: LEFT)
	u8   Timer;         // Temporizador de IA / cooldown de passos
	bool Active;        // Entidade ativa no slot
	bool Interacted;    // Flag se já foi interagida
	u8   State;         // Estado (0: padrão/fechado, 1: aberto/acionado)
} EntityInstance;

extern EntityInstance g_Entities[MAX_ACTIVE_ENTITIES];
extern u8 g_ActiveEntityCount;

// Inicializa o subsistema de entidades
void ENTITY_Init(void);

// Desativa todas as entidades e oculta os sprites 1..8 no VDP
void ENTITY_Cleanup(void);

// Carrega as entidades da sala ativa e inicializa seus sprites no VDP
void ENTITY_LoadRoomEntities(u8 count, const BinaryEntity* entities);

// Atualiza ciclo de vida e IAs de todas as entidades ativas a cada quadro
void ENTITY_Update(void);

// Busca entidade ativa na coordenada de grid (X, Y). Retorna NULL se vazio.
EntityInstance* ENTITY_GetAt(u8 tileX, u8 tileY);

// Checa se a coordenada contém uma entidade sólida bloqueante para movimento
bool ENTITY_IsSolidAt(u8 tileX, u8 tileY);

// Executa interação de ação em um tile específico (Espaço / Gatilho do Joystick)
// Retorna TRUE se interagiu com alguma entidade
bool ENTITY_InteractAt(u8 tileX, u8 tileY);

// Checa se o herói pisou em um gatilho de piso (BEHAVIOR_TRIGGER)
void ENTITY_CheckStepTrigger(u8 tileX, u8 tileY);
