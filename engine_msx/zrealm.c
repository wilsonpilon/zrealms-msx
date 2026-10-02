// ____________________________
// Z-Realm (zrealm-msx) - Engine Core para MSX 2 & MSX-DOS 2
// Fase 2: Bootstrap MSX-DOS 2, Memory Mapper (Página 2) e V9938 SCREEN 4
//─────────────────────────────────────────────────────────────────────────────
#include "core.h"
#include "dos.h"
#include "dos_mapper.h"
#include "vdp.h"
#include "keyboard.h"
#include "mapper.h"
#include "vdp_screen4.h"
#include "system.h"
#include "loader.h"

#define Halt() __asm__("halt")

static void PrintHex(u8 val)
{
	static const c8 hex[] = "0123456789ABCDEF";
	DOS_CharOutput(hex[val >> 4]);
	DOS_CharOutput(hex[val & 0x0F]);
}

// Pequeno atraso calibrado em software (não depende de interrupção VDP)
static void WaitDelay(u16 count)
{
	volatile u16 i;
	for (i = 0; i < count; i++)
	{
		__asm nop __endasm;
	}
}

// Ponto de entrada do executável .COM no MSX-DOS 2
void main(void)
{
	const MasterHeader* header;
	BinaryTileset* tileset;

	DOS_StringOutput("=== Z-Realm MSX2 Engine ===\r\n$");

	// 1. Carrega o jogo binário (HEADER.BIN e GAME.DAT) via Loader do MSX-DOS 2
	if (!LOADER_LoadGame("HEADER.BIN", "GAME.DAT"))
	{
		DOS_StringOutput("Erro ao carregar dados do jogo!\r\n$");
		return;
	}

	header = LOADER_GetMasterHeader();
	DOS_StringOutput("Jogo carregado! Segmentos: $");
	PrintHex((u8)header->SegmentCount);
	DOS_StringOutput(" Recursos: $");
	PrintHex((u8)header->ResourceCount);
	DOS_StringOutput(" Sala Inicial: $");
	PrintHex((u8)header->InitialRoomID);
	DOS_StringOutput("\r\n$");

	// 2. Mapeia e carrega o tileset inicial na VRAM
	tileset = LOADER_GetTileset(header->InitialTileset);
	if (!tileset)
	{
		DOS_StringOutput("Erro: Tileset inicial nao encontrado!\r\n$");
		MAPPER_Cleanup();
		return;
	}

	// 3. Inicializa o modo de vídeo V9938 SCREEN 4 (Graphic 3)
	VDP_InitScreen4();

	// 4. Copia os padrões e cores do tileset para os 3 bancos da SCREEN 4
	VDP_LoadTilesetAllBanks(tileset->PatternTable, tileset->ColorTable);
	VDP_ClearHUDAndDialogue(255);

	// 5. Loop de Demonstração:
	// Alterna entre a Sala 1 (Entrada) e a Sala 2 (Câmara dos Pilares) carregadas da RAM do Mapper
	{
		u8 cycle;
		for (cycle = 0; cycle < 2; cycle++)
		{
			u16 t;
			u16 roomId = (cycle == 0) ? header->InitialRoomID : 2;
			BinaryRoom* room = LOADER_GetRoom(roomId);
			if (room)
			{
				VDP_DrawRoomViewport(room->TileMatrix);
			}

			// Exibe cada sala por ~3.3 segundos (30 passos de ~110ms)
			for (t = 0; t < 30; t++)
			{
				if (Keyboard_IsKeyPressed(KEY_ESC))
				{
					goto cleanup;
				}
				if (Keyboard_IsKeyPressed(KEY_SPACE))
				{
					// Pular para o próximo ciclo se apertar espaço
					break;
				}
				WaitDelay(6000);
			}
		}
	}

cleanup:
	DOS_InterSlotCall(g_EXPTBL[0], R_INITXT);
	MAPPER_Cleanup();
	DOS_StringOutput("Z-Realm: Execucao finalizada com sucesso! RAM liberada.\r\n$");
}
