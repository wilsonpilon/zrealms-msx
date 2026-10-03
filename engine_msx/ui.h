// ____________________________
// Z-Realm (zrealm-msx) - Interface de Usuário, HUD e Diálogos (MSX 2)
// Subfase 4.4: Caixa de Diálogos & HUD
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "vdp.h"

// Índices de Tiles no Banco 2 da SCREEN 4 (Linhas 16 a 23)
// Tiles 0..127: Reservados para a sala (linhas 16 e 17)
// Tiles 128..223: Fonte ASCII (código ASCII 32 a 127 -> tile = c + 96)
#define UI_FONT_TILE_OFFSET    96

// Tiles de Moldura e Ícones na faixa 224..239
#define UI_TILE_FRAME_TL       224 // Canto superior esquerdo ┌
#define UI_TILE_FRAME_TR       225 // Canto superior direito ┐
#define UI_TILE_FRAME_BL       226 // Canto inferior esquerdo └
#define UI_TILE_FRAME_BR       227 // Canto inferior direito ┘
#define UI_TILE_FRAME_H        228 // Barra horizontal ─
#define UI_TILE_FRAME_V        229 // Barra vertical │
#define UI_TILE_HEART          230 // Ícone de Coração (HP) ♥
#define UI_TILE_STAR           231 // Ícone de Estrela (MP) ★
#define UI_TILE_KEY            232 // Ícone de Chave 🗝
#define UI_TILE_ARROW_DOWN     233 // Indicador de próxima página ▼
#define UI_TILE_BLANK          255 // Espaço vazio preto

// Configurações da Janela de Diálogo (Linhas 19 a 23 da tela)
#define UI_DIALOG_ROW_TOP      19
#define UI_DIALOG_ROW_TEXT0    20
#define UI_DIALOG_ROW_TEXT1    21
#define UI_DIALOG_ROW_TEXT2    22
#define UI_DIALOG_ROW_BOTTOM   23
#define UI_DIALOG_LINE_WIDTH   30
#define UI_DIALOG_MAX_LINES    3
#define UI_DIALOG_MAX_PAGES    4

// Inicializa o subsistema de UI: carrega fontes e ícones no Banco 2 do VDP
void UI_Init(void);

// Recarrega os padrões e cores da UI no Banco 2 (após troca de tileset de sala)
void UI_ReloadFont(void);

// Atualiza o display do HUD na linha 18 (HP, MP, Nível, Chaves)
void UI_UpdateHUD(void);

// Abre a caixa de diálogo com o texto formatado (suporta até 4 páginas)
void UI_ShowDialogue(const c8* message);

// Retorna se há uma caixa de diálogo aberta atualmente
bool UI_IsDialogueActive(void);

// Atualiza a lógica de diálogo (avanço de página e fechamento via ESPAÇO/Joystick)
void UI_UpdateDialogue(void);

// Fecha a caixa de diálogo e restaura o painel padrão
void UI_CloseDialogue(void);
