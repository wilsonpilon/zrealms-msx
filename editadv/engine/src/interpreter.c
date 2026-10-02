// _____________________________________________________________________________
//
//  Interpreter Module Implementation (Capítulo 8 e 10)
// _____________________________________________________________________________

#include "interpreter.h"
#include "ui.h"

static u8 s_rnd = 42;
static u8 GetRandom8(void) { s_rnd = (s_rnd * 17 + 53); return s_rnd; }
#define GET_RANDOM_8() GetRandom8()

// Mensagens padrão do sistema (Capítulo 9.2)
static const char* const g_SystemMessages[] = {
    "",                                                                         // 0
    "", "", "", "", "", "", "", "", "",                                         // 1..9
    "Bem-vindo à aventura!",                                                    // 10: MSG_INTRO
    "Achei o que você queria.",                                                 // 11: MSG_ACHEI
    "Está muito escuro aqui. É melhor arranjar alguma luz ou teremos problemas.", // 12: MSG_ESCURO
    "Perdão, não entendi...",                                                   // 13: MSG_NAO_ENTENDI
    "É impossível ir nesta direção.",                                           // 14: MSG_MOVIMENTO_INVALIDO
    "Isto não é possível.",                                                     // 15: MSG_NAO_POSSIVEL
    "Nós não temos isso.",                                                      // 16: MSG_NAO_TEMOS
    "Nós já temos isso.",                                                       // 17: MSG_JA_TEMOS
    "Eu não estou vendo isso por aqui.",                                        // 18: MSG_NAO_ESTOU_VENDO
    "É apenas um objeto comum.",                                                // 19: MSG_OBJETO_COMUM
    "Não dá para carregar mais nada.",                                          // 20: MSG_CARGA_MAXIMA
    "Não cabe mais nada dentro."                                                // 21: MSG_OBJ3_LOTADO
};

// -----------------------------------------------------------------------------
// Interpreter_GetMessageText
// -----------------------------------------------------------------------------
const char* Interpreter_GetMessageText(u8 msg_id, const Game_Database* db)
{
    u8 i;
    // Checa mensagens customizadas do autor no banco de dados primeiro
    if (db != NULL && db->mensagens != NULL)
    {
        for (i = 0; i < db->num_mensagens; i++)
        {
            if (db->mensagens[i] != NULL && db->mensagens[i]->id == msg_id)
            {
                return db->mensagens[i]->text;
            }
        }
    }
    // Caso não encontre mensagem customizada, usa as mensagens de sistema (10..21)
    if (msg_id >= 10 && msg_id <= 21)
    {
        return g_SystemMessages[msg_id];
    }

    return "";
}

// -----------------------------------------------------------------------------
// Interpreter_GetObjectName
// Retorna apenas a primeira palavra antes da barra '/'
// -----------------------------------------------------------------------------
const char* Interpreter_GetObjectName(u8 obj_id, const Game_Database* db)
{
    static char s_obj_name_buf[24];
    u8 i;

    if (obj_id == 0 || db == NULL || db->objetos == NULL) return "";

    for (i = 0; i < db->num_objetos; i++)
    {
        const Game_Object* obj = db->objetos[i];
        if (obj != NULL && obj->id == obj_id)
        {
            const char* p = obj->nome;
            u8 k = 0;
            while (*p && *p != '/' && k < (sizeof(s_obj_name_buf) - 1))
            {
                s_obj_name_buf[k++] = *p++;
            }
            s_obj_name_buf[k] = '\0';
            return s_obj_name_buf;
        }
    }
    return "";
}

// -----------------------------------------------------------------------------
// Interpreter_DescribeCurrentRoom (OP_DESC)
// -----------------------------------------------------------------------------
void Interpreter_DescribeCurrentRoom(const Game_Database* db, Game_State* state)
{
    u8 room_id = state->registers[REG_POSICAO];
    u8 i;
    const Game_Position* pos = NULL;

    // Se estiver no escuro e sem fonte de luz acesa
    if (state->registers[REG_ILUMINACAO] != 0 && state->registers[REG_ESTADO_OBJ2] == 0)
    {
        UI_PrintCenter(Interpreter_GetMessageText(MSG_ESCURO, db));
        UI_NewLineCenter();
        return;
    }

    if (db != NULL && db->posicoes != NULL)
    {
        for (i = 0; i < db->num_posicoes; i++)
        {
            if (db->posicoes[i] != NULL && db->posicoes[i]->id == room_id)
            {
                pos = db->posicoes[i];
                break;
            }
        }
    }

    if (pos != NULL && pos->desc != NULL)
    {
        UI_PrintCenter(pos->desc);
        UI_NewLineCenter();
    }

    // Saídas visíveis da posição
    if (pos != NULL)
    {
        static const char* const dir_names[DIR_COUNT] = { "NORTE", "SUL", "LESTE", "OESTE" };
        bool has_exit = FALSE;
        u8 d;

        UI_NewLineCenter();
        UI_PrintCenter("Sa\241das vis\241veis: ");
        for (d = 0; d < DIR_COUNT; d++)
        {
            if (pos->saidas[d] != 0)
            {
                if (has_exit)
                {
                    UI_PrintCenter(", ");
                }
                UI_PrintCenter(dir_names[d]);
                has_exit = TRUE;
            }
        }
        if (!has_exit)
        {
            UI_PrintCenter("NENHUMA");
        }
        UI_NewLineCenter();
    }

    // Lista sumária de objetos visíveis no local
    {
        bool first = TRUE;
        for (i = 1; i <= db->num_objetos; i++)
        {
            // O objeto 1 ("LOCAL") não deve ser listado como item no chão
            if (i == OBJ_ID_LOCAL) continue;

            if (state->registers[REG_OBJETO_OFFSET + i] == room_id)
            {
                if (first)
                {
                    UI_NewLineCenter();
                    UI_PrintCenter("Neste local tem:");
                    UI_NewLineCenter();
                    first = FALSE;
                }
                UI_PrintCenter("- ");
                UI_PrintCenter(Interpreter_GetObjectName(i, db));
                UI_NewLineCenter();
            }
        }
    }
}

// -----------------------------------------------------------------------------
// Interpreter_CallFunction
// -----------------------------------------------------------------------------
void Interpreter_CallFunction(u8 func_id, const Game_Database* db, Game_State* state)
{
    u8 i;
    if (db == NULL || db->funcoes == NULL || func_id == 0) return;

    for (i = 0; i < db->num_funcoes; i++)
    {
        const Game_Function* f = db->funcoes[i];
        if (f != NULL && f->id == func_id)
        {
            Interpreter_Execute(f->instructions, f->instruction_count, db, state);
            return;
        }
    }
}

// -----------------------------------------------------------------------------
// Interpreter_Execute
// Loop de execução das instruções de bytecode
// -----------------------------------------------------------------------------
void Interpreter_Execute(const Game_Instruction* instructions, u8 count, const Game_Database* db, Game_State* state)
{
    u8 pc = 0;

    if (instructions == NULL || count == 0) return;

    state->flag_espera_cmd = FALSE;

    while (pc < count && !state->flag_espera_cmd && !state->flag_fim && !state->flag_reiniciar)
    {
        const Game_Instruction* inst = &instructions[pc];
        u8 op = inst->op;
        u8 p1 = inst->p1;
        u8 p2 = inst->p2;
        u8 p3 = inst->p3;

        switch (op)
        {
            case OP_NOP:
                // Sem efeito
                break;

            case OP_MSG:
                UI_PrintCenter(Interpreter_GetMessageText(p1, db));
                UI_NewLineCenter();
                break;

            case OP_NVC:
                // Se a próxima instrução for de término (FIM ou NEU), apenas insere nova linha e avança
                if (pc + 1 < count && (instructions[pc + 1].op == OP_FIM || instructions[pc + 1].op == OP_NEU))
                {
                    UI_NewLineCenter();
                    break;
                }
                state->flag_espera_cmd = TRUE;
                return;

            case OP_LLIST:
            {
                u8 room_id = state->registers[REG_POSICAO];
                u8 i;
                bool found = FALSE;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    if (i == OBJ_ID_LOCAL) continue;
                    if (state->registers[REG_OBJETO_OFFSET + i] == room_id)
                    {
                        if (!found)
                        {
                            UI_PrintCenter("Neste local tem:");
                            UI_NewLineCenter();
                            found = TRUE;
                        }
                        UI_PrintCenter("- ");
                        UI_PrintCenter(Interpreter_GetObjectName(i, db));
                        UI_NewLineCenter();
                    }
                }
                if (!found)
                {
                    UI_PrintCenter("Não há nada de especial aqui.");
                    UI_NewLineCenter();
                }
                state->flag_espera_cmd = TRUE;
                return;
            }

            case OP_CLIST:
            {
                u8 i;
                bool found = FALSE;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    if (state->registers[REG_OBJETO_OFFSET + i] == OBJ_SIT_CARREGADO)
                    {
                        if (!found)
                        {
                            UI_PrintCenter("Você está carregando:");
                            UI_NewLineCenter();
                            found = TRUE;
                        }
                        UI_PrintCenter("- ");
                        UI_PrintCenter(Interpreter_GetObjectName(i, db));
                        if (i == OBJ_ID_CONTAINER && state->registers[REG_OBJETOS_NO_OBJ3] > 0)
                        {
                            UI_PrintCenter(" (com itens)");
                        }
                        UI_NewLineCenter();
                    }
                }
                if (!found)
                {
                    UI_PrintCenter("Você não está carregando nada.");
                    UI_NewLineCenter();
                }
                state->flag_espera_cmd = TRUE;
                return;
            }

            case OP_DLIST:
            {
                u8 i;
                bool found = FALSE;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    u8 sit = state->registers[REG_OBJETO_OFFSET + i];
                    if (sit == OBJ_SIT_EM_OBJ3_ABERTO || sit == OBJ_SIT_EM_OBJ3_FECHADO)
                    {
                        if (!found)
                        {
                            UI_PrintCenter("Dentro tem:");
                            UI_NewLineCenter();
                            found = TRUE;
                        }
                        UI_PrintCenter("- ");
                        UI_PrintCenter(Interpreter_GetObjectName(i, db));
                        UI_NewLineCenter();
                    }
                }
                if (!found)
                {
                    UI_PrintCenter("Está vazio.");
                    UI_NewLineCenter();
                }
                state->flag_espera_cmd = TRUE;
                return;
            }

            case OP_OBJ:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                UI_PrintCenter(Interpreter_GetObjectName(target, db));
                break;
            }

            case OP_INC:
                state->registers[p1]++;
                break;

            case OP_DEC:
                if (state->registers[p1] > 0)
                {
                    state->registers[p1]--;
                }
                break;

            case OP_LDR:
                state->registers[p1] = p2;
                break;

            case OP_SOMA:
                state->registers[p1] += p2;
                break;

            case OP_RND:
                if (p2 > 0)
                {
                    state->registers[p1] += (u8)(GET_RANDOM_8() % (p2 + 1));
                }
                break;

            case OP_REG_EQ:
                if (state->registers[p1] == p2)
                {
                    pc = p3;
                    continue;
                }
                break;

            case OP_REG_GT:
                if (state->registers[p1] > p2)
                {
                    pc = p3;
                    continue;
                }
                break;

            case OP_REG_LT:
                if (state->registers[p1] < p2)
                {
                    pc = p3;
                    continue;
                }
                break;

            case OP_AQUI:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                if (state->registers[REG_OBJETO_OFFSET + target] == state->registers[REG_POSICAO])
                {
                    pc = p2;
                    continue;
                }
                break;
            }

            case OP_LOCAL:
                if (state->registers[REG_POSICAO] == p1)
                {
                    pc = p2;
                    continue;
                }
                break;

            case OP_TEMOS:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                u8 sit = state->registers[REG_OBJETO_OFFSET + target];
                if (sit == OBJ_SIT_CARREGADO ||
                    ((sit == OBJ_SIT_EM_OBJ3_ABERTO || sit == OBJ_SIT_EM_OBJ3_FECHADO) &&
                     state->registers[REG_OBJETO_OFFSET + OBJ_ID_CONTAINER] == OBJ_SIT_CARREGADO))
                {
                    pc = p2;
                    continue;
                }
                break;
            }

            case OP_SOLTA:
            {
                if (p1 == 100) // Solta todos os objetos carregados
                {
                    u8 i;
                    for (i = 1; i <= db->num_objetos; i++)
                    {
                        if (state->registers[REG_OBJETO_OFFSET + i] == OBJ_SIT_CARREGADO)
                        {
                            state->registers[REG_OBJETO_OFFSET + i] = state->registers[REG_POSICAO];
                        }
                    }
                    state->registers[REG_OBJETOS_CARREGADOS] = 0;
                }
                else
                {
                    u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                    u8 sit = state->registers[REG_OBJETO_OFFSET + target];
                    if (sit == OBJ_SIT_CARREGADO)
                    {
                        state->registers[REG_OBJETO_OFFSET + target] = state->registers[REG_POSICAO];
                        if (state->registers[REG_OBJETOS_CARREGADOS] > 0)
                        {
                            state->registers[REG_OBJETOS_CARREGADOS]--;
                        }
                    }
                    else if (sit == OBJ_SIT_EM_OBJ3_ABERTO || sit == OBJ_SIT_EM_OBJ3_FECHADO)
                    {
                        state->registers[REG_OBJETO_OFFSET + target] = state->registers[REG_POSICAO];
                        if (state->registers[REG_OBJETOS_NO_OBJ3] > 0)
                        {
                            state->registers[REG_OBJETOS_NO_OBJ3]--;
                        }
                    }
                }
                break;
            }

            case OP_PEGA:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                u8 old_sit = state->registers[REG_OBJETO_OFFSET + target];
                if (old_sit == OBJ_SIT_EM_OBJ3_ABERTO || old_sit == OBJ_SIT_EM_OBJ3_FECHADO)
                {
                    if (state->registers[REG_OBJETOS_NO_OBJ3] > 0)
                    {
                        state->registers[REG_OBJETOS_NO_OBJ3]--;
                    }
                }
                state->registers[REG_OBJETO_OFFSET + target] = OBJ_SIT_CARREGADO;
                state->registers[REG_OBJETOS_CARREGADOS]++;
                break;
            }

            case OP_CRIA:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                state->registers[REG_OBJETO_OFFSET + target] = state->registers[REG_POSICAO];
                break;
            }

            case OP_APAG:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                if (state->registers[REG_OBJETO_OFFSET + target] == OBJ_SIT_CARREGADO)
                {
                    if (state->registers[REG_OBJETOS_CARREGADOS] > 0)
                    {
                        state->registers[REG_OBJETOS_CARREGADOS]--;
                    }
                }
                state->registers[REG_OBJETO_OFFSET + target] = OBJ_SIT_INEXISTENTE;
                break;
            }

            case OP_GOSUB:
                state->gosub_ret_pc = pc + 1;
                Interpreter_CallFunction(p1, db, state);
                break;

            case OP_LIBR:
            {
                u8 i;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    if (state->registers[REG_OBJETO_OFFSET + i] == OBJ_SIT_EM_OBJ3_FECHADO)
                    {
                        state->registers[REG_OBJETO_OFFSET + i] = OBJ_SIT_EM_OBJ3_ABERTO;
                    }
                }
                break;
            }

            case OP_TRC:
            {
                u8 i;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    if (state->registers[REG_OBJETO_OFFSET + i] == OBJ_SIT_EM_OBJ3_ABERTO)
                    {
                        state->registers[REG_OBJETO_OFFSET + i] = OBJ_SIT_EM_OBJ3_FECHADO;
                    }
                }
                break;
            }

            case OP_POE:
            {
                u8 target = (p1 != 0) ? p1 : state->obj_evidencia;
                if (state->registers[REG_OBJETO_OFFSET + target] == OBJ_SIT_CARREGADO)
                {
                    if (state->registers[REG_OBJETOS_CARREGADOS] > 0)
                    {
                        state->registers[REG_OBJETOS_CARREGADOS]--;
                    }
                }
                state->registers[REG_OBJETO_OFFSET + target] = OBJ_SIT_EM_OBJ3_ABERTO;
                state->registers[REG_OBJETOS_NO_OBJ3]++;
                break;
            }

            case OP_ESV:
            {
                u8 i;
                for (i = 1; i <= db->num_objetos; i++)
                {
                    u8 sit = state->registers[REG_OBJETO_OFFSET + i];
                    if (sit == OBJ_SIT_EM_OBJ3_ABERTO || sit == OBJ_SIT_EM_OBJ3_FECHADO)
                    {
                        state->registers[REG_OBJETO_OFFSET + i] = state->registers[REG_POSICAO];
                    }
                }
                state->registers[REG_OBJETOS_NO_OBJ3] = 0;
                break;
            }

            case OP_OK:
                UI_PrintCenter("Ok.");
                UI_NewLineCenter();
                state->flag_espera_cmd = TRUE;
                return;

            case OP_REGN:
                UI_PrintNumberCenter(state->registers[p1]);
                break;

            case OP_NVF:
                state->flag_espera_cmd = TRUE;
                return;

            case OP_REF:
                return;

            case OP_FIM:
                state->flag_fim = TRUE;
                return;

            case OP_NEU:
                state->flag_reiniciar = TRUE;
                return;

            case OP_DESC:
                Interpreter_DescribeCurrentRoom(db, state);
                state->flag_espera_cmd = TRUE;
                return;

            case OP_RET:
                return;

            case OP_GOTO:
                pc = p1;
                continue;

            case OP_PAUSA:
                if (p1 == 0)
                {
                    UI_NewLineCenter();
                    UI_PrintCenter("[ Pressione uma tecla ]");
                    UI_WaitKey();
                    UI_ClearCenter();
                }
                else
                {
                    UI_PauseSeconds(p1);
                }
                break;

            case OP_FLAG:
                state->registers[p1] = p2;
                break;

            case OP_EVID:
                if (p1 == 1) state->obj_evidencia = state->obj_buffer[0];
                else if (p1 == 2) state->obj_evidencia = state->obj_buffer[1];
                break;

            case OP_CLS:
                UI_ClearCenter();
                break;

            case OP_EVD_EQ:
                if (state->obj_evidencia == p1)
                {
                    pc = p2;
                    continue;
                }
                break;

            case OP_CHRS:
                UI_PrintCharCenter((char)p1);
                break;

            case OP_PRT:
                UI_PrintCenter(Interpreter_GetMessageText(p1, db));
                UI_PrintCenter(" ");
                UI_PrintCenter(Interpreter_GetObjectName(p2, db));
                UI_NewLineCenter();
                state->flag_espera_cmd = TRUE;
                return;

            case OP_DNT:
            {
                u8 sit = state->registers[REG_OBJETO_OFFSET + p1];
                if (sit == OBJ_SIT_EM_OBJ3_ABERTO || sit == OBJ_SIT_EM_OBJ3_FECHADO)
                {
                    pc = p2;
                    continue;
                }
                break;
            }

            case OP_CMD:
                if (p1 > 0 && p1 <= db->num_comandos && db->comandos != NULL)
                {
                    const Game_Command* cmd = db->comandos[p1 - 1];
                    if (cmd != NULL)
                    {
                        Interpreter_Execute(cmd->instructions, cmd->instruction_count, db, state);
                    }
                }
                break;

            default:
                break;
        }

        pc++;
    }
}
