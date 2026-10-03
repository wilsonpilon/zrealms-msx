// ____________________________
// Z-Realm (zrealm-msx) - Máquina Virtual de Eventos (MSX 2 Bytecode VM)
// Subfase 4.3: Interpretador Z80 de Opcodes, Flags Globais e Inventário
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"
#include "loader.h"
#include "world.h"
#include "hero.h"
#include "bios.h"

// Opcodes da VM de Eventos do Z-Realm
#define VM_OP_NOP        0x00 // 1 byte
#define VM_OP_MSG        0x01 // 3 bytes: string_id: u16
#define VM_OP_GIVE_ITEM  0x02 // 3 bytes: item_id: u16
#define VM_OP_TAKE_ITEM  0x03 // 3 bytes: item_id: u16
#define VM_OP_SET_FLAG   0x04 // 3 bytes: flag_id: u8, val: u8
#define VM_OP_CHECK_FLAG 0x05 // 4 bytes: flag_id: u8, jump_offset: i16
#define VM_OP_TELEPORT   0x06 // 5 bytes: room_id: u16, x: u8, y: u8
#define VM_OP_HEAL       0x07 // 2 bytes: hp: u8
#define VM_OP_DAMAGE     0x08 // 2 bytes: hp: u8
#define VM_OP_PLAY_SFX   0x09 // 2 bytes: sfx_id: u8
#define VM_OP_END        0xFF // 1 byte

// Limites da Máquina Virtual
#define VM_MAX_FLAGS        256
#define VM_MAX_INVENTORY    16
#define VM_SCRIPT_BUF_SZ    256
#define VM_MAX_MSG_LEN      128

// Slot do inventário do jogador
typedef struct
{
	u16 ItemID;
	u8  Quantity;
} VMInventorySlot;

// Atributos de RPG do Herói gerenciados pela VM
typedef struct
{
	u16 HP;
	u16 MaxHP;
	u16 MP;
	u16 MaxMP;
	u8  Level;
	u16 Gold;
} VMHeroStats;

// Estado Global da Máquina Virtual
extern u8              g_VMFlags[VM_MAX_FLAGS];
extern VMInventorySlot g_VMInventory[VM_MAX_INVENTORY];
extern u8              g_VMInventoryCount;
extern VMHeroStats     g_HeroStats;
extern c8              g_VMActiveMessage[VM_MAX_MSG_LEN];
extern u16             g_VMActiveStringID;
extern bool            g_VMHasMessage;
extern u8              g_VMLastSFX;

// Inicializa a Máquina Virtual, limpando flags, inventário e definindo stats iniciais
void VM_Init(void);

// Carrega o script especificado do Memory Mapper para a RAM local e o interpreta
bool VM_ExecuteScript(u16 scriptID);

// Executa diretamente um buffer de bytecode na RAM
bool VM_ExecuteBytecode(const u8* bytecode, u16 length);

// Manipulação de Flags Globais (0..255)
void VM_SetFlag(u8 flagID, u8 val);
u8   VM_GetFlag(u8 flagID);

// Manipulação de Inventário
bool VM_GiveItem(u16 itemID, u8 quantity);
bool VM_TakeItem(u16 itemID, u8 quantity);
bool VM_HasItem(u16 itemID);
u8   VM_GetItemQuantity(u16 itemID);

// Reproduz efeito sonoro da VM
void VM_PlaySFX(u8 sfxID);

// Limpa mensagem ativa
void VM_ClearMessage(void);
