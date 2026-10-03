// ____________________________
// Z-Realm (zrealm-msx) - Controlador e Entidade do Herói (MSX 2)
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "keyboard.h"
#include "joystick.h"
#include "vdp_screen4.h"
#include "world.h"

// Cooldown em frames (a 50/60 Hz) entre passos consecutivos ao manter a tecla pressionada
#define HERO_STEP_COOLDOWN_FRAMES  6

typedef enum
{
	HERO_DIR_UP = 0,
	HERO_DIR_RIGHT = 1,
	HERO_DIR_DOWN = 2,
	HERO_DIR_LEFT = 3
} HeroDirection;

typedef struct
{
	u8            TileX;        // Coordenada no Grid (0..31)
	u8            TileY;        // Coordenada no Grid (0..17)
	u8            PixelX;       // Coordenada física horizontal na tela (0..248)
	u8            PixelY;       // Coordenada física vertical na tela (0..136)
	HeroDirection Direction;      // Direção atual
	u8            StepCooldown;   // Temporizador em frames para repetição de passos
	u8            ActionCooldown; // Temporizador em frames para tecla de ação (Espaço / Trigger)
	bool          Moved;          // Flag indicando movimentação recente
} Hero;

extern Hero g_Hero;

// Inicializa o herói com posição e carrega o sprite no VDP
void HERO_Init(u8 startTileX, u8 startTileY, u8 spriteResourceID);

// Atualiza leitura de controles, checagem de colisão e deslocamento no grid
void HERO_Update(void);

// Atualiza a posição do sprite do herói no VDP (Sprite 0)
void HERO_Draw(void);

// Posiciona o herói diretamente em uma coordenada do grid
void HERO_SetPosition(u8 tileX, u8 tileY);
