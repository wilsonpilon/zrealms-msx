// ____________________________
// Z-Realm (zrealm-msx) - Interface de Usuário, HUD e Diálogos (MSX 2)
// Subfase 4.4: Caixa de Diálogos & HUD
//─────────────────────────────────────────────────────────────────────────────
#include "ui.h"
#include "vdp_screen4.h"
#include "hero.h"
#include "vm.h"
#include "keyboard.h"
#include "joystick.h"
#include "dos.h"
#include "font/font_tsm9900.h"

// Padrões dos tiles de moldura e ícones (8 bytes por tile, 10 tiles no total)
// Mapeados nos índices 224..233 do Banco 2
static const u8 s_UIPatterns[10][8] =
{
	// 224: UI_TILE_FRAME_TL ┌
	{ 0x00, 0x00, 0x00, 0x1F, 0x18, 0x18, 0x18, 0x18 },
	// 225: UI_TILE_FRAME_TR ┐
	{ 0x00, 0x00, 0x00, 0xF8, 0x18, 0x18, 0x18, 0x18 },
	// 226: UI_TILE_FRAME_BL └
	{ 0x18, 0x18, 0x18, 0x1F, 0x00, 0x00, 0x00, 0x00 },
	// 227: UI_TILE_FRAME_BR ┘
	{ 0x18, 0x18, 0x18, 0xF8, 0x00, 0x00, 0x00, 0x00 },
	// 228: UI_TILE_FRAME_H ─
	{ 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00 },
	// 229: UI_TILE_FRAME_V │
	{ 0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18 },
	// 230: UI_TILE_HEART ♥
	{ 0x00, 0x66, 0xFF, 0xFF, 0xFF, 0x7E, 0x3C, 0x18 },
	// 231: UI_TILE_STAR ★
	{ 0x18, 0x18, 0x3C, 0x7E, 0xFF, 0x7E, 0x3C, 0x18 },
	// 232: UI_TILE_KEY 🗝
	{ 0x00, 0x38, 0x44, 0x38, 0x10, 0x18, 0x10, 0x18 },
	// 233: UI_TILE_ARROW_DOWN ▼
	{ 0x00, 0x00, 0xFF, 0x7E, 0x3C, 0x18, 0x00, 0x00 }
};

// Cores dos tiles de UI (1 atributo por linha de varredura, Fg/Bg)
static const u8 s_UIColors[10] =
{
	0x71, // 224: TL - Ciano sobre Preto
	0x71, // 225: TR - Ciano sobre Preto
	0x71, // 226: BL - Ciano sobre Preto
	0x71, // 227: BR - Ciano sobre Preto
	0x71, // 228: H  - Ciano sobre Preto
	0x71, // 229: V  - Ciano sobre Preto
	0x91, // 230: Coração - Vermelho claro sobre Preto
	0x51, // 231: Estrela - Azul claro sobre Preto
	0xA1, // 232: Chave - Amarelo escuro sobre Preto
	0xB1  // 233: Seta - Amarelo claro sobre Preto
};

// Estado da caixa de diálogo
static bool s_DialogueActive = FALSE;
static u8   s_TotalPages = 0;
static u8   s_CurrentPage = 0;
static c8   s_DialogueLines[UI_DIALOG_MAX_PAGES][UI_DIALOG_MAX_LINES][UI_DIALOG_LINE_WIDTH + 1];
static bool s_DialogueDebounce = FALSE;

// Protótipos internos
static void UI_DrawStandbyPanel(void);
static void UI_DrawCurrentPage(void);
static void Format3Digits(c8* dest, u16 val);
static void Format2Digits(c8* dest, u8 val);

void UI_Init(void)
{
__asm
	.globl _g_HeroStats
	.globl _g_VMHasMessage
__endasm;

	UI_ReloadFont();
	UI_UpdateHUD();
	UI_DrawStandbyPanel();
}

void UI_ReloadFont(void)
{
	u16 i;
	u8 blank[8];

	for (i = 0; i < 8; i++)
	{
		blank[i] = 0x00;
	}

	// 1. Tile 128: Espaço (ASCII 32)
	VDP_WriteVRAM_16K(blank, 0x1400, 8);
	VDP_FillVRAM_16K(0xF1, 0x3400, 8); // Branco sobre Preto

	// 2. Tiles 129..222: Glifos ASCII 33 (!) a 126 (~) (94 caracteres = 752 bytes)
	VDP_WriteVRAM_16K(&g_Font_TMS9900[4], 0x1408, 752);
	VDP_FillVRAM_16K(0xF1, 0x3408, 752); // Branco sobre Preto

	// 3. Tile 223: DEL / Vazio
	VDP_WriteVRAM_16K(blank, 0x16F8, 8);
	VDP_FillVRAM_16K(0xF1, 0x36F8, 8);

	// 4. Tiles 224..233: Molduras e Ícones de UI (10 tiles = 80 bytes)
	for (i = 0; i < 10; i++)
	{
		u16 patAddr = 0x1700 + (i * 8);
		u16 colAddr = 0x3700 + (i * 8);
		VDP_WriteVRAM_16K(s_UIPatterns[i], patAddr, 8);
		VDP_FillVRAM_16K(s_UIColors[i], colAddr, 8);
	}

	// 5. Tile 255: Espaço vazio completamente preto
	VDP_WriteVRAM_16K(blank, 0x17F8, 8);
	VDP_FillVRAM_16K(0x11, 0x37F8, 8); // Preto sobre Preto
}

static void UI_PrintChar(u8 x, u8 y, c8 c)
{
	u8 tile;
	if (c >= 32 && c <= 126)
	{
		tile = (u8)c + UI_FONT_TILE_OFFSET;
	}
	else
	{
		tile = UI_TILE_BLANK;
	}
	VDP_PrintTile(x, y, tile);
}

static void UI_PrintString(u8 x, u8 y, const c8* str)
{
	while (*str && x < 32)
	{
		UI_PrintChar(x, y, *str);
		x++;
		str++;
	}
}

static void Format3Digits(c8* dest, u16 val)
{
	u8 d = '0';
	while (val >= 100)
	{
		val -= 100;
		d++;
	}
	dest[0] = (c8)d;

	d = '0';
	while (val >= 10)
	{
		val -= 10;
		d++;
	}
	dest[1] = (c8)d;
	dest[2] = (c8)('0' + val);
}

static void Format2Digits(c8* dest, u8 val)
{
	u8 d = '0';
	while (val >= 10)
	{
		val -= 10;
		d++;
	}
	dest[0] = (c8)d;
	dest[1] = (c8)('0' + val);
}

void UI_UpdateHUD(void)
{
	c8 numBuf[4];
	u8 keyCount;

	// Col 0..8: ♥ HP: 100/100
	VDP_PrintTile(0, 18, UI_TILE_HEART);
	UI_PrintChar(1, 18, ' ');
	Format3Digits(numBuf, g_HeroStats.HP);
	numBuf[3] = '\0';
	UI_PrintString(2, 18, numBuf);
	UI_PrintChar(5, 18, '/');
	Format3Digits(numBuf, g_HeroStats.MaxHP);
	numBuf[3] = '\0';
	UI_PrintString(6, 18, numBuf);

	// Col 9..10: Espaços
	UI_PrintChar(9, 18, ' ');
	UI_PrintChar(10, 18, ' ');

	// Col 11..19: ★ MP: 050/050
	VDP_PrintTile(11, 18, UI_TILE_STAR);
	UI_PrintChar(12, 18, ' ');
	Format3Digits(numBuf, g_HeroStats.MP);
	numBuf[3] = '\0';
	UI_PrintString(13, 18, numBuf);
	UI_PrintChar(16, 18, '/');
	Format3Digits(numBuf, g_HeroStats.MaxMP);
	numBuf[3] = '\0';
	UI_PrintString(17, 18, numBuf);

	// Col 20..21: Espaços
	UI_PrintChar(20, 18, ' ');
	UI_PrintChar(21, 18, ' ');

	// Col 22..26: LV:01
	UI_PrintString(22, 18, "LV:");
	Format2Digits(numBuf, g_HeroStats.Level);
	numBuf[2] = '\0';
	UI_PrintString(25, 18, numBuf);

	// Col 27: Espaço
	UI_PrintChar(27, 18, ' ');

	// Col 28..31: 🗝:1 
	VDP_PrintTile(28, 18, UI_TILE_KEY);
	UI_PrintChar(29, 18, ':');
	keyCount = VM_GetItemQuantity(1); // Item 1 = Chave
	if (keyCount > 9) keyCount = 9;
	UI_PrintChar(30, 18, (c8)('0' + keyCount));
	UI_PrintChar(31, 18, ' ');
}

static void UI_DrawStandbyPanel(void)
{
	u8 x;

	// Linha 19: Moldura superior
	VDP_PrintTile(0, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_TL);
	for (x = 1; x < 31; x++)
	{
		VDP_PrintTile(x, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_H);
	}
	VDP_PrintTile(31, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_TR);

	// Linha 20: Painel de repouso (Linha 1)
	VDP_PrintTile(0, UI_DIALOG_ROW_TEXT0, UI_TILE_FRAME_V);
	UI_PrintString(1, UI_DIALOG_ROW_TEXT0, "                              ");
	VDP_PrintTile(31, UI_DIALOG_ROW_TEXT0, UI_TILE_FRAME_V);

	// Linha 21: Painel de repouso (Linha 2 - Título)
	VDP_PrintTile(0, UI_DIALOG_ROW_TEXT1, UI_TILE_FRAME_V);
	UI_PrintString(1, UI_DIALOG_ROW_TEXT1, "    Z-REALM: CATACUMBAS       ");
	VDP_PrintTile(31, UI_DIALOG_ROW_TEXT1, UI_TILE_FRAME_V);

	// Linha 22: Painel de repouso (Linha 3 - Dica de Ação)
	VDP_PrintTile(0, UI_DIALOG_ROW_TEXT2, UI_TILE_FRAME_V);
	UI_PrintString(1, UI_DIALOG_ROW_TEXT2, "    [ESPACO] Interagir / Acao ");
	VDP_PrintTile(31, UI_DIALOG_ROW_TEXT2, UI_TILE_FRAME_V);

	// Linha 23: Moldura inferior
	VDP_PrintTile(0, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_BL);
	for (x = 1; x < 31; x++)
	{
		VDP_PrintTile(x, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_H);
	}
	VDP_PrintTile(31, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_BR);
}

void UI_ShowDialogue(const c8* message)
{
	u8 page, line, col, w;
	u16 i;

	if (!message || message[0] == '\0')
	{
		return;
	}

	// 1. Limpa o buffer de páginas
	for (page = 0; page < UI_DIALOG_MAX_PAGES; page++)
	{
		for (line = 0; line < UI_DIALOG_MAX_LINES; line++)
		{
			s_DialogueLines[page][line][0] = '\0';
		}
	}

	// 2. Quebra o texto respeitando palavras (word-wrap até 30 colunas por linha)
	page = 0;
	line = 0;
	col = 0;
	i = 0;

	while (message[i] != '\0' && page < UI_DIALOG_MAX_PAGES)
	{
		u16 wordStart = i;
		u8 wordLen;

		while (message[i] != ' ' && message[i] != '\0')
		{
			i++;
		}
		wordLen = (u8)(i - wordStart);

		// Se a palavra ultrapassar a linha, avança para a próxima linha/página
		if (col > 0 && (col + wordLen > UI_DIALOG_LINE_WIDTH))
		{
			s_DialogueLines[page][line][col] = '\0';
			line++;
			col = 0;
			if (line >= UI_DIALOG_MAX_LINES)
			{
				line = 0;
				page++;
				if (page >= UI_DIALOG_MAX_PAGES)
				{
					break;
				}
			}
		}

		for (w = 0; w < wordLen && col < UI_DIALOG_LINE_WIDTH; w++)
		{
			s_DialogueLines[page][line][col++] = message[wordStart + w];
		}
		s_DialogueLines[page][line][col] = '\0';

		if (message[i] == ' ')
		{
			if (col < UI_DIALOG_LINE_WIDTH)
			{
				s_DialogueLines[page][line][col++] = ' ';
				s_DialogueLines[page][line][col] = '\0';
			}
			i++;
		}
	}

	s_TotalPages = (line > 0 || col > 0) ? (page + 1) : page;
	if (s_TotalPages == 0)
	{
		s_TotalPages = 1;
	}

	s_CurrentPage = 0;
	s_DialogueActive = TRUE;
	s_DialogueDebounce = TRUE; // Evita que a tecla que acionou feche imediatamente

	// 3. Desenha a primeira página na tela
	UI_DrawCurrentPage();
}

static void UI_DrawCurrentPage(void)
{
	u8 line, col, len, x;

	// Linha 19: Moldura superior
	VDP_PrintTile(0, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_TL);
	for (x = 1; x < 31; x++)
	{
		VDP_PrintTile(x, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_H);
	}
	VDP_PrintTile(31, UI_DIALOG_ROW_TOP, UI_TILE_FRAME_TR);

	// Linhas 20 a 22: Linhas de texto
	for (line = 0; line < UI_DIALOG_MAX_LINES; line++)
	{
		u8 row = UI_DIALOG_ROW_TEXT0 + line;
		const c8* text = s_DialogueLines[s_CurrentPage][line];
		len = 0;

		VDP_PrintTile(0, row, UI_TILE_FRAME_V);

		// Imprime texto da linha
		while (text[len] != '\0' && len < UI_DIALOG_LINE_WIDTH)
		{
			UI_PrintChar(1 + len, row, text[len]);
			len++;
		}

		// Preenche restante com espaços
		for (col = len; col < UI_DIALOG_LINE_WIDTH; col++)
		{
			UI_PrintChar(1 + col, row, ' ');
		}

		VDP_PrintTile(31, row, UI_TILE_FRAME_V);
	}

	// Linha 23: Moldura inferior com indicador de avanço/fechamento
	VDP_PrintTile(0, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_BL);

	if (s_CurrentPage + 1 < s_TotalPages)
	{
		// Mais páginas: "─────────────── [ESPACO] ▼"
		for (x = 1; x < 20; x++)
		{
			VDP_PrintTile(x, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_H);
		}
		UI_PrintString(20, UI_DIALOG_ROW_BOTTOM, "[ESPACO] ");
		VDP_PrintTile(29, UI_DIALOG_ROW_BOTTOM, UI_TILE_ARROW_DOWN);
		VDP_PrintTile(30, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_H);
	}
	else
	{
		// Última página: "─────────────────── [ESPACO]"
		for (x = 1; x < 22; x++)
		{
			VDP_PrintTile(x, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_H);
		}
		UI_PrintString(22, UI_DIALOG_ROW_BOTTOM, "[ESPACO]");
		VDP_PrintTile(30, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_H);
	}

	VDP_PrintTile(31, UI_DIALOG_ROW_BOTTOM, UI_TILE_FRAME_BR);
}

bool UI_IsDialogueActive(void)
{
	return s_DialogueActive;
}

void UI_UpdateDialogue(void)
{
	u8 joy;
	bool actionPressed;

	if (!s_DialogueActive)
	{
		return;
	}

	joy = Joystick_Read(JOY_PORT_1);
	actionPressed = Keyboard_IsKeyPressed(KEY_SPACE) || 
	                IS_JOY_PRESSED(joy, JOY_INPUT_TRIGGER_A);

	if (actionPressed)
	{
		if (!s_DialogueDebounce)
		{
			s_DialogueDebounce = TRUE;

			if (s_CurrentPage + 1 < s_TotalPages)
			{
				s_CurrentPage++;
				UI_DrawCurrentPage();
				DOS_Beep();
			}
			else
			{
				UI_CloseDialogue();
			}
		}
	}
	else
	{
		s_DialogueDebounce = FALSE;
	}
}

void UI_CloseDialogue(void)
{
	s_DialogueActive = FALSE;
	g_VMHasMessage = FALSE;
	UI_DrawStandbyPanel();
}
