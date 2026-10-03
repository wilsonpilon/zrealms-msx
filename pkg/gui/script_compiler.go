package gui

import (
	"github.com/zrealm-msx/zrealm/pkg/script"
)

// Opcodes da Máquina Virtual de Eventos (MSX 2 Bytecode VM)
const (
	OpNop       = script.OpNop
	OpMsg       = script.OpMsg
	OpGiveItem  = script.OpGiveItem
	OpTakeItem  = script.OpTakeItem
	OpSetFlag   = script.OpSetFlag
	OpCheckFlag = script.OpCheckFlag
	OpTeleport  = script.OpTeleport
	OpHeal      = script.OpHeal
	OpDamage    = script.OpDamage
	OpPlaySFX   = script.OpPlaySFX
	OpEnd       = script.OpEnd
)

// OpcodeInfo armazena metadados de uma instrução
type OpcodeInfo = script.OpcodeInfo

// OpcodeCatalog disponibiliza o catálogo de instruções para a GUI
var OpcodeCatalog = script.OpcodeCatalog

// CompileScript compila código fonte textual em bytecode binário executável pelo MSX.
var CompileScript = script.CompileScript

// DisassembleScript converte bytecode binário de volta para mnemônicos textuais legíveis.
var DisassembleScript = script.DisassembleScript

// FormatBytecodeHex retorna representação em string hexadecimal dos bytes.
var FormatBytecodeHex = script.FormatBytecodeHex
