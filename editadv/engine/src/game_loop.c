// _____________________________________________________________________________
//
//  Game Loop Module Implementation (Capítulo 12)
// _____________________________________________________________________________

#include "game_loop.h"
#include "ui.h"
#include "parser.h"
#include "interpreter.h"

#if defined(MSXGL)
    #include "core.h"
    #include "string.h"
#else
    #include <string.h>
    #include <stdio.h>
#endif

// -----------------------------------------------------------------------------
// Game_Init
// -----------------------------------------------------------------------------
void Game_Init(const Game_Database* db, Game_State* state)
{
    u16 i;

    // Zera todos os registradores e flags
    for (i = 0; i < REG_COUNT; i++)
    {
        state->registers[i] = 0;
    }
    state->obj_evidencia = 0;
    state->obj_buffer[0] = 0;
    state->obj_buffer[1] = 0;
    state->verbo_atual = 0;
    state->gosub_ret_func = 0;
    state->gosub_ret_pc = 0;
    state->flag_espera_cmd = FALSE;
    state->flag_fim = FALSE;
    state->flag_reiniciar = FALSE;

    if (db == NULL) return;

    // Define posição inicial
    state->registers[REG_POSICAO] = db->posicao_inicial;

    // Transfere situação inicial de cada objeto para a tabela de registradores (Reg 100 + ID)
    if (db->objetos != NULL)
    {
        for (i = 0; i < db->num_objetos; i++)
        {
            const Game_Object* obj = db->objetos[i];
            if (obj != NULL && obj->id > 0 && obj->id <= OBJ_MAX)
            {
                state->registers[REG_OBJETO_OFFSET + obj->id] = obj->situacao_inicial;

                if (obj->situacao_inicial == OBJ_SIT_CARREGADO)
                {
                    state->registers[REG_OBJETOS_CARREGADOS]++;
                }
                else if (obj->situacao_inicial == OBJ_SIT_EM_OBJ3_ABERTO ||
                         obj->situacao_inicial == OBJ_SIT_EM_OBJ3_FECHADO)
                {
                    state->registers[REG_OBJETOS_NO_OBJ3]++;
                }
            }
        }
    }

    // Configura interface gráfica
    UI_Init();
    UI_SetTopText(db->titulo);

    // Carrega atalhos SHIFT pré-calculados do banco de dados (se fornecidos)
    if (db->atalhos_shift != NULL)
    {
        for (i = 0; i < 10; i++)
        {
            if (db->atalhos_shift[i] != 0)
            {
                g_ShiftShortcuts[i] = db->atalhos_shift[i];
            }
        }
    }

    // 1. Apresentação da História (Prólogo - MSG 10 / MSG_INTRO)
    {
        const char* intro = Interpreter_GetMessageText(MSG_INTRO, db);
        if (intro != NULL && intro[0] != '\0')
        {
            UI_ClearCenter();
            UI_PrintCenter(intro);
            UI_NewLineCenter();
            UI_NewLineCenter();
            UI_PrintCenter("[ Pressione uma tecla para continuar ]");
            UI_WaitKey();
        }
    }

    // 2. Tela de Instruções e Orientações do Sistema (MSG 197 se presente no jogo)
    {
        const char* guide = Interpreter_GetMessageText(197, db);
        if (guide != NULL && guide[0] != '\0')
        {
            UI_ClearCenter();
            UI_PrintCenter(guide);
            UI_NewLineCenter();
            UI_NewLineCenter();
            UI_PrintCenter("[ Pressione uma tecla para jogar ]");
            UI_WaitKey();
        }
    }

    UI_ClearCenter();

    // 3. Executa a FUNÇÃO 1 (Reset do Jogo - Inicialização dos registradores)
    Interpreter_CallFunction(FUNC_RESET, db, state);
    state->flag_espera_cmd = FALSE;

    // 4. Descreve o local de início (com saídas visíveis)
    Interpreter_DescribeCurrentRoom(db, state);
}

// -----------------------------------------------------------------------------
// Helper: Obtém a saída de uma posição
// -----------------------------------------------------------------------------
static u8 GetRoomExit(const Game_Database* db, u8 room_id, Game_Direction dir)
{
    u8 i;
    if (db == NULL || db->posicoes == NULL) return 0;

    for (i = 0; i < db->num_posicoes; i++)
    {
        const Game_Position* pos = db->posicoes[i];
        if (pos != NULL && pos->id == room_id)
        {
            if (dir < DIR_COUNT)
            {
                return pos->saidas[dir];
            }
        }
    }
    return 0;
}

// -----------------------------------------------------------------------------
// Helper: Obtém o objeto pelo ID
// -----------------------------------------------------------------------------
static const Game_Object* GetObject(const Game_Database* db, u8 obj_id)
{
    u8 i;
    if (db == NULL || db->objetos == NULL || obj_id == 0) return NULL;

    for (i = 0; i < db->num_objetos; i++)
    {
        const Game_Object* obj = db->objetos[i];
        if (obj != NULL && obj->id == obj_id)
        {
            return obj;
        }
    }
    return NULL;
}

// -----------------------------------------------------------------------------
// Game_Run
// Executa rigorosamente o Game Loop documentado no Capítulo 12
// -----------------------------------------------------------------------------
void Game_Run(const Game_Database* db, Game_State* state)
{
    char input_line[40];
    char echo_text[40];

    Game_Init(db, state);

    while (!state->flag_fim)
    {
        if (state->flag_reiniciar)
        {
            UI_NewLineCenter();
            UI_PrintCenter("*** AVISO: REINICIANDO O JOGO ***");
            UI_NewLineCenter();
            UI_PrintCenter("[ Pressione tecla para recome\207ar ]");
            UI_WaitKey();
            Game_Init(db, state);
            continue;
        }

        // =====================================================================
        // PASSO 1: Executa a FUNÇÃO 5 (se 1ª instrução for diferente de NOP)
        // =====================================================================
        if (db != NULL && db->funcoes != NULL)
        {
            u8 i;
            for (i = 0; i < db->num_funcoes; i++)
            {
                const Game_Function* f = db->funcoes[i];
                if (f != NULL && f->id == FUNC_PRE_COMANDO)
                {
                    if (f->instruction_count > 0 && f->instructions[0].op != OP_NOP)
                    {
                        Interpreter_Execute(f->instructions, f->instruction_count, db, state);
                    }
                    break;
                }
            }
        }

        if (state->flag_fim || state->flag_reiniciar) continue;

        // =====================================================================
        // PASSO 2: Recebe a frase comando do jogador
        // =====================================================================
        UI_ReadLine(input_line, sizeof(input_line));

        // Se o jogador apenas pressionou ENTER sem digitar nada:
        // Exibe novamente a descrição do local (Capítulo 6)
        if (input_line[0] == '\0')
        {
            state->registers[REG_JOGADAS_L]++;
            if (state->registers[REG_JOGADAS_L] == 0)
            {
                state->registers[REG_JOGADAS_H]++;
            }
            UI_ClearCenter();
            Interpreter_DescribeCurrentRoom(db, state);
            continue;
        }

        // Incrementa o contador de jogadas (Reg 2 / Reg 3)
        state->registers[REG_JOGADAS_L]++;
        if (state->registers[REG_JOGADAS_L] == 0)
        {
            state->registers[REG_JOGADAS_H]++;
        }

        // =====================================================================
        // PASSO 3: Incrementa REGISTRADOR 4 (sempre que for diferente de zero)
        // =====================================================================
        if (state->registers[REG_CONTADOR_4] != 0)
        {
            state->registers[REG_CONTADOR_4]++;
        }

        // =====================================================================
        // PASSO 4: Decrementa REGISTRADOR 5 e testa se é zero
        // (Sempre que for diferente de zero. Se zerar, executa a FUNÇÃO 2)
        // =====================================================================
        if (state->registers[REG_CONTADOR_5] != 0)
        {
            state->registers[REG_CONTADOR_5]--;
            if (state->registers[REG_CONTADOR_5] == 0)
            {
                UI_NewLineCenter();
                UI_PrintCenter("*** AVISO: VOC\210 SUCUMBIU \205 SEDE! ***");
                UI_NewLineCenter();
                Interpreter_CallFunction(FUNC_TIMER_BOMBA, db, state);
                state->flag_fim = TRUE;
                continue;
            }
        }

        // Checagem de Passos no Escuro (Capítulo 3 / Capítulo 12)
        if (state->registers[REG_ILUMINACAO] != 0 && state->registers[REG_ESTADO_OBJ2] == 0)
        {
            state->registers[REG_PASSOS_ESCURO]++;
            if (state->registers[REG_PASSOS_ESCURO] >= 5)
            {
                UI_NewLineCenter();
                UI_PrintCenter("*** AVISO: MORTE NO ESCURO! ***");
                UI_NewLineCenter();
                Interpreter_CallFunction(FUNC_ESCURO, db, state);
                state->flag_fim = TRUE;
                continue;
            }
        }
        else
        {
            state->registers[REG_PASSOS_ESCURO] = 0;
        }

        // =====================================================================
        // PASSO 5: Reconhece o conteúdo da frase comando e executa o comando
        // =====================================================================
        {
            bool verb_found = Parser_Parse(input_line, db, state, echo_text, sizeof(echo_text));

            // Atualiza Campo Superior com a frase interpretada
            if (echo_text[0] != '\0')
            {
                UI_SetTopText(echo_text);
            }

            if (!verb_found)
            {
                UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_ENTENDI, db));
                UI_NewLineCenter();
                continue;
            }

            // Define o primeiro objeto no registrador de evidência por padrão
            if (state->obj_buffer[0] != 0)
            {
                state->obj_evidencia = state->obj_buffer[0];
            }

            // -----------------------------------------------------------------
            // 5.1 Busca por comando customizado na tabela
            // -----------------------------------------------------------------
            {
                bool cmd_executed = FALSE;
                u8 c;

                if (db != NULL && db->comandos != NULL)
                {
                    for (c = 0; c < db->num_comandos; c++)
                    {
                        const Game_Command* cmd = db->comandos[c];
                        if (cmd != NULL &&
                            cmd->verbo == state->verbo_atual &&
                            cmd->objeto1 == state->obj_buffer[0] &&
                            cmd->objeto2 == state->obj_buffer[1])
                        {
                            Interpreter_Execute(cmd->instructions, cmd->instruction_count, db, state);
                            cmd_executed = TRUE;
                            break;
                        }
                    }
                }

                if (cmd_executed)
                {
                    continue;
                }
            }

            // -----------------------------------------------------------------
            // 5.2 Se for verbo direcional (Norte, Sul, Leste, Oeste, NE, NO, SE, SO)
            // -----------------------------------------------------------------
            if ((state->verbo_atual >= VERBO_NORTE && state->verbo_atual <= VERBO_OESTE) ||
                (state->verbo_atual >= VERBO_NORDESTE && state->verbo_atual <= VERBO_SUDOESTE))
            {
                Game_Direction dir;
                u8 saida;

                if (state->verbo_atual <= VERBO_OESTE)
                {
                    dir = (Game_Direction)(state->verbo_atual - VERBO_NORTE);
                }
                else
                {
                    dir = (Game_Direction)(DIR_NORDESTE + (state->verbo_atual - VERBO_NORDESTE));
                }

                saida = GetRoomExit(db, state->registers[REG_POSICAO], dir);

                if (saida == 0)
                {
                    // Sem saída naquela direção
                    UI_PrintCenter(Interpreter_GetMessageText(MSG_MOVIMENTO_INVALIDO, db));
                    UI_NewLineCenter();
                }
                else if (saida <= POSICAO_MAX)
                {
                    // Movimento direto
                    state->registers[REG_POSICAO] = saida;
                    UI_ClearCenter();
                    Interpreter_DescribeCurrentRoom(db, state);
                }
                else
                {
                    // Movimento Condicional (Capítulo 6: valor > 100 executa Função = valor - 100)
                    u8 func_id = saida - POSICAO_COND_OFFSET;
                    Interpreter_CallFunction(func_id, db, state);
                }
                continue;
            }

            // -----------------------------------------------------------------
            // 5.3 Verificação das Ações Padrão via Byte de Consistência
            // -----------------------------------------------------------------
            {
                u8 obj_id = state->obj_buffer[0];
                const Game_Object* obj = GetObject(db, obj_id);

                switch (state->verbo_atual)
                {
                    case VERBO_PEGUE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_PEGAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_PEGAR, db, state);
                        }
                        else if (obj == NULL)
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_ESTOU_VENDO, db));
                            UI_NewLineCenter();
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_COLOQUE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_COLOCAR_OBJ3))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_COLOCAR_OBJ3, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_TROQUE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_TROCAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_TROCAR, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_COMPRE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_COMPRAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_COMPRAR, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_ROUBE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_ROUBAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_ROUBAR, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_TIRE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_TIRAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_TIRAR, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_QUEBRE:
                        if (obj != NULL && (obj->consistencia & OBJ_CONSIST_QUEBRAR))
                        {
                            Interpreter_CallFunction(FUNC_PADRAO_QUEBRAR, db, state);
                        }
                        else
                        {
                            UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                            UI_NewLineCenter();
                        }
                        break;

                    case VERBO_SOLTE:
                        Interpreter_CallFunction(FUNC_PADRAO_SOLTAR, db, state);
                        break;

                    case VERBO_EXAMINE:
                        Interpreter_CallFunction(FUNC_PADRAO_EXAMINAR, db, state);
                        break;

                    case VERBO_PROCURE:
                        Interpreter_CallFunction(FUNC_PADRAO_PROCURAR, db, state);
                        break;

                    case VERBO_TEMOS:
                    {
                        static const Game_Instruction cmd_clist[] = { { OP_CLIST, 0, 0, 0 } };
                        Interpreter_Execute(cmd_clist, 1, db, state);
                        break;
                    }

                    case VERBO_RECOMECE:
                        state->flag_reiniciar = TRUE;
                        break;

                    default:
                        UI_PrintCenter(Interpreter_GetMessageText(MSG_NAO_POSSIVEL, db));
                        UI_NewLineCenter();
                        break;
                }
            }
        }

        // =====================================================================
        // PASSO 6: Volta ao PASSO 1 (loop automático)
        // =====================================================================
    }

    // Preserva o aviso e a causa de morte na tela central (sem limpar a tela)
    UI_NewLineCenter();
    UI_PrintCenter("========================================");
    UI_PrintCenter("Fim da partida.");
    UI_NewLineCenter();
}
