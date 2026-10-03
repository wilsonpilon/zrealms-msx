// ____________________________
// Z-Realm (zrealm-msx) - Engine Core para MSX 2 & MSX-DOS 2
// Fase 4: Gameplay Engine, Movimentação no Grid & Colisão (Subfase 4.1)
//─────────────────────────────────────────────────────────────────────────────
#include "core.h"
#include "dos.h"
#include "dos_mapper.h"
#include "vdp.h"
#include "keyboard.h"
#include "joystick.h"
#include "mapper.h"
#include "vdp_screen4.h"
#include "system.h"
#include "loader.h"
#include "world.h"
#include "hero.h"
#include "entity.h"
#include "vm.h"

#define Halt() __asm__("halt")

static void PrintHex(u8 val)
{
	static const c8 hex[] = "0123456789ABCDEF";
	DOS_CharOutput(hex[val >> 4]);
	DOS_CharOutput(hex[val & 0x0F]);
}

// Ponto de entrada do executável .COM no MSX-DOS 2
void main(void)
{
	const MasterHeader* header;
	u8 startX, startY;

	DOS_StringOutput("=== Z-Realm MSX2 Engine (Gameplay Core) ===\r\n$");

	// 1. Carrega o banco binário do jogo via handles de arquivo do DOS 2
	if (!LOADER_LoadGame("HEADER.BIN", "GAME.DAT"))
	{
		DOS_StringOutput("Erro fatal: Falha ao carregar HEADER.BIN ou GAME.DAT!\r\n$");
		return;
	}

	header = LOADER_GetMasterHeader();
	DOS_StringOutput("Dados carregados com sucesso! Segmentos: $");
	PrintHex((u8)header->SegmentCount);
	DOS_StringOutput(" Recursos: $");
	PrintHex((u8)header->ResourceCount);
	DOS_StringOutput(" Sala Inicial: $");
	PrintHex((u8)header->InitialRoomID);
	DOS_StringOutput("\r\n$");

	// 2. Inicializa o subsistema de mundo, de entidades e a VM de eventos
	WORLD_Init();
	ENTITY_Init();
	VM_Init();

	// 3. Inicializa o processador de vídeo V9938 em SCREEN 4 (Graphic 3)
	VDP_InitScreen4();
	VDP_ClearHUDAndDialogue(255); // Preenche HUD e diálogo com tile vazio

	// 4. Carrega a sala inicial no Memory Mapper e desenha na tela
	if (!WORLD_LoadRoom(header->InitialRoomID))
	{
		DOS_StringOutput("Erro: Nao foi possivel carregar a sala inicial!\r\n$");
		goto cleanup;
	}

	// 5. Inicializa o herói na coordenada de partida com o Sprite 1
	startX = (header->InitialHeroX < VIEWPORT_WIDTH) ? header->InitialHeroX : 16;
	startY = (header->InitialHeroY < VIEWPORT_HEIGHT) ? header->InitialHeroY : 9;
	HERO_Init(startX, startY, 1);

	// 6. Loop Principal de Gameplay (Sincronizado a 50/60 Hz no V-Blank)
	EnableInterrupt();

	while (TRUE)
	{
		// Aguarda o próximo ciclo de interrupção vertical (V-Blank)
		Halt();

		// Atualiza herói (controles, física, colisão e ação)
		HERO_Update();

		// Atualiza ciclo de vida e IAs de todas as entidades ativas
		ENTITY_Update();

		// Tecla ESC para encerrar a partida e retornar ao sistema operacional
		if (Keyboard_IsKeyPressed(KEY_ESC))
		{
			break;
		}
	}

cleanup:
	// Oculta o sprite do herói e todas as entidades ativas
	VDP_Screen4_HideSprite(0);
	ENTITY_Cleanup();

	// Restaura o modo de texto SCREEN 0 padrão do MSX-DOS 2 via BIOS
	DOS_InterSlotCall(g_EXPTBL[0], R_INITXT);

	// Libera todos os segmentos do Memory Mapper alocados via EXTBIOS
	MAPPER_Cleanup();

	DOS_StringOutput("Z-Realm: Sessao finalizada com sucesso! RAM liberada.\r\n$");
}
