// ____________________________
// Z-Realm (zrealm-msx) - Máquina Virtual de Eventos (MSX 2 Bytecode VM)
// Subfase 4.3: Interpretador Z80 de Opcodes, Flags Globais e Inventário
//─────────────────────────────────────────────────────────────────────────────
#include "vm.h"
#include "dos.h"
#include "ui.h"

u8              g_VMFlags[VM_MAX_FLAGS];
VMInventorySlot g_VMInventory[VM_MAX_INVENTORY];
u8              g_VMInventoryCount;
VMHeroStats     g_HeroStats;
c8              g_VMActiveMessage[VM_MAX_MSG_LEN];
u16             g_VMActiveStringID;
bool            g_VMHasMessage;
u8              g_VMLastSFX;

// Buffer de execução em RAM local (Página 1) imune a chaveamentos de segmento
static u8 s_VMScriptBuffer[VM_SCRIPT_BUF_SZ];

void VM_Init(void)
{
	u16 i;
	for (i = 0; i < VM_MAX_FLAGS; i++)
	{
		g_VMFlags[i] = 0;
	}

	for (i = 0; i < VM_MAX_INVENTORY; i++)
	{
		g_VMInventory[i].ItemID = 0;
		g_VMInventory[i].Quantity = 0;
	}
	g_VMInventoryCount = 0;

	g_HeroStats.HP = 75;
	g_HeroStats.MaxHP = 100;
	g_HeroStats.MP = 30;
	g_HeroStats.MaxMP = 30;
	g_HeroStats.Level = 1;
	g_HeroStats.Gold = 0;

	VM_ClearMessage();
	g_VMLastSFX = 0;
}

void VM_ClearMessage(void)
{
	g_VMActiveMessage[0] = 0;
	g_VMActiveStringID = 0xFFFF;
	g_VMHasMessage = FALSE;
}

void VM_SetFlag(u8 flagID, u8 val)
{
	g_VMFlags[flagID] = val;
}

u8 VM_GetFlag(u8 flagID)
{
	return g_VMFlags[flagID];
}

bool VM_GiveItem(u16 itemID, u8 quantity)
{
	u8 i;
	if (itemID == 0 || quantity == 0)
	{
		return FALSE;
	}

	// 1. Procura se já possui o item no inventário
	for (i = 0; i < g_VMInventoryCount; i++)
	{
		if (g_VMInventory[i].ItemID == itemID)
		{
			g_VMInventory[i].Quantity += quantity;
			return TRUE;
		}
	}

	// 2. Adiciona em novo slot se houver capacidade
	if (g_VMInventoryCount < VM_MAX_INVENTORY)
	{
		g_VMInventory[g_VMInventoryCount].ItemID = itemID;
		g_VMInventory[g_VMInventoryCount].Quantity = quantity;
		g_VMInventoryCount++;
		return TRUE;
	}

	return FALSE; // Inventário cheio
}

bool VM_TakeItem(u16 itemID, u8 quantity)
{
	u8 i, j;
	if (itemID == 0 || quantity == 0)
	{
		return FALSE;
	}

	for (i = 0; i < g_VMInventoryCount; i++)
	{
		if (g_VMInventory[i].ItemID == itemID)
		{
			if (g_VMInventory[i].Quantity <= quantity)
			{
				// Remove o slot e compacta o inventário
				for (j = i; j < g_VMInventoryCount - 1; j++)
				{
					g_VMInventory[j].ItemID = g_VMInventory[j + 1].ItemID;
					g_VMInventory[j].Quantity = g_VMInventory[j + 1].Quantity;
				}
				g_VMInventoryCount--;
				g_VMInventory[g_VMInventoryCount].ItemID = 0;
				g_VMInventory[g_VMInventoryCount].Quantity = 0;
			}
			else
			{
				g_VMInventory[i].Quantity -= quantity;
			}
			return TRUE;
		}
	}

	return FALSE; // Item não encontrado
}

bool VM_HasItem(u16 itemID)
{
	return (VM_GetItemQuantity(itemID) > 0);
}

u8 VM_GetItemQuantity(u16 itemID)
{
	u8 i;
	if (itemID == 0)
	{
		return 0;
	}

	for (i = 0; i < g_VMInventoryCount; i++)
	{
		if (g_VMInventory[i].ItemID == itemID)
		{
			return g_VMInventory[i].Quantity;
		}
	}
	return 0;
}

void VM_PlaySFX(u8 sfxID)
{
	g_VMLastSFX = sfxID;
	DOS_Beep();
}

bool VM_ExecuteScript(u16 scriptID)
{
	u16 length;

	if (scriptID == 0 || scriptID == 0xFFFF)
	{
		return FALSE;
	}

	// Copia o script do Memory Mapper para a RAM local (Página 1)
	length = LOADER_CopyScript(scriptID, s_VMScriptBuffer, VM_SCRIPT_BUF_SZ);
	if (length == 0)
	{
		return FALSE;
	}

	return VM_ExecuteBytecode(s_VMScriptBuffer, length);
}

bool VM_ExecuteBytecode(const u8* bytecode, u16 length)
{
	u16 pc = 0;

	if (!bytecode || length == 0)
	{
		return FALSE;
	}

	while (pc < length)
	{
		u8 op = bytecode[pc];

		switch (op)
		{
			case VM_OP_NOP:
				pc++;
				break;

			case VM_OP_MSG:
			{
				u16 strID;
				if (pc + 2 >= length) return FALSE;
				strID = (u16)bytecode[pc + 1] | ((u16)bytecode[pc + 2] << 8);
				g_VMActiveStringID = strID;
				LOADER_CopyString(strID, g_VMActiveMessage, VM_MAX_MSG_LEN);
				g_VMHasMessage = TRUE;
				UI_ShowDialogue(g_VMActiveMessage);
				pc += 3;
				break;
			}

			case VM_OP_GIVE_ITEM:
			{
				u16 itemID;
				if (pc + 2 >= length) return FALSE;
				itemID = (u16)bytecode[pc + 1] | ((u16)bytecode[pc + 2] << 8);
				VM_GiveItem(itemID, 1);
				UI_UpdateHUD();
				pc += 3;
				break;
			}

			case VM_OP_TAKE_ITEM:
			{
				u16 itemID;
				if (pc + 2 >= length) return FALSE;
				itemID = (u16)bytecode[pc + 1] | ((u16)bytecode[pc + 2] << 8);
				VM_TakeItem(itemID, 1);
				UI_UpdateHUD();
				pc += 3;
				break;
			}

			case VM_OP_SET_FLAG:
			{
				u8 flagID, val;
				if (pc + 2 >= length) return FALSE;
				flagID = bytecode[pc + 1];
				val    = bytecode[pc + 2];
				g_VMFlags[flagID] = val;
				pc += 3;
				break;
			}

			case VM_OP_CHECK_FLAG:
			{
				u8  flagID;
				i16 offset;
				u16 nextPC;
				if (pc + 3 >= length) return FALSE;
				flagID = bytecode[pc + 1];
				offset = (i16)((u16)bytecode[pc + 2] | ((u16)bytecode[pc + 3] << 8));
				nextPC = pc + 4;

				if (g_VMFlags[flagID] != 0)
				{
					pc = (u16)((i32)nextPC + (i32)offset);
				}
				else
				{
					pc = nextPC;
				}
				break;
			}

			case VM_OP_TELEPORT:
			{
				u16 roomID;
				u8  x, y;
				if (pc + 4 >= length) return FALSE;
				roomID = (u16)bytecode[pc + 1] | ((u16)bytecode[pc + 2] << 8);
				x      = bytecode[pc + 3];
				y      = bytecode[pc + 4];

				WORLD_LoadRoom(roomID);
				HERO_SetPosition(x, y);
				HERO_Draw();
				pc += 5;
				break;
			}

			case VM_OP_HEAL:
			{
				u8 hp;
				if (pc + 1 >= length) return FALSE;
				hp = bytecode[pc + 1];
				g_HeroStats.HP += hp;
				if (g_HeroStats.HP > g_HeroStats.MaxHP)
				{
					g_HeroStats.HP = g_HeroStats.MaxHP;
				}
				UI_UpdateHUD();
				pc += 2;
				break;
			}

			case VM_OP_DAMAGE:
			{
				u8 hp;
				if (pc + 1 >= length) return FALSE;
				hp = bytecode[pc + 1];
				if (g_HeroStats.HP <= hp)
				{
					g_HeroStats.HP = 0;
				}
				else
				{
					g_HeroStats.HP -= hp;
				}
				UI_UpdateHUD();
				pc += 2;
				break;
			}

			case VM_OP_PLAY_SFX:
			{
				u8 sfxID;
				if (pc + 1 >= length) return FALSE;
				sfxID = bytecode[pc + 1];
				VM_PlaySFX(sfxID);
				pc += 2;
				break;
			}

			case VM_OP_END:
				return TRUE;

			default:
				// Opcode desconhecido: interrompe execução defensivamente
				return FALSE;
		}
	}

	return TRUE;
}
