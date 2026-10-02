// _____________________________________________________________________________
//
//  Clean-Room Reverse-Engineered Adventure Engine for MSX (MSXgl / SDCC)
//  Based on "Editor de Adventures" by Renato Degiovani
// _____________________________________________________________________________

#ifndef GAME_TYPES_H
#define GAME_TYPES_H

#ifndef MSXGL_TYPES
#define MSXGL_TYPES
#if defined(MSXGL)
    #include "core.h"
#else
    // Standalone / SDCC / PC mock compatibility types
    typedef unsigned char       u8;
    typedef signed char         i8;
    typedef unsigned short      u16;
    typedef signed short        i16;
    typedef unsigned char       bool;
    #ifndef TRUE
        #define TRUE            (1)
        #define FALSE           (0)
    #endif
    #ifndef NULL
        #define NULL            ((void*)0)
    #endif
#endif
#endif

// =============================================================================
// 1. INTERFACE GRÁFICA & CAMPOS DA TELA (Capítulo 2)
// =============================================================================
// A tela é dividida horizontalmente em 3 campos principais:
// - CAMPO SUPERIOR: Exibe título do jogo e eco da frase interpretada do jogador.
// - CAMPO CENTRAL: Descrições de posições, objetos e mensagens de resposta.
// - CAMPO INFERIOR: Entrada interativa de comandos do jogador.
// =============================================================================
#define SCREEN_TEXT_WIDTH           40      // Modo texto Screen 0 padrão (40 colunas)
#define SCREEN_TEXT_HEIGHT          24      // 24 linhas

#define SCREEN_ROW_TOP_START        0       // Linha do Campo Superior (Título e Eco)
#define SCREEN_ROW_TOP_END          0
#define SCREEN_ROW_DIVIDER_1        1       // Barra divisória superior horizontal

#define SCREEN_ROW_CENTER_START     2       // Linha inicial do Campo Central
#define SCREEN_ROW_CENTER_END       21      // Linha final do Campo Central (linhas 2 a 21)
#define SCREEN_ROW_DIVIDER_2        22      // Barra divisória inferior horizontal

#define SCREEN_ROW_BOTTOM_START     23      // Linha do Campo Inferior (Prompt / Comandos)
#define SCREEN_ROW_BOTTOM_END       23

#define SCREEN_MARGIN_LEFT          2       // Margem esquerda da área de texto útil (coluna 2)
#define SCREEN_MARGIN_RIGHT         37      // Margem direita da área de texto útil (coluna 37)
#define SCREEN_BAR_WIDTH            36      // 36 caracteres de barra divisória (colunas 2 a 37)

#define SCREEN_CHAR_MARKER          0x18    // Marcador vertical na coluna 2 (linhas 0 e 23)
#define SCREEN_CHAR_DIVIDER_TOP     0x1B    // Meia-barra inferior usada na linha 1 (colunas 2 a 37)
#define SCREEN_CHAR_DIVIDER_BOTTOM  0x1A    // Meia-barra superior usada na linha 22 (colunas 2 a 37)


// =============================================================================
// 2. REGISTRADORES (Capítulo 3)
// =============================================================================
// O sistema dispõe de uma tabela de 256 registradores de 1 byte (0..255).
// Os registradores 1 a 14 possuem propósitos pré-definidos na arquitetura.
// Os registradores 15 a 99 são livres para uso geral pelo autor da aventura.
// Os registradores 100 a 199 mapeiam diretamente a situação dos objetos (100 + ID).
// =============================================================================
#define REG_COUNT                   256

#define REG_POSICAO                 1       // Posição atual do jogador (1..99)
#define REG_JOGADAS_L               2       // Contador de jogadas (byte menos significativo) - inc a cada ENTER
#define REG_JOGADAS_H               3       // Contador de jogadas (byte mais significativo)
#define REG_CONTADOR_4              4       // Contador especial: incrementado a cada frase se != 0
#define REG_CONTADOR_5              5       // Contador especial / bomba: decrementado a cada frase se != 0.
                                            // Quando atinge zero, executa a FUNÇÃO 2 automaticamente.
#define REG_PASSOS_ESCURO           6       // Passos no escuro: incrementado a cada comando no escuro.
                                            // Ao atingir o valor 5, executa a FUNÇÃO 3 automaticamente.
#define REG_OBJETOS_NO_OBJ3         7       // Quantidade de objetos dentro do Objeto 3 (normalmente máx 3)
#define REG_OBJETOS_CARREGADOS      8       // Quantidade de objetos sendo carregados pelo jogador (normalmente máx 5)
#define REG_ILUMINACAO              9       // Flag de iluminação do local: 0 = claro, 1 = escuro
#define REG_ESTADO_OBJ2             10      // Estado do Objeto 2 (fonte de luz): 0 = apagado, 1 = aceso
#define REG_RELOGIO_MIN             11      // Minutos do relógio de jogo
#define REG_RELOGIO_HORA            12      // Horas do relógio de jogo
#define REG_RELOGIO_DIA             13      // Dias do relógio de jogo
#define REG_RELOGIO_AJUSTE          14      // Cadência do relógio interno

#define REG_USER_START              15      // Início dos registradores de uso geral (15..99)
#define REG_USER_END                99

#define REG_OBJETO_OFFSET           100     // Acesso à situação do objeto X através de REG[100 + X]

// Limites padrão do sistema (endereços históricos 6593H / 6594H no MSX original)
#define DEFAULT_MAX_CARRIED         5       // Máximo de objetos na mão
#define DEFAULT_MAX_IN_CONTAINER    3       // Máximo de objetos dentro do objeto 3

// =============================================================================
// 3. VERBOS PRÉ-DEFINIDOS (Capítulo 4)
// =============================================================================
// Os verbos possuem uma distribuição obrigatória na tabela interna de 1 a 34.
// Verbos adicionais do autor iniciam a partir do índice 35 até 200.
// =============================================================================
typedef enum {
    VERBO_NORTE        = 1,
    VERBO_SUL          = 2,
    VERBO_LESTE        = 3,
    VERBO_OESTE        = 4,
    VERBO_GRAVE        = 5,
    VERBO_RECUPERE     = 6,
    VERBO_ENTRE        = 7,
    VERBO_SUBA         = 8,
    VERBO_SAIA         = 9,
    VERBO_DESCA        = 10,
    VERBO_HORAS        = 11,
    VERBO_QUANTO       = 12,
    VERBO_TEMOS        = 13,
    VERBO_RECOMECE     = 14,
    VERBO_HA           = 15,
    VERBO_GARIMPE      = 16,
    VERBO_PENSE        = 17,
    VERBO_GRITE        = 18,
    VERBO_CORRA        = 19,
    VERBO_PEGUE        = 20,
    VERBO_COLOQUE      = 21,
    VERBO_TROQUE       = 22,
    VERBO_COMPRE       = 23,
    VERBO_ROUBE        = 24,
    VERBO_TIRE         = 25,
    VERBO_QUEBRE       = 26,
    VERBO_SOLTE        = 27,
    VERBO_EXAMINE      = 28,
    VERBO_PROCURE      = 29,
    VERBO_OFERECA      = 30,
    VERBO_FACA         = 31,
    VERBO_JOGUE        = 32,
    VERBO_CONSERTE     = 33,
    VERBO_VENDA        = 34,
    VERBO_BEBA         = 35,
    VERBO_NORDESTE     = 36,
    VERBO_NOROESTE     = 37,
    VERBO_SUDESTE      = 38,
    VERBO_SUDOESTE     = 39,
    VERBO_ENCHA        = 40,
    VERBO_CUSTOM_START = 41
} Game_VerbId;

typedef struct {
    u8          id;     // ID do verbo (ex: 50, 51, 52)
    const char* nome;   // Nome e sinônimos separados por barra: "VERBO/SIN1/SIN2"
} Game_Verb;

// =============================================================================
// 4. OBJETOS, SITUAÇÃO E BYTE DE CONSISTÊNCIA (Capítulo 5)
// =============================================================================
// O sistema suporta até 99 objetos indexados de 1 a 99.
// Objetos 1, 2 e 3 possuem papéis estruturais fixos na engine.
// =============================================================================
#define OBJ_MAX                     99

#define OBJ_ID_LOCAL                1       // Objeto 1: "LOCAL" (reconhecimento do ambiente)
#define OBJ_ID_LUZ                  2       // Objeto 2: Objeto para iluminação (vela, tocha, lanterna)
#define OBJ_ID_CONTAINER            3       // Objeto 3: Objeto recipiente (mala, saco, mochila)

// Valores da Situação de um Objeto
#define OBJ_SIT_INEXISTENTE         0       // 0: Não existe (ou palavra de reconhecimento sintático)
// 1 a 99: Presente no local (sala) indicado pelo valor, listável no comando EXAMINE LOCAL.
#define OBJ_SIT_SALA_MIN            1
#define OBJ_SIT_SALA_MAX            99
// 101 a 199: Presente no local (valor - 100), porém oculto (não listado em EXAMINE LOCAL).
#define OBJ_SIT_OCULTO_OFFSET       100
#define OBJ_SIT_CARREGADO           250     // Na mão / inventário do jogador
#define OBJ_SIT_EM_OBJ3_ABERTO      251     // Dentro do objeto 3 liberado para ser pego (objeto 3 aberto)
#define OBJ_SIT_EM_OBJ3_FECHADO     253     // Dentro do objeto 3 trancado/fechado (não acessível)

// Bits do Byte de Consistência (bits 0 a 6 indicam predisposição para ações comuns)
// Se bit == 1, executa a respectiva FUNÇÃO padrão caso não haja comando específico de interceptação.
#define OBJ_CONSIST_PEGAR           (1 << 0) // Bit 0: Pode ser pego -> FUNÇÃO 6
#define OBJ_CONSIST_COLOCAR_OBJ3    (1 << 1) // Bit 1: Pode ser colocado no Obj 3 -> FUNÇÃO 7
#define OBJ_CONSIST_TROCAR          (1 << 2) // Bit 2: Pode ser trocado -> FUNÇÃO 8
#define OBJ_CONSIST_COMPRAR         (1 << 3) // Bit 3: Pode ser comprado -> FUNÇÃO 9
#define OBJ_CONSIST_ROUBAR          (1 << 4) // Bit 4: Pode ser roubado -> FUNÇÃO 10
#define OBJ_CONSIST_TIRAR           (1 << 5) // Bit 5: Pode ser tirado de algum lugar -> FUNÇÃO 11
#define OBJ_CONSIST_QUEBRAR         (1 << 6) // Bit 6: Pode ser quebrado -> FUNÇÃO 12

typedef struct {
    u8          id;                 // Identificador (1..99)
    u8          situacao_inicial;   // Situação inicial no início da partida
    u8          consistencia;       // Byte de consistência (máscara de bits 0..6)
    const char* nome;               // Nome e sinônimos separados por barra (ex: "VELA/VELAS")
    const char* desc;               // Descrição detalhada do objeto
} Game_Object;

// =============================================================================
// 5. POSIÇÕES / SALAS (Capítulo 6 - Suporte a 8 Pontos Cardeais)
// =============================================================================
// O sistema suporta até 99 posições indexadas de 1 a 99.
// Cada sala possui 8 saídas cardeais na ordem N / S / L / O / NE / NO / SE / SO.
// Regras das saídas:
// - 0: Sem passagem naquela direção (emite mensagem MSG 15: "É impossível ir...")
// - 1 a 99: Passagem direta para a posição correspondente (Reg 1 = valor).
// - > 100: Movimento condicional! Executa a FUNÇÃO = (valor - 100).
// =============================================================================
#define POSICAO_MAX                 99
#define POSICAO_COND_OFFSET         100

typedef enum {
    DIR_NORTE    = 0,
    DIR_SUL      = 1,
    DIR_LESTE    = 2,
    DIR_OESTE    = 3,
    DIR_NORDESTE = 4,
    DIR_NOROESTE = 5,
    DIR_SUDESTE  = 6,
    DIR_SUDOESTE = 7,
    DIR_COUNT    = 8
} Game_Direction;

typedef struct {
    u8          id;                 // Identificador da posição (1..99)
    u8          saidas[DIR_COUNT];  // Saídas [Norte, Sul, Leste, Oeste, Nordeste, Noroeste, Sudeste, Sudoeste]
    const char* desc;               // Descrição do ambiente (impressa em DESC ou ENTER)
} Game_Position;

// =============================================================================
// 6. INSTRUÇÕES DO SISTEMA / BYTECODES (Capítulo 8)
// =============================================================================
// Conjunto completo de 45 instruções do interpretador de bytecode.
// =============================================================================
typedef enum {
    OP_NOP = 0,     // NOP          : Sem efeito; não produz resultado.
    OP_MSG,         // MSG X        : Imprime mensagem X.
    OP_NVC,         // NVC          : Aguarda novo comando do jogador.
    OP_LLIST,       // LLIST        : Lista objetos no local corrente e aguarda novo comando.
    OP_CLIST,       // CLIST        : Lista objetos carregados (inventário) e aguarda novo comando.
    OP_DLIST,       // DLIST        : Lista objetos dentro do objeto 3 e aguarda novo comando.
    OP_OBJ,         // OBJ X        : Imprime o nome do objeto X (se X=0, usa objeto em evidência).
    OP_INC,         // INC X        : Incrementa o registrador X.
    OP_DEC,         // DEC X        : Decrementa o registrador X.
    OP_LDR,         // LDR X, Y     : Grava o registrador X com o valor Y.
    OP_SOMA,        // SOMA X, Y    : Soma ao registrador X o valor Y.
    OP_RND,         // RND X, Y     : Soma ao registrador X um valor aleatório entre 0 e Y.
    OP_REG_EQ,      // REG= X, Y, Z : Se Reg[X] == Y salta para instrução/label Z.
    OP_REG_GT,      // REG> X, Y, Z : Se Reg[X] > Y salta para instrução/label Z.
    OP_REG_LT,      // REG< X, Y, Z : Se Reg[X] < Y salta para instrução/label Z.
    OP_AQUI,        // AQUI X, Y    : Se objeto X está no local, salta para label Y (se X=0 usa evidência).
    OP_LOCAL,       // LOCAL X, Y   : Se local corrente == X, salta para label Y.
    OP_TEMOS,       // TEMOS X, Y   : Se objeto X está sendo carregado, salta para label Y (se X=0 usa evidência).
    OP_SOLTA,       // SOLTA X      : Solta objeto X (se X=100, solta todos; se X=0 usa evidência).
    OP_PEGA,        // PEGA X       : Carrega o objeto X (se X=0 usa evidência).
    OP_CRIA,        // CRIA X       : Coloca objeto X no local corrente (se X=0 usa evidência).
    OP_APAG,        // APAG X       : Apaga objeto X do local/inventário (se X=0 usa evidência).
    OP_GOSUB,       // GOSUB X      : Desvia execução para a FUNÇÃO X como subrotina (1 nível).
    OP_LIBR,        // LIBR         : Libera todos os objetos dentro do objeto 3 (situação 251).
    OP_TRC,         // TRC          : Tranca todos os objetos dentro do objeto 3 (situação 253).
    OP_POE,         // POE X        : Coloca objeto X dentro do objeto 3 (se X=0 usa evidência).
    OP_ESV,         // ESV          : Esvazia o objeto 3 (coloca objetos no local corrente).
    OP_OK,          // OK           : Imprime "Ok." e aguarda novo comando.
    OP_REGN,        // REGN X       : Imprime valor decimal do registrador X.
    OP_NVF,         // NVF          : Aguarda novo comando sem testar a FUNÇÃO 5.
    OP_REF,         // REF          : Retorna da FUNÇÃO 4 (tempo real).
    OP_FIM,         // FIM          : Finaliza a partida.
    OP_NEU,         // NEU          : Reinicia a partida (reset completo).
    OP_DESC,        // DESC         : Descreve o local corrente e aguarda novo comando.
    OP_RET,         // RET          : Retorna da chamada GOSUB.
    OP_GOTO,        // GOTO X       : Salta incondicionalmente para o label/instrução X.
    OP_PAUSA,       // PAUSA X      : Atraso de X segundos.
    OP_FLAG,        // FLAG X, Y    : Grava flag/registrador X com o valor Y.
    OP_EVID,        // EVID X       : Põe em evidência o objeto do buffer X (1=primeiro objeto, 2=segundo).
    OP_CLS,         // CLS          : Limpa a tela (campo central).
    OP_EVD_EQ,      // EVD= X, Y    : Se objeto em evidência == X, salta para label Y.
    OP_CHRS,        // CHRS X       : Imprime o caractere de código ASCII X.
    OP_PRT,         // PRT X, Y     : Imprime mensagem X seguida do nome do objeto Y e aguarda comando.
    OP_DNT,         // DNT X, Y     : Se objeto X estiver dentro do objeto 3, salta para label Y.
    OP_CMD,         // CMD X        : Equivale a executar o comando de índice X.
    OP_MAX
} Game_Opcode;

// Estrutura compacta de uma instrução compilada (4 bytes por instrução, ideal para Z80)
typedef struct {
    u8 op;  // Código de operação (Game_Opcode)
    u8 p1;  // Parâmetro 1 (Registrador, ID de Mensagem, ID de Objeto, Label, etc.)
    u8 p2;  // Parâmetro 2 (Valor literal, Label secundário, ID de Objeto)
    u8 p3;  // Parâmetro 3 (Label de destino para instruções de comparação REG=, REG>, REG<)
} Game_Instruction;

// =============================================================================
// 7. COMANDOS (Capítulo 7)
// =============================================================================
// Cabeçalho de reconhecimento estruturado: VERBO + OBJETO1 + OBJETO2.
// Ausência de objetos preenche 0 no cabeçalho.
// =============================================================================
typedef struct {
    u8                      verbo;              // ID do verbo requerido (1..200)
    u8                      objeto1;            // ID do 1º objeto requerido (0 = nenhum)
    u8                      objeto2;            // ID do 2º objeto requerido (0 = nenhum)
    const Game_Instruction* instructions;       // Sequência de instruções compiladas
    u8                      instruction_count;  // Quantidade de instruções
} Game_Command;

// =============================================================================
// 8. FUNÇÕES (Capítulo 10)
// =============================================================================
// Procedimentos disparados por eventos, saídas condicionais ou ações padrão.
// Funções 1 a 20 possuem papéis específicos e pré-definidos na arquitetura.
// =============================================================================
typedef enum {
    FUNC_RESET              = 1,    // Reset do jogo: executada ao iniciar a partida.
    FUNC_TIMER_BOMBA        = 2,    // Executada quando Reg 5 atinge zero (bomba-relógio).
    FUNC_ESCURO             = 3,    // Executada quando Reg 6 atinge 5 (passos no escuro).
    FUNC_REALTIME           = 4,    // Executada constantemente no loop de tempo real (se instrução 0 != NOP).
    FUNC_PRE_COMANDO        = 5,    // Executada no topo do Game Loop antes de ler comando (se instrução 0 != NOP).
    FUNC_PADRAO_PEGAR       = 6,    // Ação padrão: pegar objeto comum (Byte de Consistência bit 0).
    FUNC_PADRAO_COLOCAR_OBJ3= 7,    // Ação padrão: colocar no objeto 3 (bit 1).
    FUNC_PADRAO_TROCAR      = 8,    // Ação padrão: trocar objeto (bit 2).
    FUNC_PADRAO_COMPRAR     = 9,    // Ação padrão: comprar objeto (bit 3).
    FUNC_PADRAO_ROUBAR      = 10,   // Ação padrão: roubar objeto (bit 4).
    FUNC_PADRAO_TIRAR       = 11,   // Ação padrão: tirar objeto de recipiente (bit 5).
    FUNC_PADRAO_QUEBRAR     = 12,   // Ação padrão: quebrar objeto (bit 6).
    FUNC_PADRAO_SOLTAR      = 13,   // Ação padrão: soltar objeto.
    FUNC_PADRAO_EXAMINAR    = 14,   // Ação padrão: examinar objeto.
    FUNC_PADRAO_PROCURAR    = 15,   // Ação padrão: procurar objetos.
    FUNC_PADRAO_OFERECER    = 16,   // Ação padrão: dar / oferecer.
    FUNC_PADRAO_FAZER       = 17,   // Ação padrão: fazer / construir.
    FUNC_PADRAO_JOGAR       = 18,   // Ação padrão: jogar / arremessar.
    FUNC_PADRAO_CONSERTAR   = 19,   // Ação padrão: consertar.
    FUNC_PADRAO_VENDER      = 20,   // Ação padrão: vender.
    FUNC_USER_START         = 21    // Funções customizadas do autor (21..200)
} Game_FunctionId;

typedef struct {
    u8                      id;                 // ID da função (1..200)
    const Game_Instruction* instructions;       // Sequência de instruções
    u8                      instruction_count;  // Quantidade de instruções
} Game_Function;

// =============================================================================
// 9. MENSAGENS (Capítulo 9)
// =============================================================================
// Mensagens pré-definidas utilizadas autonomamente pelo sistema.
// =============================================================================
typedef enum {
    MSG_INTRO               = 10,   // Apresentação inicial do jogo
    MSG_ACHEI               = 11,   // "Achei o que você queria."
    MSG_ESCURO              = 12,   // "Está muito escuro aqui. É melhor arranjar alguma luz..."
    MSG_NAO_ENTENDI         = 13,   // "Perdão, não entendi..."
    MSG_MOVIMENTO_INVALIDO  = 14,   // "É impossível ir nesta direção."
    MSG_NAO_POSSIVEL        = 15,   // "Isto não é possível."
    MSG_NAO_TEMOS           = 16,   // "Nós não temos (objeto)."
    MSG_JA_TEMOS            = 17,   // "Nós já temos (objeto)."
    MSG_NAO_ESTOU_VENDO     = 18,   // "Eu não estou vendo..."
    MSG_OBJETO_COMUM        = 19,   // "É apenas (objeto)."
    MSG_CARGA_MAXIMA        = 20,   // "Não dá para carregar mais nada."
    MSG_OBJ3_LOTADO         = 21    // "Não cabe mais nada."
} Game_SystemMessageId;

typedef struct {
    u8          id;     // ID da mensagem (1..200)
    const char* text;   // Texto da mensagem
} Game_Message;

// =============================================================================
// 10. ESTADO EM EXECUÇÃO DO JOGO (Runtime State)
// =============================================================================
typedef struct {
    u8   registers[REG_COUNT];  // Array de 256 registradores
    u8   obj_evidencia;         // Objeto atualmente em evidência
    u8   obj_buffer[2];         // Objeto 1 e Objeto 2 reconhecidos na frase atual
    u8   verbo_atual;           // Verbo reconhecido na frase atual
    u8   gosub_ret_func;        // ID da função de retorno para chamada GOSUB
    u8   gosub_ret_pc;          // PC de retorno da subrotina GOSUB
    bool flag_espera_cmd;       // True quando instrução aguarda novo comando (NVC, OK, etc.)
    bool flag_fim;              // True quando o jogo terminou (instrução FIM)
    bool flag_reiniciar;        // True quando o jogo deve reiniciar (instrução NEU)
} Game_State;

// =============================================================================
// 11. ESTRUTURA GERAL DA BASE DE DADOS DO JOGO
// =============================================================================
typedef struct {
    const char*                 titulo;
    u8                          posicao_inicial;
    u8                          max_carregados;
    u8                          max_no_obj3;
    
    const Game_Position* const* posicoes;       // Tabela de ponteiros de posições (1..99)
    u8                          num_posicoes;

    const Game_Object* const*   objetos;        // Tabela de ponteiros de objetos (1..99)
    u8                          num_objetos;

    const Game_Command* const*  comandos;       // Tabela de comandos customizados
    u8                          num_comandos;

    const Game_Function* const* funcoes;        // Tabela de funções (1..200)
    u8                          num_funcoes;

    const Game_Message* const*  mensagens;      // Tabela de mensagens
    u8                          num_mensagens;

    const Game_Verb* const*     verbos;         // Tabela de verbos customizados
    u8                          num_verbos;

    const u8*                   atalhos_shift;  // 10 atalhos de acentos calculados (Shift+0..9)
} Game_Database;

#endif // GAME_TYPES_H
