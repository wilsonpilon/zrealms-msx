// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Memory Mapper para MSX-DOS 2
// Arquitetura: Janela Dinâmica na Página 2 (8000h - BFFFh)
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "dos_mapper.h"

// Máximo de segmentos gerenciados pelo cache em RAM (32 x 16KB = 512KB de RAM)
#define MAPPER_MAX_SEGMENTS 32

// Endereço base da Página 2 na memória do Z80
#define MAPPER_PAGE2_BASE   0x8000

// Informações de status do subsistema de memória
typedef struct
{
	u8 TotalSegments;     // Segmentos físicos alocados com sucesso
	u8 ActiveSegment;     // Segmento lógico atualmente visível na Página 2
	u8 OriginalPage2;     // Segmento original da Página 2 no momento da inicialização
} MapperState;

extern MapperState g_MapperState;

// Inicializa o subsistema de Mapper, descobre a RAM disponível e aloca os segmentos de trabalho.
// Retorna TRUE se conseguiu alocar o mínimo necessário para a execução.
bool MAPPER_Init(u8 minRequiredSegments);

// Chaveia o segmento lógico especificado diretamente para a Página 2 (8000h - BFFFh).
void MAPPER_SetPage2(u8 logicalSegment);

// Retorna o segmento lógico atualmente mapeado na Página 2.
u8 MAPPER_GetPage2(void);

// Obtém um ponteiro para a base da Página 2 (0x8000)
inline void* MAPPER_GetPage2Address(void) { return (void*)MAPPER_PAGE2_BASE; }

// Libera graciosamente todos os segmentos alocados de volta ao MSX-DOS 2 e restaura a Página 2.
void MAPPER_Cleanup(void);
