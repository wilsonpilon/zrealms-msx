// ____________________________
// Z-Realm (zrealm-msx) - Controlador de Vídeo V9938 SCREEN 4 (Graphic 3)
//─────────────────────────────────────────────────────────────────────────────
#include "vdp_screen4.h"

// Inicializa a SCREEN 4 (Graphic 3 do V9938)
void VDP_InitScreen4(void)
{
	// Desliga a tela temporariamente para evitar "snow" e artefatos visuais
	VDP_EnableDisplay(FALSE);

	// Configura o modo SCREEN 4 no VDP
	VDP_SetMode(VDP_MODE_GRAPHIC3);

	// Limpa a tabela de nomes (768 bytes com 0x00)
	VDP_FillVRAM_16K(0x00, VRAM_LAYOUT_TABLE, 768);

	// Ativa sprites de 16x16 (Modo 2)
	VDP_SetSpriteFlag(VDP_SPRITE_SIZE_16);

	// Oculta todos os sprites (define Y = 216 / 0xD8 para desativar processador de sprites)
	VDP_FillVRAM_16K(216, VRAM_SPRITE_ATTR, 32 * 4);

	// Reativa o display e habilita interrupções de sincronismo vertical (V-Blank)
	VDP_EnableDisplay(TRUE);
	VDP_EnableVBlank(TRUE);
}

// Carrega o conjunto de 256 tiles para as três seções verticais da tela
void VDP_LoadTilesetAllBanks(const u8* patternData, const u8* colorData)
{
	u8 bank;

	for (bank = 0; bank < 3; bank++)
	{
		u16 patternAddr = VRAM_PATTERN_TABLE + (bank * 0x0800);
		u16 colorAddr   = VRAM_COLOR_TABLE   + (bank * 0x0800);

		// 2048 bytes de padrões por banco
		VDP_WriteVRAM_16K(patternData, patternAddr, 2048);

		// 2048 bytes de atributos de cores por banco
		VDP_WriteVRAM_16K(colorData, colorAddr, 2048);
	}
}

// Escreve os 576 bytes da matriz de tiles diretamente na Name Table do VDP
void VDP_DrawRoomViewport(const u8* tileMatrix)
{
	// Escreve 576 bytes contíguos iniciando na posição (0, 0) da Name Table (0x1800)
	VDP_WriteVRAM_16K(tileMatrix, VRAM_LAYOUT_TABLE, VIEWPORT_TILE_COUNT);
}

// Limpa as áreas de HUD (linhas 18-19) e Diálogo (linhas 20-23)
void VDP_ClearHUDAndDialogue(u8 clearTileIndex)
{
	// Linhas 18 a 23 = 6 linhas * 32 colunas = 192 tiles
	u16 hudStartAddr = VRAM_LAYOUT_TABLE + VIEWPORT_TILE_COUNT;
	VDP_FillVRAM_16K(clearTileIndex, hudStartAddr, 192);
}

// Imprime um tile avulso na tela
void VDP_PrintTile(u8 x, u8 y, u8 tileIndex)
{
	if (x >= 32 || y >= 24)
		return;

	u16 addr = VRAM_LAYOUT_TABLE + (y * 32) + x;
	VDP_Poke_16K(tileIndex, addr);
}

// Carrega padrão (32 bytes) e cores (16 bytes) para um sprite de 16x16 (Modo 2)
void VDP_LoadSprite(u8 spriteIndex, const u8* patternData, const u8* colorData)
{
	u16 patternAddr;
	u16 colorAddr;

	if (spriteIndex >= 32)
		return;

	patternAddr = VRAM_SPRITE_PATTERN + ((u16)spriteIndex * 32);
	colorAddr   = VRAM_SPRITE_COLOR   + ((u16)spriteIndex * 16);

	VDP_WriteVRAM_16K(patternData, patternAddr, 32);
	VDP_WriteVRAM_16K(colorData, colorAddr, 16);
}

// Posiciona um sprite de 16x16 na tela (coordenadas de pixels 0..255, 0..211)
void VDP_SetSpritePos(u8 spriteIndex, u8 x, u8 y)
{
	u8 attr[4];
	u16 attrAddr;

	if (spriteIndex >= 32)
		return;

	attrAddr = VRAM_SPRITE_ATTR + ((u16)spriteIndex * 4);

	// No VDP V9938, a linha de exibição na tela é Y + 1 (portanto, para exibir em y usa-se y - 1)
	attr[0] = (u8)(y - 1);
	attr[1] = x;
	attr[2] = spriteIndex * 4; // Em 16x16, cada sprite ocupa 4 padrões de 8x8
	attr[3] = 0;               // Sem flags especiais

	VDP_WriteVRAM_16K(attr, attrAddr, 4);
}

// Oculta um sprite desativando-o no VDP (Y = 216)
void VDP_Screen4_HideSprite(u8 spriteIndex)
{
	if (spriteIndex >= 32)
		return;

	VDP_Poke_16K(216, VRAM_SPRITE_ATTR + ((u16)spriteIndex * 4));
}

