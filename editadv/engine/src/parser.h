// _____________________________________________________________________________
//
//  Parser Module (Capítulo 4, 5 e 7)
//  Syntactic Analyzer for Player Sentences: VERB + OBJ1 + OBJ2
// _____________________________________________________________________________

#ifndef PARSER_H
#define PARSER_H

#include "game_types.h"

// Analisa a frase do jogador e preenche os registradores de parsing do Game_State:
// - state->verbo_atual (1..200)
// - state->obj_buffer[0] (1..99 ou 0 se ausente)
// - state->obj_buffer[1] (1..99 ou 0 se ausente)
//
// Retorna TRUE se pelo menos o verbo foi identificado com sucesso.
// Se parsed_echo != NULL, preenche com a forma canônica da frase interpretada
// para exibição no Campo Superior (Capítulo 2).
bool Parser_Parse(const char* input, const Game_Database* db, Game_State* state, char* parsed_echo, u8 echo_max_len);

// Helper para verificar se uma palavra é um artigo ou conector ignorável
bool Parser_IsNoiseWord(const char* word);

#endif // PARSER_H
