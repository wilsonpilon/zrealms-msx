// _____________________________________________________________________________
//
//  UI Module Implementation - 3-Field Screen Layout (Capítulo 2)
// _____________________________________________________________________________

#include "ui.h"

#if defined(MSXGL)
    #include "msxgl.h"
    #include "font_custom.h"
#elif defined(__SDCC)
    // MSX BIOS Standard Entry Points (Z80)
    static void Bios_InitText(void) __naked
    {
    __asm
        call 0x005F // INITXT: Screen 0
        ret
    __endasm;
    }

    static void Bios_Cls(void) __naked
    {
    __asm
        call 0x00C3 // CLS
        ret
    __endasm;
    }

    static void Bios_SetCursor(u8 col, u8 row) __naked
    {
        col; row; // Suppress unreferenced warning
    __asm
        ld hl, #2
        add hl, sp
        ld a, (hl)    // col
        inc a         // BIOS POSIT usa 1-based (1..40)
        ld h, a
        inc hl
        ld a, (hl)    // row
        inc a         // BIOS POSIT usa 1-based (1..24)
        ld l, a
        call 0x00C6   // POSIT: H=col, L=row
        ret
    __endasm;
    }

    static void Bios_Chput(char c) __naked
    {
        c;
    __asm
        ld hl, #2
        add hl, sp
        ld a, (hl)
        call 0x00A2   // CHPUT
        ret
    __endasm;
    }

    static char Bios_Chget(void) __naked
    {
    __asm
        call 0x009F   // CHGET
        ld l, a
        ret
    __endasm;
    }

    static void Bios_WaitFrame(void) __naked
    {
    __asm
        ei
        halt
        ret
    __endasm;
    }
#else
    #include <stdio.h>
    #include <string.h>
    #include <ctype.h>
#endif

static u8 g_CenterCursorX = SCREEN_MARGIN_LEFT;
static u8 g_CenterCursorY = SCREEN_ROW_CENTER_START;

// -----------------------------------------------------------------------------
// UI_Init
// -----------------------------------------------------------------------------
void UI_Init(void)
{
#if defined(MSXGL)
    VDP_SetMode(VDP_MODE_SCREEN0);
    VDP_ClearVRAM();
    VDP_SetColor(0xF1);
    // Carrega o alfabeto customizado de vram.dat para a VRAM (0x0800 ~ 0x0FFF)
    VDP_WriteVRAM_16K(g_FontCustom, 0x0800, 2048);

    // Carrega o glifo customizado para Ü maiúsculo na posição 0x9F
    {
        static const u8 s_GlyphUmlautU[8] = { 0x50, 0x00, 0x88, 0x88, 0x88, 0x88, 0xF8, 0x00 };
        VDP_WriteVRAM_16K(s_GlyphUmlautU, 0x0800 + (0x9F * 8), 8);
    }

    // Configura taxa de repetição do teclado da BIOS para digitação confortável
    // REPCNT (0xF3F7): atraso inicial antes de iniciar repetição (~50 frames / ~0.8s)
    // RPTTICK (0xF3F8): intervalo de repetição (8 frames)
    *(volatile u8*)0xF3F7 = 50;
    *(volatile u8*)0xF3F8 = 8;
#elif defined(__SDCC)
    Bios_InitText();
    Bios_Cls();
    *(volatile u8*)0xF3F7 = 50;
    *(volatile u8*)0xF3F8 = 8;
#endif
    UI_DrawLayout();
    UI_ClearCenter();
}

// -----------------------------------------------------------------------------
// UI_DrawLayout
// Desenha as barras divisórias nas linhas 1 (0x1B) e 22 (0x1A)
// e os marcadores de início (0x18) nas linhas 0 e 23, coluna 2
// -----------------------------------------------------------------------------
// Macro para calcular endereço linear na VRAM da SCREEN 0 (Name Table em 0x0000)
#define VRAM_TEXT_ADDR(col, row) ((u16)((row) * SCREEN_TEXT_WIDTH + (col)))

// -----------------------------------------------------------------------------
// UI_DrawLayout
// Desenha as barras divisórias nas linhas 1 (0x1B) e 22 (0x1A)
// e os marcadores de início (0x18) nas linhas 0 e 23, coluna 2
// -----------------------------------------------------------------------------
void UI_DrawLayout(void)
{
#if defined(MSXGL)
    // Linha 0: Marcador inicial na coluna 2
    VDP_Poke_16K(SCREEN_CHAR_MARKER, VRAM_TEXT_ADDR(SCREEN_MARGIN_LEFT, SCREEN_ROW_TOP_START));

    // Linha 1: Divisória superior (0x1B das colunas 2 a 37)
    VDP_FillVRAM_16K(SCREEN_CHAR_DIVIDER_TOP, VRAM_TEXT_ADDR(SCREEN_MARGIN_LEFT, SCREEN_ROW_DIVIDER_1), SCREEN_BAR_WIDTH);

    // Linha 22: Divisória inferior (0x1A das colunas 2 a 37)
    VDP_FillVRAM_16K(SCREEN_CHAR_DIVIDER_BOTTOM, VRAM_TEXT_ADDR(SCREEN_MARGIN_LEFT, SCREEN_ROW_DIVIDER_2), SCREEN_BAR_WIDTH);

    // Linha 23: Marcador inicial na coluna 2
    VDP_Poke_16K(SCREEN_CHAR_MARKER, VRAM_TEXT_ADDR(SCREEN_MARGIN_LEFT, SCREEN_ROW_BOTTOM_START));
#elif defined(__SDCC)
    u8 x;
    Bios_SetCursor(SCREEN_MARGIN_LEFT, SCREEN_ROW_TOP_START);
    Bios_Chput(SCREEN_CHAR_MARKER);

    Bios_SetCursor(0, SCREEN_ROW_DIVIDER_1);
    Bios_Chput(' ');
    Bios_Chput(' ');
    for (x = SCREEN_MARGIN_LEFT; x <= SCREEN_MARGIN_RIGHT; x++) Bios_Chput(SCREEN_CHAR_DIVIDER_TOP);
    Bios_Chput(' ');
    Bios_Chput(' ');

    Bios_SetCursor(0, SCREEN_ROW_DIVIDER_2);
    Bios_Chput(' ');
    Bios_Chput(' ');
    for (x = SCREEN_MARGIN_LEFT; x <= SCREEN_MARGIN_RIGHT; x++) Bios_Chput(SCREEN_CHAR_DIVIDER_BOTTOM);
    Bios_Chput(' ');
    Bios_Chput(' ');

    Bios_SetCursor(SCREEN_MARGIN_LEFT, SCREEN_ROW_BOTTOM_START);
    Bios_Chput(SCREEN_CHAR_MARKER);
#endif
}

// Helper para decodificar caracteres UTF-8 de 1 ou 2 bytes para os códigos da fonte customizada (vram.dat)
static u8 UI_DecodeUTF8Char(const char** pStr)
{
    u8 c;
    if (pStr == NULL || *pStr == NULL) return 0;
    c = (u8)**pStr;
    if (c == '\0') return 0;
    (*pStr)++;
    if (c == 0xC3)
    {
        u8 c2 = (u8)**pStr;
        if (c2 == '\0') return 0;
        (*pStr)++;
        switch (c2)
        {
            // Maiúsculas acentuadas
            case 0x80: return 0x8F; // À
            case 0x81: return 0x84; // Á
            case 0x82: return 0x8C; // Â
            case 0x83: return 0xB0; // Ã
            case 0x87: return 0x80; // Ç
            case 0x88: return 0x90; // È
            case 0x89: return 0x90; // É
            case 0x8A: return 0x8D; // Ê
            case 0x8C: return 0x89; // Ì
            case 0x8D: return 0x89; // Í
            case 0x92: return 0x8A; // Ò
            case 0x93: return 0x8A; // Ó
            case 0x94: return 0x8E; // Ô
            case 0x95: return 0xB4; // Õ
            case 0x99: return 0x8B; // Ù
            case 0x9A: return 0x8B; // Ú

            // Minúsculas acentuadas
            case 0xA0: return 0x85; // à
            case 0xA1: return 0xA0; // á
            case 0xA2: return 0x83; // â
            case 0xA3: return 0xB1; // ã
            case 0xA7: return 0x87; // ç
            case 0xA8: return 0x82; // è
            case 0xA9: return 0x82; // é
            case 0xAA: return 0x88; // ê
            case 0xAC: return 0xA1; // ì
            case 0xAD: return 0xA1; // í
            case 0xB2: return 0x95; // ò
            case 0xB3: return 0xA2; // ó
            case 0xB4: return 0x93; // ô
            case 0xB5: return 0xB6; // õ (código 0xB6 em vram.dat para õ minúsculo)
            case 0xB9: return 0xA3; // ù
            case 0xBA: return 0xA3; // ú
            default: return c2;
        }
    }
    else if (c == 0xC2)
    {
        u8 c2 = (u8)**pStr;
        if (c2 == '\0') return 0;
        (*pStr)++;
        if (c2 == 0xA0) return ' '; // Non-breaking space
        if (c2 == 0xAA) return 0x9B; // ª
        if (c2 == 0xBA) return 0x9C; // º
        return c2;
    }
    if (c == '!') return 0x5B; // Ponto de exclamação na fonte customizada vram.dat (0x21 é Á)
    return c;
}

// -----------------------------------------------------------------------------
// UI_SetTopText
// Escreve o texto a partir da coluna 4 (após marcador na col 2 e espaço na col 3)
// -----------------------------------------------------------------------------
void UI_SetTopText(const char* text)
{
#if defined(MSXGL)
    u16 addr = VRAM_TEXT_ADDR(4, SCREEN_ROW_TOP_START);
    u8 max_chars = SCREEN_MARGIN_RIGHT - 4 + 1; // 38 - 4 + 1 = 35
    u8 decoded_buf[36];
    u8 len = 0;

    // Limpa colunas 4 a 39 instantaneamente
    VDP_FillVRAM_16K(' ', addr, SCREEN_TEXT_WIDTH - 4);

    if (text != NULL)
    {
        while (*text && len < max_chars)
        {
            u8 ch = UI_DecodeUTF8Char(&text);
            if (ch == '\0') break;
            decoded_buf[len++] = ch;
        }
        if (len > 0)
        {
            VDP_WriteVRAM_16K(decoded_buf, addr, len);
        }
    }
#elif defined(__SDCC)
    u8 x;
    Bios_SetCursor(4, SCREEN_ROW_TOP_START);
    for (x = 4; x < SCREEN_TEXT_WIDTH; x++) Bios_Chput(' ');

    if (text != NULL)
    {
        Bios_SetCursor(4, SCREEN_ROW_TOP_START);
        for (x = 4; *text && x <= SCREEN_MARGIN_RIGHT; x++)
        {
            u8 ch = UI_DecodeUTF8Char(&text);
            if (ch == '\0') break;
            Bios_Chput((char)ch);
        }
    }
#else
    if (text) printf("[TOPO] %s\n", text);
#endif
}

// -----------------------------------------------------------------------------
// UI_ClearCenter
// Limpa linhas 2 a 21 instantaneamente em VRAM
// -----------------------------------------------------------------------------
void UI_ClearCenter(void)
{
#if defined(MSXGL)
    // Limpa todas as 20 linhas centrais (linhas 2 a 21 = 20 * 40 = 800 bytes) com Fill rápido
    VDP_FillVRAM_16K(' ', VRAM_TEXT_ADDR(0, SCREEN_ROW_CENTER_START), (SCREEN_ROW_CENTER_END - SCREEN_ROW_CENTER_START + 1) * SCREEN_TEXT_WIDTH);
#elif defined(__SDCC)
    u8 y, x;
    for (y = SCREEN_ROW_CENTER_START; y <= SCREEN_ROW_CENTER_END; y++)
    {
        Bios_SetCursor(0, y);
        for (x = 0; x < SCREEN_TEXT_WIDTH; x++) Bios_Chput(' ');
    }
#else
    printf("\n--- [CENTRO LIMPO] ---\n");
#endif
    g_CenterCursorX = SCREEN_MARGIN_LEFT;
    g_CenterCursorY = SCREEN_ROW_CENTER_START;
}
// Helper para verificar se a tecla SHIFT está pressionada no MSX (matriz NEWKEY linha 6, bit 0)
static u8 UI_IsShiftPressed(void)
{
#if defined(MSXGL) || defined(__SDCC)
    volatile u8* pNewKeyRow6 = (volatile u8*)0xFBEB;
    return ((*pNewKeyRow6 & 0x01) == 0);
#else
    return 0;
#endif
}

// Helper para verificar se QUALQUER tecla física está pressionada (linhas 0..5 e 7..8)
static bool UI_IsAnyKeyPressed(void)
{
#if defined(MSXGL) || defined(__SDCC)
    volatile u8* pNewKey = (volatile u8*)0xFBE5;
    u8 r;
    for (r = 0; r <= 8; r++)
    {
        if (r == 6) continue; // Ignora modificadoras (SHIFT, CTRL, GRAPH, CODE)
        if (pNewKey[r] != 0xFF) return TRUE;
    }
    return FALSE;
#else
    return FALSE;
#endif
}

#if defined(MSXGL) && (TARGET_TYPE == TYPE_DOS)
static u8 DOS_CheckKey(void) __naked
{
__asm
    push ix
    ld   c, #0x0B   ; DOS_FUNC_CONST: A = 0x00 se vazio, 0xFF se caractere pronto
    call 0x0005
    pop  ix
    ld   l, a
    ret
__endasm;
}
#endif

// -----------------------------------------------------------------------------
// UI_ReadRawKey
// Lê uma tecla do teclado de forma confiável em ROM e MSX-DOS
// -----------------------------------------------------------------------------
#if defined(MSXGL) && (TARGET_TYPE == TYPE_DOS)
static u8 UI_ReadRawKey(void) __naked
{
__asm
    push ix
    ld   c, #0x07   ; MSX-DOS BDOS Direct Console Input without echo (espera tecla)
    call 0x0005     ; Chama BDOS
    pop  ix
    ld   l, a       ; Retorna caractere em L para o SDCC
    ret
__endasm;
}
#elif defined(MSXGL)
static u8 UI_ReadRawKey(void)
{
    return (u8)BIOS_GetCharacter();
}
#elif defined(__SDCC)
static u8 UI_ReadRawKey(void)
{
    return (u8)Bios_Chget();
}
#else
static u8 UI_ReadRawKey(void)
{
    return (u8)getchar();
}
#endif

// -----------------------------------------------------------------------------
// UI_NewLineCenter
// -----------------------------------------------------------------------------
void UI_NewLineCenter(void)
{
    g_CenterCursorX = SCREEN_MARGIN_LEFT;
    g_CenterCursorY++;

    if (g_CenterCursorY > SCREEN_ROW_CENTER_END)
    {
#if defined(MSXGL)
        {
            const char* msg = "-- Pressione uma tecla --";
            u8 col = SCREEN_MARGIN_LEFT;
            while (*msg) {
                VDP_Poke_16K(*msg++, VRAM_TEXT_ADDR(col++, SCREEN_ROW_CENTER_END));
            }
        }
        UI_WaitKey();
        UI_ClearCenter();
#elif defined(__SDCC)
        Bios_SetCursor(SCREEN_MARGIN_LEFT, SCREEN_ROW_CENTER_END);
        {
            const char* msg = "-- Pressione tecla --";
            while (*msg) Bios_Chput(*msg++);
        }
        UI_WaitKey();
        UI_ClearCenter();
#else
        printf("\n");
        g_CenterCursorY = SCREEN_ROW_CENTER_START;
#endif
    }
}

// -----------------------------------------------------------------------------
// UI_PrintCharCenter
// Imprime diretamente na VRAM respeitando margens 2..37 com quebra automática
// -----------------------------------------------------------------------------
void UI_PrintCharCenter(char c)
{
    if (c == '\n' || c == '\r')
    {
        UI_NewLineCenter();
        return;
    }

#if defined(MSXGL)
    VDP_Poke_16K((u8)c, VRAM_TEXT_ADDR(g_CenterCursorX, g_CenterCursorY));
#elif defined(__SDCC)
    Bios_SetCursor(g_CenterCursorX, g_CenterCursorY);
    Bios_Chput(c);
#else
    putchar(c);
#endif

    g_CenterCursorX++;
    if (g_CenterCursorX > SCREEN_MARGIN_RIGHT)
    {
        UI_NewLineCenter();
    }
}

// -----------------------------------------------------------------------------
// UI_PrintCenter
// Com suporte a decodificação UTF-8, quebra correta de palavras e escrita em bloco VRAM
// -----------------------------------------------------------------------------
void UI_PrintCenter(const char* text)
{
    u8 word_buf[48];
    if (text == NULL) return;

    while (*text)
    {
        if (*text == '\n' || *text == '\r')
        {
            UI_NewLineCenter();
            text++;
            continue;
        }
        if (*text == ' ')
        {
            UI_PrintCharCenter(' ');
            text++;
            continue;
        }

        // Decodifica a palavra inteira para word_buf respeitando UTF-8
        {
            u8 word_len = 0;
            while (*text && *text != ' ' && *text != '\n' && *text != '\r' && word_len < (sizeof(word_buf) - 1))
            {
                u8 ch = UI_DecodeUTF8Char(&text);
                if (ch == '\0') break;
                word_buf[word_len++] = ch;
            }

            // Se não cabe na linha atual mas cabe numa linha nova, quebra antes
            if (word_len <= (SCREEN_MARGIN_RIGHT - SCREEN_MARGIN_LEFT + 1) &&
                (g_CenterCursorX + word_len - 1) > SCREEN_MARGIN_RIGHT)
            {
                UI_NewLineCenter();
            }

#if defined(MSXGL)
            // Se a palavra cabe na linha, escreve ela em bloco diretamente na VRAM
            if ((g_CenterCursorX + word_len - 1) <= SCREEN_MARGIN_RIGHT)
            {
                VDP_WriteVRAM_16K(word_buf, VRAM_TEXT_ADDR(g_CenterCursorX, g_CenterCursorY), word_len);
                g_CenterCursorX += word_len;
                if (g_CenterCursorX > SCREEN_MARGIN_RIGHT)
                {
                    UI_NewLineCenter();
                }
            }
            else
            {
                u8 k;
                for (k = 0; k < word_len; k++)
                {
                    UI_PrintCharCenter((char)word_buf[k]);
                }
            }
#elif defined(__SDCC)
            {
                u8 k;
                for (k = 0; k < word_len; k++)
                {
                    UI_PrintCharCenter((char)word_buf[k]);
                }
            }
#else
            {
                u8 k;
                for (k = 0; k < word_len; k++)
                {
                    putchar(word_buf[k]);
                }
            }
#endif
        }
    }
}

// -----------------------------------------------------------------------------
// UI_PrintNumberCenter
// -----------------------------------------------------------------------------
void UI_PrintNumberCenter(u8 value)
{
    char buf[4];
    u8 idx = 0;
    if (value >= 100)
    {
        buf[idx++] = '0' + (value / 100);
        value %= 100;
        buf[idx++] = '0' + (value / 10);
        buf[idx++] = '0' + (value % 10);
    }
    else if (value >= 10)
    {
        buf[idx++] = '0' + (value / 10);
        buf[idx++] = '0' + (value % 10);
    }
    else
    {
        buf[idx++] = '0' + value;
    }
    buf[idx] = '\0';
    UI_PrintCenter(buf);
}


// Tabela de 10 atalhos padrão ativos para Shift+0..9
// Padrão de fábrica: 1:Á 2:É 3:Í 4:Ó 5:Ú 6:Ã 7:Õ 8:Ê 9:Ô e 0:Ç
u8 g_ShiftShortcuts[10] = {
    0x80, // Shift+0: Ç
    0x84, // Shift+1: Á
    0x90, // Shift+2: É
    0x89, // Shift+3: Í
    0x8A, // Shift+4: Ó
    0x8B, // Shift+5: Ú
    0xB0, // Shift+6: Ã
    0xB4, // Shift+7: Õ
    0x8D, // Shift+8: Ê
    0x8E  // Shift+9: Ô
};

// 13 caracteres acentuados em estrita ordem alfabética:
// À (0x8F), Á (0x84), Â (0x8C), Ã (0xB0), Ç (0x80), É (0x90), Ê (0x8D), Í (0x89), Ó (0x8A), Ô (0x8E), Õ (0xB4), Ú (0x8B), Ü (0x9F)
static const u8 s_AccentCodes[13] = {
    0x8F, // 0: À
    0x84, // 1: Á
    0x8C, // 2: Â
    0xB0, // 3: Ã
    0x80, // 4: Ç
    0x90, // 5: É
    0x8D, // 6: Ê
    0x89, // 7: Í
    0x8A, // 8: Ó
    0x8E, // 9: Ô
    0xB4, // 10: Õ
    0x8B, // 11: Ú
    0x9F  // 12: Ü
};

// 10 atalhos das teclas Shift+1..9 e Shift+0
static const u8 s_ShortcutSlots[10] = {
    1, 2, 3, 4, 5, 6, 7, 8, 9, 0
};
static const u8 s_ShortcutKeyNames[10] = {
    '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'
};

#if defined(MSXGL)
// Desenha a moldura de quadro usando caracteres gráficos do MSX (0x81, 0x9A, 0xA6, 0xA7, 0x5F, 0x5E)
static void UI_DrawWindowFrame(u8 x, u8 y, u8 width, u8 height)
{
    u8 i, row;
    u16 addr;

    // Borda superior: Canto Sup-Esq (0x81), Barra Horiz (0x5F), Canto Sup-Dir (0x9A)
    addr = VRAM_TEXT_ADDR(x, y);
    VDP_Poke_16K(0x81, addr++);
    for (i = 0; i < width - 2; i++) VDP_Poke_16K(0x5F, addr++);
    VDP_Poke_16K(0x9A, addr);

    // Laterais verticais (0x5E) e interior limpo
    for (row = 1; row < height - 1; row++)
    {
        addr = VRAM_TEXT_ADDR(x, y + row);
        VDP_Poke_16K(0x5E, addr++);
        for (i = 0; i < width - 2; i++) VDP_Poke_16K(' ', addr++);
        VDP_Poke_16K(0x5E, addr);
    }

    // Borda inferior: Canto Inf-Esq (0xA6), Barra Horiz (0x5F), Canto Inf-Dir (0xA7)
    addr = VRAM_TEXT_ADDR(x, y + height - 1);
    VDP_Poke_16K(0xA6, addr++);
    for (i = 0; i < width - 2; i++) VDP_Poke_16K(0x5F, addr++);
    VDP_Poke_16K(0xA7, addr);
}

// Desenha linha divisória interna conectada às bordas laterais
static void UI_DrawWindowDivider(u8 x, u8 y, u8 width)
{
    u8 i;
    u16 addr = VRAM_TEXT_ADDR(x, y);
    VDP_Poke_16K(0x5E, addr++);
    for (i = 0; i < width - 2; i++) VDP_Poke_16K(0x5F, addr++);
    VDP_Poke_16K(0x5E, addr);
}

// Imprime texto centralizado dentro de uma linha da janela
static void UI_PrintWindowText(u8 x, u8 y, u8 width, const char* text)
{
    u8 len = 0;
    const char* p = text;
    u8 col, offset;
    while (*p++) len++;
    if (len > width - 2) len = width - 2;
    offset = (width - 2 - len) / 2;

    for (col = 0; col < width - 2; col++)
    {
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 1 + col, y));
    }
    for (col = 0; col < len; col++)
    {
        VDP_Poke_16K((u8)text[col], VRAM_TEXT_ADDR(x + 1 + offset + col, y));
    }
}

// Desenha uma célula de acento na grade de 5 colunas x 3 linhas
static void UI_DrawAccentCell(u8 win_x, u8 start_y, u8 idx, bool is_selected)
{
    u8 r = idx / 5;
    u8 c = idx % 5;
    u8 x = win_x + 1 + c * 6;
    u8 y = start_y + r;
    u8 ch = s_AccentCodes[idx];

    // Célula de 6 colunas rigorosamente alinhada
    if (is_selected)
    {
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x,     y));
        VDP_Poke_16K('[', VRAM_TEXT_ADDR(x + 1, y));
        VDP_Poke_16K(ch,  VRAM_TEXT_ADDR(x + 2, y));
        VDP_Poke_16K(']', VRAM_TEXT_ADDR(x + 3, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 4, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 5, y));
    }
    else
    {
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x,     y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 1, y));
        VDP_Poke_16K(ch,  VRAM_TEXT_ADDR(x + 2, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 3, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 4, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 5, y));
    }
}

// Desenha uma célula de atalho na grade de 5 colunas x 2 linhas
// state: 0 = normal (" 1:Á  "), 1 = cursor em foco ("[1:Á] "), 2 = editando (">1:Á< ")
static void UI_DrawShortcutCell(u8 win_x, u8 start_y, u8 sidx, u8 state)
{
    u8 r = sidx / 5;
    u8 c = sidx % 5;
    u8 x = win_x + 1 + c * 6;
    u8 y = start_y + r;
    u8 slot = s_ShortcutSlots[sidx];
    u8 key_char = s_ShortcutKeyNames[sidx];
    u8 acc_char = g_ShiftShortcuts[slot];

    if (state == 1) // Cursor no atalho (Fase 1)
    {
        VDP_Poke_16K('[', VRAM_TEXT_ADDR(x,     y));
        VDP_Poke_16K(key_char, VRAM_TEXT_ADDR(x + 1, y));
        VDP_Poke_16K(':', VRAM_TEXT_ADDR(x + 2, y));
        VDP_Poke_16K(acc_char, VRAM_TEXT_ADDR(x + 3, y));
        VDP_Poke_16K(']', VRAM_TEXT_ADDR(x + 4, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 5, y));
    }
    else if (state == 2) // Atalho selecionado sendo editado (Fase 2)
    {
        VDP_Poke_16K('>', VRAM_TEXT_ADDR(x,     y));
        VDP_Poke_16K(key_char, VRAM_TEXT_ADDR(x + 1, y));
        VDP_Poke_16K(':', VRAM_TEXT_ADDR(x + 2, y));
        VDP_Poke_16K(acc_char, VRAM_TEXT_ADDR(x + 3, y));
        VDP_Poke_16K('<', VRAM_TEXT_ADDR(x + 4, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 5, y));
    }
    else // Atalho comum
    {
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x,     y));
        VDP_Poke_16K(key_char, VRAM_TEXT_ADDR(x + 1, y));
        VDP_Poke_16K(':', VRAM_TEXT_ADDR(x + 2, y));
        VDP_Poke_16K(acc_char, VRAM_TEXT_ADDR(x + 3, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 4, y));
        VDP_Poke_16K(' ', VRAM_TEXT_ADDR(x + 5, y));
    }
}
#endif

// -----------------------------------------------------------------------------
// UI_SelectAccentedChar
// Exibe janela sobreposta com as 13 letras acentuadas (tecla TAB)
// Centralizada, com moldura de quadro MSX e navegação interativa por cursor
// -----------------------------------------------------------------------------
u8 UI_SelectAccentedChar(void)
{
#if defined(MSXGL)
    static u8 s_VramBuf[400];
    u8 key, i;
    u8 cur_idx = 0;
    u8 result = 0;
    const u8 win_x = 4;
    const u8 win_y = 6;
    const u8 win_w = 32;
    const u8 win_h = 10;
    u16 base_addr = VRAM_TEXT_ADDR(0, win_y);

    // Salva 10 linhas da VRAM (linhas 6 a 15 = 400 bytes)
    VDP_ReadVRAM_16K(base_addr, s_VramBuf, 400);

    // Desenha a moldura de quadro do MSX
    UI_DrawWindowFrame(win_x, win_y, win_w, win_h);
    UI_PrintWindowText(win_x, win_y + 1, win_w, "TABELA DE ACENTOS");
    UI_DrawWindowDivider(win_x, win_y + 2, win_w);

    // Desenha todos os 13 acentos na ordem alfabética
    for (i = 0; i < 13; i++)
    {
        UI_DrawAccentCell(win_x, win_y + 3, i, (i == cur_idx));
    }

    UI_DrawWindowDivider(win_x, win_y + 6, win_w);
    UI_PrintWindowText(win_x, win_y + 7, win_w, "[Setas] Navega  [ENTER] Escolhe");
    UI_PrintWindowText(win_x, win_y + 8, win_w, "[ESC / TAB] Cancela");

    while (UI_IsAnyKeyPressed()) Halt();

    while (1)
    {
        key = UI_ReadRawKey();

        if (key == 27 || key == 9) // ESC ou TAB cancela
        {
            result = 0;
            break;
        }
        if (key == 13) // ENTER escolhe o acento atual
        {
            result = s_AccentCodes[cur_idx];
            break;
        }

        // Navegação com setas do teclado MSX
        if (key == 29) // Esquerda
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx > 0) ? cur_idx - 1 : 12;
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 28) // Direita
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx < 12) ? cur_idx + 1 : 0;
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 30) // Cima
        {
            if (cur_idx >= 5)
            {
                UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
                cur_idx -= 5;
                UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
            }
        }
        else if (key == 31) // Baixo
        {
            if (cur_idx + 5 <= 12)
            {
                UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
                cur_idx += 5;
                UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
            }
        }
        // Atalhos diretos por letra base
        else if (key == 'a' || key == 'A')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx < 3) ? cur_idx + 1 : 0;
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 'c' || key == 'C')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = 4; // Ç
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 'e' || key == 'E')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx == 5) ? 6 : 5; // É ou Ê
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 'i' || key == 'I')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = 7; // Í
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 'o' || key == 'O')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx >= 8 && cur_idx < 10) ? cur_idx + 1 : 8; // Ó, Ô, Õ
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
        else if (key == 'u' || key == 'U')
        {
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, FALSE);
            cur_idx = (cur_idx == 11) ? 12 : 11; // Ú ou Ü
            UI_DrawAccentCell(win_x, win_y + 3, cur_idx, TRUE);
        }
    }

    // Restaura a tela original
    VDP_WriteVRAM_16K(s_VramBuf, base_addr, 400);
    while (UI_IsAnyKeyPressed()) Halt();

    return result;
#else
    return 0;
#endif
}

// -----------------------------------------------------------------------------
// UI_ConfigShiftShortcuts
// Exibe janela para configurar os 10 atalhos SHIFT (tecla SELECT)
// Fluxo em 2 fases: 1º escolhe o atalho e dá ENTER -> vai pro quadro de cima,
// escolhe o acento e dá ENTER, associando imediatamente ao atalho.
// -----------------------------------------------------------------------------
void UI_ConfigShiftShortcuts(void)
{
#if defined(MSXGL)
    static u8 s_VramBuf[520];
    u8 key, i;
    u8 step = 0; // 0 = escolhendo atalho (0..9); 1 = escolhendo acento (quadro superior)
    u8 slot_idx = 0; // 0..9 (índice em s_ShortcutSlots)
    u8 accent_idx = 0; // 0..12 (índice em s_AccentCodes)
    const u8 win_x = 4;
    const u8 win_y = 4;
    const u8 win_w = 32;
    const u8 win_h = 13;
    u16 base_addr = VRAM_TEXT_ADDR(0, win_y);

    // Salva 13 linhas da VRAM (linhas 4 a 16 = 520 bytes)
    VDP_ReadVRAM_16K(base_addr, s_VramBuf, 520);

    // Desenha a moldura de quadro do MSX
    UI_DrawWindowFrame(win_x, win_y, win_w, win_h);
    UI_PrintWindowText(win_x, win_y + 1, win_w, "TABELA DE ACENTOS");

    // Desenha todos os 13 acentos (sem seleção inicial)
    for (i = 0; i < 13; i++)
    {
        UI_DrawAccentCell(win_x, win_y + 2, i, FALSE);
    }

    UI_DrawWindowDivider(win_x, win_y + 5, win_w);
    UI_PrintWindowText(win_x, win_y + 6, win_w, "ATALHOS (SHIFT + 1..9, 0)");

    // Desenha os 10 atalhos (slot 0 selecionado inicialmente)
    for (i = 0; i < 10; i++)
    {
        UI_DrawShortcutCell(win_x, win_y + 7, i, (i == slot_idx ? 1 : 0));
    }

    UI_DrawWindowDivider(win_x, win_y + 9, win_w);
    UI_PrintWindowText(win_x, win_y + 10, win_w, "Escolha o atalho e tecle ENTER");
    UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas/0-9] Navega   ESC: Sair");

    while (UI_IsAnyKeyPressed()) Halt();

    while (1)
    {
        key = UI_ReadRawKey();

        if (step == 0) // ================= FASE 1: Escolhendo atalho (0..9) =================
        {
            if (key == 27 || key == 24) // ESC ou SELECT conclui
            {
                break;
            }

            if (key == 13) // ENTER confirma atalho escolhido e sobe para o quadro de acentos!
            {
                step = 1;
                // Posiciona cursor inicial de acento no acento atual daquele atalho
                {
                    u8 cur_char = g_ShiftShortcuts[s_ShortcutSlots[slot_idx]];
                    accent_idx = 0;
                    for (i = 0; i < 13; i++)
                    {
                        if (s_AccentCodes[i] == cur_char) { accent_idx = i; break; }
                    }
                }

                // Destaca atalho como "em edição" (>X:Y<) e ativa cursor no quadro de acentos ([Z])
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 2);
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
                UI_PrintWindowText(win_x, win_y + 10, win_w, "Escolha o acento e tecle ENTER");
                UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas] Navega   ESC: Voltar");
                continue;
            }

            // Seleção direta de tecla numérica '1'..'9'
            if (key >= '1' && key <= '9')
            {
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 0);
                slot_idx = key - '1';
                step = 1;
                {
                    u8 cur_char = g_ShiftShortcuts[s_ShortcutSlots[slot_idx]];
                    accent_idx = 0;
                    for (i = 0; i < 13; i++)
                    {
                        if (s_AccentCodes[i] == cur_char) { accent_idx = i; break; }
                    }
                }
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 2);
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
                UI_PrintWindowText(win_x, win_y + 10, win_w, "Escolha o acento e tecle ENTER");
                UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas] Navega   ESC: Voltar");
                continue;
            }
            if (key == '0')
            {
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 0);
                slot_idx = 9; // Slot 0
                step = 1;
                {
                    u8 cur_char = g_ShiftShortcuts[s_ShortcutSlots[slot_idx]];
                    accent_idx = 0;
                    for (i = 0; i < 13; i++)
                    {
                        if (s_AccentCodes[i] == cur_char) { accent_idx = i; break; }
                    }
                }
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 2);
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
                UI_PrintWindowText(win_x, win_y + 10, win_w, "Escolha o acento e tecle ENTER");
                UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas] Navega   ESC: Voltar");
                continue;
            }

            // Navegação por setas nos atalhos
            if (key == 29) // Esquerda
            {
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 0);
                slot_idx = (slot_idx > 0) ? slot_idx - 1 : 9;
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 1);
            }
            else if (key == 28) // Direita
            {
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 0);
                slot_idx = (slot_idx < 9) ? slot_idx + 1 : 0;
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 1);
            }
            else if (key == 30 || key == 31) // Cima / Baixo (alterna entre linha 1..5 e 6..0)
            {
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 0);
                slot_idx = (slot_idx + 5) % 10;
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 1);
            }
        }
        else // ================= FASE 2: Escolhendo acento no quadro de cima =================
        {
            if (key == 27) // ESC volta para escolher outro atalho sem alterar
            {
                step = 0;
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 1);
                UI_PrintWindowText(win_x, win_y + 10, win_w, "Escolha o atalho e tecle ENTER");
                UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas/0-9] Navega   ESC: Sair");
                continue;
            }

            if (key == 13) // ENTER associa o acento ao atalho e ativa para o jogo!
            {
                step = 0;
                u8 target_slot = s_ShortcutSlots[slot_idx];
                g_ShiftShortcuts[target_slot] = s_AccentCodes[accent_idx];

                // Atualiza tela: desliga cursor do quadro de cima, atualiza atalho
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                UI_DrawShortcutCell(win_x, win_y + 7, slot_idx, 1);
                UI_PrintWindowText(win_x, win_y + 10, win_w, "Atalho alterado com sucesso!");
                UI_PrintWindowText(win_x, win_y + 11, win_w, "[Setas/0-9] Navega   ESC: Sair");
                continue;
            }

            // Navegação com setas no quadro de acentos
            if (key == 29) // Esquerda
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx > 0) ? accent_idx - 1 : 12;
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 28) // Direita
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx < 12) ? accent_idx + 1 : 0;
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 30) // Cima
            {
                if (accent_idx >= 5)
                {
                    UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                    accent_idx -= 5;
                    UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
                }
            }
            else if (key == 31) // Baixo
            {
                if (accent_idx + 5 <= 12)
                {
                    UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                    accent_idx += 5;
                    UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
                }
            }
            // Atalhos diretos por letra no quadro de acentos
            else if (key == 'a' || key == 'A')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx < 3) ? accent_idx + 1 : 0;
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 'c' || key == 'C')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = 4; // Ç
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 'e' || key == 'E')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx == 5) ? 6 : 5; // É ou Ê
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 'i' || key == 'I')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = 7; // Í
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 'o' || key == 'O')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx >= 8 && accent_idx < 10) ? accent_idx + 1 : 8; // Ó, Ô, Õ
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
            else if (key == 'u' || key == 'U')
            {
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, FALSE);
                accent_idx = (accent_idx == 11) ? 12 : 11; // Ú ou Ü
                UI_DrawAccentCell(win_x, win_y + 2, accent_idx, TRUE);
            }
        }
    }

    // Restaura tela original
    VDP_WriteVRAM_16K(s_VramBuf, base_addr, 520);
    while (UI_IsAnyKeyPressed()) Halt();
#endif
}

void UI_ReadLine(char* buffer, u8 max_len)
{
    u8 len = 0;
    u8 dead_key = 0;
    u16 prompt_base = VRAM_TEXT_ADDR(4, SCREEN_ROW_BOTTOM_START);

    if (max_len == 0 || buffer == NULL) return;

#if defined(MSXGL) || defined(__SDCC)
    // Limpa a área de comando da linha 23
    #if defined(MSXGL)
        VDP_Poke_16K(SCREEN_CHAR_MARKER, VRAM_TEXT_ADDR(SCREEN_MARGIN_LEFT, SCREEN_ROW_BOTTOM_START));
        VDP_FillVRAM_16K(' ', prompt_base, SCREEN_TEXT_WIDTH - 4);
        VDP_Poke_16K('_', prompt_base); // Cursor inicial
    #else
        Bios_SetCursor(SCREEN_MARGIN_LEFT, SCREEN_ROW_BOTTOM_START);
        Bios_Chput(SCREEN_CHAR_MARKER);
        Bios_Chput(' ');
        {
            u8 x;
            for (x = 4; x < SCREEN_TEXT_WIDTH; x++) Bios_Chput(' ');
        }
        Bios_SetCursor(4, SCREEN_ROW_BOTTOM_START);
    #endif

    while (TRUE)
    {
        u8 ch = UI_ReadRawKey();

        if (ch == 13 || ch == 10) // ENTER
        {
            #if defined(MSXGL)
                // Remove o cursor visual
                VDP_Poke_16K(' ', prompt_base + len);
            #endif
            break;
        }
        else if (ch == 8 || ch == 127) // Backspace / DEL
        {
            if (dead_key != 0)
            {
                dead_key = 0;
                continue;
            }
            if (len > 0)
            {
                #if defined(MSXGL)
                    VDP_Poke_16K(' ', prompt_base + len); // Apaga cursor na pos antiga
                    len--;
                    VDP_Poke_16K('_', prompt_base + len); // Desenha cursor na pos nova
                #else
                    len--;
                    Bios_SetCursor(4 + len, SCREEN_ROW_BOTTOM_START);
                    Bios_Chput(' ');
                #endif
            }
            continue;
        }
        else if (ch == 9) // TAB: Menu de seleção de caractere acentuado
        {
            u8 sel = UI_SelectAccentedChar();
            if (sel != 0 && len < (max_len - 1) && (4 + len) <= SCREEN_MARGIN_RIGHT)
            {
                buffer[len] = (char)sel;
                #if defined(MSXGL)
                    VDP_Poke_16K(sel, prompt_base + len);
                    len++;
                    if ((4 + len) <= SCREEN_MARGIN_RIGHT)
                    {
                        VDP_Poke_16K('_', prompt_base + len);
                    }
                #else
                    Bios_SetCursor(4 + len, SCREEN_ROW_BOTTOM_START);
                    Bios_Chput((char)sel);
                    len++;
                #endif
            }
            continue;
        }
        else if (ch == 24) // SELECT: Configuração dos atalhos SHIFT
        {
            UI_ConfigShiftShortcuts();
            #if defined(MSXGL)
                if ((4 + len) <= SCREEN_MARGIN_RIGHT)
                {
                    VDP_Poke_16K('_', prompt_base + len);
                }
            #endif
            continue;
        }

        // --- Suporte a Atalhos de Teclado para Ç e Acentos ---

        // 1. Tecla direta de Ç ou códigos estendidos (BIOS nacional / Latin-1 / CP437)
        if (ch == 0x80 || ch == 0x87 || ch == 0xC7 || ch == 0xE7)
        {
            ch = 0x80; // Ç
        }
        else if (ch == 0x90 || ch == 0x82 || ch == 0xC9 || ch == 0xE9)
        {
            ch = 0x90; // É
        }
        else if (ch == 0x84 || ch == 0xA0 || ch == 0xC1 || ch == 0xE1)
        {
            ch = 0x84; // Á
        }
        else if (ch == 0x89 || ch == 0xA1 || ch == 0xCD || ch == 0xED)
        {
            ch = 0x89; // Í
        }
        else if (ch == 0x8A || ch == 0xA2 || ch == 0xD3 || ch == 0xF3)
        {
            ch = 0x8A; // Ó
        }
        else if (ch == 0x8B || ch == 0xA3 || ch == 0xDA || ch == 0xFA)
        {
            ch = 0x8B; // Ú
        }
        else if (ch == 0xB0 || ch == 0xB1 || ch == 0xC3 || ch == 0xE3)
        {
            ch = 0xB0; // Ã
        }
        else if (ch == 0xB4 || ch == 0xB5 || ch == 0xD5 || ch == 0xF5)
        {
            ch = 0xB4; // Õ
        }
        else if (ch == 0x8C || ch == 0x83 || ch == 0xC2 || ch == 0xE2)
        {
            ch = 0x8C; // Â
        }
        else if (ch == 0x8E || ch == 0x93 || ch == 0xD4 || ch == 0xF4)
        {
            ch = 0x8E; // Ô
        }
        else if (ch == 0x8D || ch == 0x88 || ch == 0xCA || ch == 0xEA)
        {
            ch = 0x8D; // Ê
        }
        // 2. Atalho rápido CTRL+C -> Ç
        else if (ch == 3)
        {
            ch = 0x80; // Ç
        }
        // 3. Teclas de atalho SHIFT + NÚMEROS (Padrão MSX do Editor de Adventures):
        // SHIFT+1..9 e SHIFT+0 mapeados conforme g_ShiftShortcuts
        else if (UI_IsShiftPressed())
        {
            if (ch == '0' || ch == ')') ch = g_ShiftShortcuts[0];
            else if (ch >= '1' && ch <= '9') ch = g_ShiftShortcuts[ch - '0'];
            else if (ch == '!') ch = g_ShiftShortcuts[1];
            else if (ch == '"' || ch == '@') ch = g_ShiftShortcuts[2];
            else if (ch == '#') ch = g_ShiftShortcuts[3];
            else if (ch == '$') ch = g_ShiftShortcuts[4];
            else if (ch == '%') ch = g_ShiftShortcuts[5];
            else if (ch == '&') ch = g_ShiftShortcuts[6];
            else if (ch == '\'') ch = g_ShiftShortcuts[7];
            else if (ch == '*' || ch == '(') ch = g_ShiftShortcuts[8];
        }
        // 4. Suporte a Dead Keys (teclado internacional / emuladores):
        // ' + C -> Ç, ' + vogal -> Á/É/Í/Ó/Ú, ~ + A/O -> Ã/Õ, ^ + vogal -> Â/Ê/Ô
        if (dead_key != 0)
        {
            u8 prev_dead = dead_key;
            dead_key = 0;
            if (prev_dead == '\'')
            {
                if (ch == 'C' || ch == 'c') ch = 0x80; // ' + C -> Ç
                else if (ch == 'A' || ch == 'a') ch = 0x84; // ' + A -> Á
                else if (ch == 'E' || ch == 'e') ch = 0x90; // ' + E -> É
                else if (ch == 'I' || ch == 'i') ch = 0x89; // ' + I -> Í
                else if (ch == 'O' || ch == 'o') ch = 0x8A; // ' + O -> Ó
                else if (ch == 'U' || ch == 'u') ch = 0x8B; // ' + U -> Ú
                else if (ch == ' ' || ch == '\'') ch = '\'';
            }
            else if (prev_dead == '~')
            {
                if (ch == 'A' || ch == 'a') ch = 0xB0; // ~ + A -> Ã
                else if (ch == 'O' || ch == 'o') ch = 0xB4; // ~ + O -> Õ
                else if (ch == ' ' || ch == '~') ch = '~';
            }
            else if (prev_dead == '^')
            {
                if (ch == 'A' || ch == 'a') ch = 0x8C; // ^ + A -> Â
                else if (ch == 'E' || ch == 'e') ch = 0x8D; // ^ + E -> Ê
                else if (ch == 'O' || ch == 'o') ch = 0x8E; // ^ + O -> Ô
                else if (ch == ' ' || ch == '^') ch = '^';
            }
            else if (prev_dead == '`')
            {
                if (ch == 'A' || ch == 'a') ch = 0x8F; // ` + A -> À
                else if (ch == ' ' || ch == '`') ch = '`';
            }
        }
        else if (ch == '\'' || ch == '~' || ch == '^' || ch == '`')
        {
            dead_key = ch;
            continue;
        }

        // Insere caractere no buffer
        if (ch >= 32 && len < (max_len - 1) && (4 + len) <= SCREEN_MARGIN_RIGHT)
        {
            if (ch >= 'a' && ch <= 'z')
            {
                ch = ch - ('a' - 'A');
            }
            buffer[len] = (char)ch;
            #if defined(MSXGL)
                VDP_Poke_16K(ch, prompt_base + len);
                len++;
                if ((4 + len) <= SCREEN_MARGIN_RIGHT)
                {
                    VDP_Poke_16K('_', prompt_base + len);
                }
            #else
                Bios_SetCursor(4 + len, SCREEN_ROW_BOTTOM_START);
                Bios_Chput((char)ch);
                len++;
                while (UI_IsAnyKeyPressed())
                {
                    Bios_WaitFrame();
                }
            #endif
        }
    }
    buffer[len] = '\0';
#else
    printf("\n> ");
    if (fgets(buffer, max_len, stdin) != NULL)
    {
        char* p = buffer;
        while (*p)
        {
            if (*p == '\n' || *p == '\r') { *p = '\0'; break; }
            if (*p >= 'a' && *p <= 'z') { *p = *p - ('a' - 'A'); }
            p++;
        }
    }
#endif
}

// -----------------------------------------------------------------------------
// UI_PauseSeconds
// -----------------------------------------------------------------------------
void UI_PauseSeconds(u8 seconds)
{
#if defined(MSXGL)
    u16 i, frames = (u16)seconds * 60;
    for (i = 0; i < frames; i++) Halt();
#elif defined(__SDCC)
    u16 i, frames = (u16)seconds * 60;
    for (i = 0; i < frames; i++) Bios_WaitFrame();
#else
    (void)seconds;
#endif
}

// -----------------------------------------------------------------------------
// UI_WaitKey
// Aguarda confiavelmente que o jogador solte teclas anteriores, pressione uma
// tecla nova e solte-a, drenando buffers residuais do DOS.
// -----------------------------------------------------------------------------
void UI_WaitKey(void)
{
#if defined(MSXGL)
    // 1. Aguarda que qualquer tecla previamente pressionada (ex: ENTER inicial) seja solta
    while (UI_IsAnyKeyPressed())
    {
        Halt();
    }

    // 2. Drena caracteres residuais do buffer do MSX-DOS
    #if (TARGET_TYPE == TYPE_DOS)
    while (DOS_CheckKey() != 0)
    {
        UI_ReadRawKey();
    }
    #endif

    // 3. Aguarda que uma nova tecla seja pressionada
    UI_ReadRawKey();

    // 4. Aguarda que o jogador solte a tecla para não disparar a tela seguinte acidentalmente
    while (UI_IsAnyKeyPressed())
    {
        Halt();
    }

    // 5. Garante que o buffer do DOS fique limpo
    #if (TARGET_TYPE == TYPE_DOS)
    while (DOS_CheckKey() != 0)
    {
        UI_ReadRawKey();
    }
    #endif
#elif defined(__SDCC)
    while (UI_IsAnyKeyPressed()) Bios_WaitFrame();
    UI_ReadRawKey();
    while (UI_IsAnyKeyPressed()) Bios_WaitFrame();
#else
    getchar();
#endif
}

