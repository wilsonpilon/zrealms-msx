// _____________________________________________________________________________
//
//  UI Module - 3-Field Screen Layout (Capítulo 2)
//  Screen 0 (40x24) Text Interface for MSX
// _____________________________________________________________________________

#ifndef UI_H
#define UI_H

#include "game_types.h"

// Inicializa o modo de vídeo (Screen 0, 40 colunas), limpa VRAM e desenha as divisórias.
void UI_Init(void);

// Desenha ou redesenha a moldura com os 3 campos e as duas barras horizontais.
void UI_DrawLayout(void);

// Define e imprime o texto do Campo Superior (Linha 0).
void UI_SetTopText(const char* text);

// Limpa todo o Campo Central (Linhas 2 a 20) e reposiciona o cursor no início (2, 0).
void UI_ClearCenter(void);

// Imprime uma string no Campo Central com quebra automática de linha e rolagem se necessário.
void UI_PrintCenter(const char* text);

// Imprime um caractere único no Campo Central.
void UI_PrintCharCenter(char c);

// Imprime uma nova linha no Campo Central.
void UI_NewLineCenter(void);

// Imprime número decimal no Campo Central.
void UI_PrintNumberCenter(u8 value);

// Lê uma linha de comando do jogador no Campo Inferior (Linha 22..23) com prompt e conversão para maiúsculas.
void UI_ReadLine(char* buffer, u8 max_len);

// Pausa a execução por 'seconds' segundos.
void UI_PauseSeconds(u8 seconds);

// Aguarda o pressionamento de uma tecla (e sua posterior liberação).
void UI_WaitKey(void);

// Tabela de 10 atalhos ativos para Shift+0..9
extern u8 g_ShiftShortcuts[10];

// Exibe a janela de seleção rápida de caracteres acentuados (tecla TAB)
u8 UI_SelectAccentedChar(void);

// Exibe a janela de configuração dos atalhos Shift+0..9 (tecla SELECT)
void UI_ConfigShiftShortcuts(void);

#endif // UI_H
