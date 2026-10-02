// _____________________________________________________________________________
//
//  Parser Module Implementation (Capítulo 4, 5 e 7)
// _____________________________________________________________________________

#include "parser.h"

#if defined(MSXGL)
    #include "core.h"
    #include "string.h"
#else
    #include <string.h>
    #include <ctype.h>
#endif

// Tabela de nomes padrão dos verbos do sistema (Capítulo 4.1 + 8 Direções Cardeais)
#define NUM_DEFAULT_VERBS 40
static const char* const g_DefaultVerbs[NUM_DEFAULT_VERBS] = {
    "NORTE/N",
    "SUL/S",
    "LESTE/L/E/ESTE",
    "OESTE/O/W/WEST",
    "GRAVE",
    "RECUPERE",
    "ENTRE/ABRA/DESTRANQUE",
    "SUBA",
    "SAIA",
    "DESCA/DESCER",
    "HORAS",
    "QUANTO",
    "TEMOS/INV/I/INVENTARIO/INVENTARIO",
    "RECOMECE/REINICIE/REINICIAR/RESTART",
    "HA",
    "GARIMPE",
    "PENSE/PENSAR/DICA/DICAS/AJUDA",
    "GRITE",
    "CORRA",
    "PEGUE/PEGAR/APANHE",
    "COLOQUE/PONHA/GUARDE/USE/AMARRE",
    "TROQUE",
    "COMPRE",
    "ROUBE",
    "TIRE",
    "QUEBRE",
    "SOLTE/LARGUE/DEIXE",
    "EXAMINE/OLHE/VER/EX",
    "PROCURE/BUSQUE",
    "OFERECA/DOE/DE",
    "FACA/CONSTRUA/ACENDA/RISQUE",
    "JOGUE/ATIRE",
    "CONSERTE/REPARE",
    "VENDA",
    "BEBA/BEBER/TOME/TOMAR",
    "NORDESTE/NE",
    "NOROESTE/NO/NW",
    "SUDESTE/SE",
    "SUDOESTE/SO/SW",
    "ENCHA/ENCHER/ABASTECA"
};

// Normaliza caracteres acentuados e cedilha para comparação case-insensitive
static char NormalizeChar(char c)
{
    u8 uc = (u8)c;
    // Cedilha (Ç, ç, Latin-1 e CP437)
    if (uc == 0x80 || uc == 0x87 || uc == 0xC7 || uc == 0xE7) return 'C';
    // A acentuado (Á, á, À, à, Ã, ã, Â, â)
    if (uc == 0x84 || uc == 0xA0 || uc == 0x8F || uc == 0x85 || uc == 0xB0 || uc == 0xB1 || uc == 0x8C || uc == 0x83 ||
        uc == 0xC1 || uc == 0xE1 || uc == 0xC0 || uc == 0xE0 || uc == 0xC3 || uc == 0xE3 || uc == 0xC2 || uc == 0xE2) return 'A';
    // E acentuado (É, é, Ê, ê)
    if (uc == 0x90 || uc == 0x82 || uc == 0x8D || uc == 0x88 ||
        uc == 0xC9 || uc == 0xE9 || uc == 0xCA || uc == 0xEA) return 'E';
    // I acentuado (Í, í)
    if (uc == 0x89 || uc == 0xA1 || uc == 0xCD || uc == 0xED) return 'I';
    // O acentuado (Ó, ó, Ô, ô, Õ, õ)
    if (uc == 0x8A || uc == 0xA2 || uc == 0x8E || uc == 0x93 || uc == 0xB4 || uc == 0xB5 || uc == 0x95 ||
        uc == 0xD3 || uc == 0xF3 || uc == 0xD4 || uc == 0xF4 || uc == 0xD5 || uc == 0xF5) return 'O';
    // U acentuado (Ú, ú, Ü, ü)
    if (uc == 0x8B || uc == 0xA3 || uc == 0xDA || uc == 0xFA || uc == 0x81 || uc == 0x9A || uc == 0xDC || uc == 0xFC || uc == 0x9F) return 'U';
    // Letras minúsculas normais
    if (uc >= 'a' && uc <= 'z') return (char)(uc - ('a' - 'A'));
    return c;
}

// Compara uma palavra com uma lista de sinônimos separados por barra: "NOME/SIN1/SIN2"
static bool MatchWord(const char* word, const char* synonyms)
{
    const char* p;
    const char* start;
    u8 wlen, slen;

    if (word == NULL || synonyms == NULL) return FALSE;
    wlen = 0;
    while (word[wlen]) wlen++;
    if (wlen == 0) return FALSE;

    p = synonyms;
    while (*p)
    {
        start = p;
        while (*p && *p != '/') p++;
        slen = (u8)(p - start);

        if (wlen == slen)
        {
            u8 i;
            bool match = TRUE;
            for (i = 0; i < wlen; i++)
            {
                if (NormalizeChar(word[i]) != NormalizeChar(start[i]))
                {
                    match = FALSE;
                    break;
                }
            }
            if (match) return TRUE;
        }

        if (*p == '/') p++;
    }

    return FALSE;
}

static u8 StrLen(const char* s)
{
    u8 len = 0;
    if (s == NULL) return 0;
    while (*s++) len++;
    return len;
}

static bool StrEqual(const char* s1, const char* s2)
{
    if (s1 == NULL || s2 == NULL) return FALSE;
    while (*s1 && (*s1 == *s2)) { s1++; s2++; }
    return (*s1 == *s2);
}

static void StrCat(char* dest, const char* src, u8 max_len)
{
    u8 dlen = StrLen(dest);
    if (dest == NULL || src == NULL) return;
    while (*src && dlen < (max_len - 1))
    {
        dest[dlen++] = *src++;
    }
    dest[dlen] = '\0';
}

// -----------------------------------------------------------------------------
// Parser_IsNoiseWord
// -----------------------------------------------------------------------------
bool Parser_IsNoiseWord(const char* word)
{
    static const char* const noise[] = {
        "AS", "OS", "UM", "UMA", "UNS", "UMAS",
        "DE", "DO", "DA", "DOS", "DAS",
        "EM", "NO", "NA", "NOS", "NAS",
        "PARA", "PRA", "COM", "POR", NULL
    };
    u8 i = 0;
    while (noise[i] != NULL)
    {
        if (StrEqual(word, noise[i])) return TRUE;
        i++;
    }
    return FALSE;
}

// -----------------------------------------------------------------------------
// Parser_FindVerb
// -----------------------------------------------------------------------------
static u8 Parser_FindVerb(const char* word, const Game_Database* db)
{
    u8 i;

    // 1. Checa verbos customizados do banco de dados primeiro
    if (db != NULL && db->verbos != NULL)
    {
        for (i = 0; i < db->num_verbos; i++)
        {
            const Game_Verb* v = db->verbos[i];
            if (v != NULL && v->nome != NULL)
            {
                if (MatchWord(word, v->nome))
                {
                    return v->id;
                }
            }
        }
    }

    // 2. Checa tabela de verbos padrão do sistema
    for (i = 0; i < NUM_DEFAULT_VERBS; i++)
    {
        if (MatchWord(word, g_DefaultVerbs[i]))
        {
            return (i + 1);
        }
    }

    return 0;
}

// -----------------------------------------------------------------------------
// Parser_FindObject
// -----------------------------------------------------------------------------
static u8 Parser_FindObject(const char* word, const Game_Database* db)
{
    u8 i;
    if (db == NULL || db->objetos == NULL) return 0;

    for (i = 0; i < db->num_objetos; i++)
    {
        const Game_Object* obj = db->objetos[i];
        if (obj != NULL && obj->nome != NULL)
        {
            if (MatchWord(word, obj->nome))
            {
                return obj->id;
            }
        }
    }
    return 0;
}

// -----------------------------------------------------------------------------
// Parser_Parse
// -----------------------------------------------------------------------------
bool Parser_Parse(const char* input, const Game_Database* db, Game_State* state, char* parsed_echo, u8 echo_max_len)
{
    char word[32];
    u8 word_idx = 0;
    u8 meaningful_words = 0;
    const char* p = input;

    state->verbo_atual = 0;
    state->obj_buffer[0] = 0;
    state->obj_buffer[1] = 0;

    if (parsed_echo != NULL && echo_max_len > 0)
    {
        parsed_echo[0] = '\0';
    }

    if (input == NULL) return FALSE;

    while (*p)
    {
        // Pula espaços
        while (*p == ' ' || *p == '\t' || *p == '\r' || *p == '\n') p++;
        if (*p == '\0') break;

        // Extrai próxima palavra
        word_idx = 0;
        while (*p && *p != ' ' && *p != '\t' && *p != '\r' && *p != '\n')
        {
            if (word_idx < (sizeof(word) - 1))
            {
                word[word_idx++] = *p;
            }
            p++;
        }
        word[word_idx] = '\0';

        // Verifica se é palavra de ruído / conectivo
        if (Parser_IsNoiseWord(word))
        {
            continue;
        }

        // Primeira palavra significativa: VERBO
        if (meaningful_words == 0)
        {
            state->verbo_atual = Parser_FindVerb(word, db);
            meaningful_words++;

            if (parsed_echo != NULL && state->verbo_atual > 0)
            {
                // Obtém o nome canônico do verbo (customizado ou padrão)
                const char* vname = "";
                if (db != NULL && db->verbos != NULL)
                {
                    u8 vi;
                    for (vi = 0; vi < db->num_verbos; vi++)
                    {
                        if (db->verbos[vi] != NULL && db->verbos[vi]->id == state->verbo_atual)
                        {
                            vname = db->verbos[vi]->nome;
                            break;
                        }
                    }
                }
                if (vname[0] == '\0' && state->verbo_atual <= NUM_DEFAULT_VERBS)
                {
                    vname = g_DefaultVerbs[state->verbo_atual - 1];
                }

                u8 k = 0;
                while (vname[k] && vname[k] != '/' && k < (echo_max_len - 1))
                {
                    parsed_echo[k] = vname[k];
                    k++;
                }
                parsed_echo[k] = '\0';
            }
        }
        // Segunda palavra significativa: OBJETO 1
        else if (meaningful_words == 1)
        {
            state->obj_buffer[0] = Parser_FindObject(word, db);
            meaningful_words++;

            if (parsed_echo != NULL && state->obj_buffer[0] > 0)
            {
                u8 cur = StrLen(parsed_echo);
                if (cur < (echo_max_len - 2))
                {
                    parsed_echo[cur++] = ' ';
                    parsed_echo[cur] = '\0';
                    StrCat(parsed_echo, word, echo_max_len);
                }
            }
        }
        // Terceira palavra significativa: OBJETO 2
        else if (meaningful_words == 2)
        {
            state->obj_buffer[1] = Parser_FindObject(word, db);
            meaningful_words++;

            if (parsed_echo != NULL && state->obj_buffer[1] > 0)
            {
                u8 cur = StrLen(parsed_echo);
                if (cur < (echo_max_len - 2))
                {
                    parsed_echo[cur++] = ' ';
                    parsed_echo[cur] = '\0';
                    StrCat(parsed_echo, word, echo_max_len);
                }
            }
            break;
        }
    }

    return (state->verbo_atual > 0);
}
