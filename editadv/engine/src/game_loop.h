// _____________________________________________________________________________
//
//  Game Loop Module (Capítulo 12)
//  Main Command Loop and Game Lifecycle
// _____________________________________________________________________________

#ifndef GAME_LOOP_H
#define GAME_LOOP_H

#include "game_types.h"

// Inicializa o estado do jogo e variáveis a partir da base de dados.
void Game_Init(const Game_Database* db, Game_State* state);

// Executa o Game Loop principal de acordo com o Capítulo 12 do manual.
void Game_Run(const Game_Database* db, Game_State* state);

#endif // GAME_LOOP_H
