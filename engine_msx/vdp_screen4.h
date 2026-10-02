// ____________________________
// Z-Realm (zrealm-msx) - Controlador de Vídeo V9938 SCREEN 4 (Graphic 3)
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "vdp.h"

// Endereços na VRAM para SCREEN 4 (Graphic 3 no V9938)
#define VRAM_PATTERN_TABLE   0x0000 // 3 bancos de 2 KB cada (0000h, 0800h, 1000h)
#define VRAM_LAYOUT_TABLE    0x1800 // Tabela de Nomes (32 x 24 = 768 bytes)
#define VRAM_COLOR_TABLE     0x2000 // 3 bancos de 2 KB cada (2000h, 2800h, 3000h)
#define VRAM_SPRITE_COLOR    0x1C00 // Tabela de Cores dos Sprites Modo 2
#define VRAM_SPRITE_ATTR     0x1E00 // Tabela de Atributos dos Sprites
#define VRAM_SPRITE_PATTERN  0x3800 // Tabela de Padrões dos Sprites (16x16)

// Dimensões do Viewport da Sala
#define VIEWPORT_WIDTH       32
#define VIEWPORT_HEIGHT      18
#define VIEWPORT_TILE_COUNT  (VIEWPORT_WIDTH * VIEWPORT_HEIGHT) // 576 tiles

// Inicializa o modo de vídeo SCREEN 4 (Graphic 3) com suporte a Sprites Modo 2
void VDP_InitScreen4(void);

// Copia tabelas de padrões e cores para todos os 3 bancos da SCREEN 4
void VDP_LoadTilesetAllBanks(const u8* patternData, const u8* colorData);

// Desenha a matriz da sala (576 bytes) na área de jogo (linhas 0 a 17)
void VDP_DrawRoomViewport(const u8* tileMatrix);

// Preenche a barra de HUD e linha de diálogo com tiles vazios ou limpos
void VDP_ClearHUDAndDialogue(u8 clearTileIndex);

// Define um caractere na tabela de nomes em coordenadas específicas (X: 0..31, Y: 0..23)
void VDP_PrintTile(u8 x, u8 y, u8 tileIndex);
