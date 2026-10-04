// ____________________________
// Z-Realm (zrealm-msx) - Gerenciador de Memory Mapper para MSX-DOS 2
// Implementação com proteção total contra vazamento de memória no DOS 2
//─────────────────────────────────────────────────────────────────────────────
#include "mapper.h"
#if (TARGET_TYPE == TYPE_DOS)
#include "dos.h"

// Tabela de segmentos físicos alocados do sistema operacional
static DOS_Segment s_AllocatedSegments[MAPPER_MAX_SEGMENTS];

// Estado global do gerenciador
MapperState g_MapperState = { 0, 0xFF, 0 };

// Inicializa e aloca os blocos de 16 KB no Memory Mapper
bool MAPPER_Init(u8 minRequiredSegments)
{
	u8 i;
	u8 slotId;
	g_MapperState.TotalSegments = 0;
	g_MapperState.ActiveSegment = 0xFF;

	// 1. Inicializa a Extended BIOS e descobre a Jump Table do Mapper no MSX-DOS 2
	if (!DOSMapper_Init())
	{
		return FALSE;
	}

	// 2. Salva o segmento atual da Página 2 para restaurar na saída
	g_MapperState.OriginalPage2 = DOSMapper_GetPage2();

__asm
	.globl _g_DOS_VarTable
__endasm;

	// Identifica o slot do mapper principal a partir da tabela de variáveis do DOS 2
	slotId = g_DOS_VarTable ? g_DOS_VarTable->Slot : 0;

	// 3. Aloca os segmentos solicitados via ALL_SEG
	for (i = 0; i < MAPPER_MAX_SEGMENTS; i++)
	{
		DOS_Segment seg;
		// Aloca segmento de usuário preferencialmente no mapper primário
		if (!DOSMapper_Alloc(DOS_ALLOC_USER, DOS_SEGSLOT_THISFIRST | slotId, &seg))
		{
			// Não há mais segmentos disponíveis
			break;
		}

		s_AllocatedSegments[g_MapperState.TotalSegments].Number = seg.Number;
		s_AllocatedSegments[g_MapperState.TotalSegments].Slot = seg.Slot;
		g_MapperState.TotalSegments++;
	}

	// 4. Verifica se conseguiu a cota mínima
	if (g_MapperState.TotalSegments < minRequiredSegments)
	{
		MAPPER_Cleanup();
		return FALSE;
	}

	return TRUE;
}

// Chaveia o segmento lógico para a Página 2 (8000h - BFFFh)
void MAPPER_SetPage2(u8 logicalSegment)
{
	if (logicalSegment >= g_MapperState.TotalSegments)
		return;

	if (g_MapperState.ActiveSegment == logicalSegment)
		return;

	DOSMapper_SetPage2(s_AllocatedSegments[logicalSegment].Number);
	g_MapperState.ActiveSegment = logicalSegment;
}

// Retorna qual segmento lógico está ativo
u8 MAPPER_GetPage2(void)
{
	return g_MapperState.ActiveSegment;
}

// Liberação completa de todos os recursos alocados
void MAPPER_Cleanup(void)
{
	u8 i;

	// Restaura a Página 2 original
	if (g_MapperState.OriginalPage2 != 0)
	{
		DOSMapper_SetPage2(g_MapperState.OriginalPage2);
	}

	// Libera cada um dos segmentos alocados de volta ao DOS 2
	for (i = 0; i < g_MapperState.TotalSegments; i++)
	{
		DOSMapper_FreeStruct(&s_AllocatedSegments[i]);
	}

	g_MapperState.TotalSegments = 0;
	g_MapperState.ActiveSegment = 0xFF;
}

#elif (TARGET_TYPE == TYPE_ROM)

// No cartucho MegaROM, todos os segmentos lógicos estão disponíveis diretamente na ROM
MapperState g_MapperState = { MAPPER_MAX_SEGMENTS, 0xFF, 0 };

bool MAPPER_Init(u8 minRequiredSegments)
{
	(void)minRequiredSegments;
	g_MapperState.TotalSegments = MAPPER_MAX_SEGMENTS;
	g_MapperState.ActiveSegment = 0xFF;

	// Configura o primeiro segmento de dados (Bank 1) na Página 2 (8000h)
	MAPPER_SetPage2(0);
	return TRUE;
}

void MAPPER_SetPage2(u8 logicalSegment)
{
	if (g_MapperState.ActiveSegment == logicalSegment)
		return;

	// No cartucho MegaROM:
	// - Bank 0 é a Engine e cabeçalho mestre em Página 1 (4000h - 7FFFh).
	// - Banks 1..N são os segmentos de dados chaveados na Página 2 (8000h - BFFFh).
	#if (ROM_MAPPER == ROM_ASCII16)
		Poke(0x77FF, (u8)(logicalSegment + 1));
	#elif (ROM_MAPPER == ROM_ASCII8)
		Poke(0x7000, (u8)((logicalSegment + 1) * 2));
		Poke(0x7800, (u8)((logicalSegment + 1) * 2 + 1));
	#elif (ROM_MAPPER == ROM_KONAMI)
		Poke(0x8000, (u8)((logicalSegment + 1) * 2));
		Poke(0xA000, (u8)((logicalSegment + 1) * 2 + 1));
	#else
		Poke(0x77FF, (u8)(logicalSegment + 1));
	#endif

	g_MapperState.ActiveSegment = logicalSegment;
}

u8 MAPPER_GetPage2(void)
{
	return g_MapperState.ActiveSegment;
}

void MAPPER_Cleanup(void)
{
	g_MapperState.ActiveSegment = 0xFF;
}

#endif
