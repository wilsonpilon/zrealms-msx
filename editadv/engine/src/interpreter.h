// _____________________________________________________________________________
//
//  Interpreter Module (Capítulo 8 e 10)
//  Bytecode Interpreter for Commands and Functions
// _____________________________________________________________________________

#ifndef INTERPRETER_H
#define INTERPRETER_H

#include "game_types.h"

// Executa uma sequência de bytecodes (Game_Instruction)
void Interpreter_Execute(const Game_Instruction* instructions, u8 count, const Game_Database* db, Game_State* state);

// Executa uma Função pelo seu ID (1..200)
void Interpreter_CallFunction(u8 func_id, const Game_Database* db, Game_State* state);

// Helper para obter o texto de uma mensagem (do sistema ou do autor)
const char* Interpreter_GetMessageText(u8 msg_id, const Game_Database* db);

// Helper para obter o nome principal do objeto
const char* Interpreter_GetObjectName(u8 obj_id, const Game_Database* db);

// Helper para descrever a sala atual (equivalente a OP_DESC)
void Interpreter_DescribeCurrentRoom(const Game_Database* db, Game_State* state);

#endif // INTERPRETER_H
