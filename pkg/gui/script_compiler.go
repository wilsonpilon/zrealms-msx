package gui

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// Opcodes da Máquina Virtual de Eventos (MSX 2 Bytecode VM)
const (
	OpNop       byte = 0x00 // Nenhuma operação
	OpMsg       byte = 0x01 // Exibe diálogo da tabela de textos (uint16 string_id)
	OpGiveItem  byte = 0x02 // Adiciona item ao inventário (uint16 item_id)
	OpTakeItem  byte = 0x03 // Remove item do inventário (uint16 item_id)
	OpSetFlag   byte = 0x04 // Altera estado global de flag (uint8 flag_id, uint8 val)
	OpCheckFlag byte = 0x05 // Desvio condicional se flag != 0 (uint8 flag_id, int16 offset)
	OpTeleport  byte = 0x06 // Teletransporte forçado (uint16 room_id, uint8 x, uint8 y)
	OpHeal      byte = 0x07 // Cura HP do herói (uint8 hp)
	OpDamage    byte = 0x08 // Causa dano no herói (uint8 hp)
	OpPlaySFX   byte = 0x09 // Toca efeito sonoro PSG (uint8 sfx_id)
	OpEnd       byte = 0xFF // Fim da execução do script
)

// OpcodeInfo armazena metadados de uma instrução
type OpcodeInfo struct {
	Opcode      byte
	Mnemonic    string
	ArgCount    int
	Description string
	Example     string
}

// OpcodeCatalog disponibiliza o catálogo de instruções para a GUI
var OpcodeCatalog = []OpcodeInfo{
	{Opcode: OpMsg, Mnemonic: "MSG", ArgCount: 1, Description: "Exibe diálogo na caixa de texto", Example: "MSG 1"},
	{Opcode: OpGiveItem, Mnemonic: "GIVE_ITEM", ArgCount: 1, Description: "Adiciona item ao inventário", Example: "GIVE_ITEM 2"},
	{Opcode: OpTakeItem, Mnemonic: "TAKE_ITEM", ArgCount: 1, Description: "Remove item do inventário", Example: "TAKE_ITEM 2"},
	{Opcode: OpSetFlag, Mnemonic: "SET_FLAG", ArgCount: 2, Description: "Define valor de uma flag (0 ou 1)", Example: "SET_FLAG 10 1"},
	{Opcode: OpCheckFlag, Mnemonic: "CHECK_FLAG", ArgCount: 2, Description: "Pula se a flag estiver ativa", Example: "CHECK_FLAG 10 label_fim"},
	{Opcode: OpTeleport, Mnemonic: "TELEPORT", ArgCount: 3, Description: "Move herói para sala (X, Y)", Example: "TELEPORT 3 16 9"},
	{Opcode: OpHeal, Mnemonic: "HEAL", ArgCount: 1, Description: "Cura pontos de vida do herói", Example: "HEAL 20"},
	{Opcode: OpDamage, Mnemonic: "DAMAGE", ArgCount: 1, Description: "Causa dano ao herói", Example: "DAMAGE 10"},
	{Opcode: OpPlaySFX, Mnemonic: "PLAY_SFX", ArgCount: 1, Description: "Reproduz efeito sonoro no PSG", Example: "PLAY_SFX 3"},
	{Opcode: OpEnd, Mnemonic: "END", ArgCount: 0, Description: "Finaliza execução do evento", Example: "END"},
	{Opcode: OpNop, Mnemonic: "NOP", ArgCount: 0, Description: "Nenhuma operação", Example: "NOP"},
}

// CompileScript compila código fonte textual em bytecode binário executável pelo MSX.
func CompileScript(source string) ([]byte, error) {
	lines := strings.Split(source, "\n")
	type parsedLine struct {
		lineNum  int
		label    string
		mnemonic string
		args     []string
	}

	var parsed []parsedLine
	labels := make(map[string]int) // label -> byte offset

	// 1º Passo: Tokenização e cálculo de endereços de labels
	currentOffset := 0
	for idx, rawLine := range lines {
		lineNum := idx + 1
		line := strings.TrimSpace(rawLine)

		// Remove comentários (# ou //)
		if commentIdx := strings.Index(line, "//"); commentIdx >= 0 {
			line = strings.TrimSpace(line[:commentIdx])
		}
		if commentIdx := strings.Index(line, "#"); commentIdx >= 0 {
			line = strings.TrimSpace(line[:commentIdx])
		}
		if line == "" {
			continue
		}

		// Verifica se é uma definição de Label (ex: "label_fim:")
		if strings.HasSuffix(line, ":") {
			lbl := strings.TrimSuffix(line, ":")
			lbl = strings.TrimSpace(lbl)
			if lbl == "" {
				return nil, fmt.Errorf("linha %d: label vazio inválido", lineNum)
			}
			if _, exists := labels[lbl]; exists {
				return nil, fmt.Errorf("linha %d: label '%s' duplicado", lineNum, lbl)
			}
			labels[lbl] = currentOffset
			continue
		}

		// Divide instrução e argumentos
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		mnemonic := strings.ToUpper(tokens[0])
		args := tokens[1:]

		// Estima tamanho em bytes da instrução
		instrSize := 1
		switch mnemonic {
		case "NOP", "END":
			instrSize = 1
		case "MSG", "GIVE_ITEM", "TAKE_ITEM":
			instrSize = 3 // opcode(1) + uint16(2)
		case "SET_FLAG":
			instrSize = 3 // opcode(1) + flag(1) + val(1)
		case "CHECK_FLAG":
			instrSize = 4 // opcode(1) + flag(1) + int16(2)
		case "TELEPORT":
			instrSize = 5 // opcode(1) + room(2) + x(1) + y(1)
		case "HEAL", "DAMAGE", "PLAY_SFX":
			instrSize = 2 // opcode(1) + val(1)
		default:
			return nil, fmt.Errorf("linha %d: mnemônico desconhecido '%s'", lineNum, mnemonic)
		}

		parsed = append(parsed, parsedLine{
			lineNum:  lineNum,
			mnemonic: mnemonic,
			args:     args,
		})
		currentOffset += instrSize
	}

	// Se o script não terminar com END, adiciona automaticamente
	if len(parsed) == 0 || parsed[len(parsed)-1].mnemonic != "END" {
		parsed = append(parsed, parsedLine{
			lineNum:  len(lines) + 1,
			mnemonic: "END",
			args:     nil,
		})
	}

	// 2º Passo: Emissão do Bytecode com resolução de parâmetros e labels
	bytecode := make([]byte, 0, currentOffset+1)

	for _, p := range parsed {
		startInstrOffset := len(bytecode)

		switch p.mnemonic {
		case "NOP":
			bytecode = append(bytecode, OpNop)

		case "END":
			bytecode = append(bytecode, OpEnd)

		case "MSG":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: MSG requer 1 argumento (string_id)", p.lineNum)
			}
			val, err := strconv.ParseUint(p.args[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("linha %d: ID de string inválido '%s'", p.lineNum, p.args[0])
			}
			bytecode = append(bytecode, OpMsg, byte(val&0xFF), byte((val>>8)&0xFF))

		case "GIVE_ITEM":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: GIVE_ITEM requer 1 argumento (item_id)", p.lineNum)
			}
			val, err := strconv.ParseUint(p.args[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("linha %d: ID de item inválido '%s'", p.lineNum, p.args[0])
			}
			bytecode = append(bytecode, OpGiveItem, byte(val&0xFF), byte((val>>8)&0xFF))

		case "TAKE_ITEM":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: TAKE_ITEM requer 1 argumento (item_id)", p.lineNum)
			}
			val, err := strconv.ParseUint(p.args[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("linha %d: ID de item inválido '%s'", p.lineNum, p.args[0])
			}
			bytecode = append(bytecode, OpTakeItem, byte(val&0xFF), byte((val>>8)&0xFF))

		case "SET_FLAG":
			if len(p.args) < 2 {
				return nil, fmt.Errorf("linha %d: SET_FLAG requer 2 argumentos (flag_id, valor)", p.lineNum)
			}
			flag, errF := strconv.ParseUint(p.args[0], 10, 8)
			val, errV := strconv.ParseUint(p.args[1], 10, 8)
			if errF != nil || errV != nil {
				return nil, fmt.Errorf("linha %d: argumentos inválidos para SET_FLAG", p.lineNum)
			}
			bytecode = append(bytecode, OpSetFlag, byte(flag), byte(val))

		case "CHECK_FLAG":
			if len(p.args) < 2 {
				return nil, fmt.Errorf("linha %d: CHECK_FLAG requer 2 argumentos (flag_id, label_destino)", p.lineNum)
			}
			flag, errF := strconv.ParseUint(p.args[0], 10, 8)
			if errF != nil {
				return nil, fmt.Errorf("linha %d: flag inválida '%s'", p.lineNum, p.args[0])
			}

			// Tenta resolver label ou número direto de offset
			var jumpOffset int16
			targetLabel := p.args[1]
			if targetPos, ok := labels[targetLabel]; ok {
				// Offset relativo a partir do final desta instrução (startInstrOffset + 4)
				nextInstrPos := startInstrOffset + 4
				jumpOffset = int16(targetPos - nextInstrPos)
			} else {
				val, errO := strconv.ParseInt(targetLabel, 10, 16)
				if errO != nil {
					return nil, fmt.Errorf("linha %d: label ou offset '%s' não encontrado", p.lineNum, targetLabel)
				}
				jumpOffset = int16(val)
			}
			bytecode = append(bytecode, OpCheckFlag, byte(flag), byte(jumpOffset&0xFF), byte((jumpOffset>>8)&0xFF))

		case "TELEPORT":
			if len(p.args) < 3 {
				return nil, fmt.Errorf("linha %d: TELEPORT requer 3 argumentos (room_id, pos_x, pos_y)", p.lineNum)
			}
			roomID, errR := strconv.ParseUint(p.args[0], 10, 16)
			posX, errX := strconv.ParseUint(p.args[1], 10, 8)
			posY, errY := strconv.ParseUint(p.args[2], 10, 8)
			if errR != nil || errX != nil || errY != nil || posX >= 32 || posY >= 18 {
				return nil, fmt.Errorf("linha %d: argumentos inválidos para TELEPORT (X: 0..31, Y: 0..17)", p.lineNum)
			}
			bytecode = append(bytecode, OpTeleport, byte(roomID&0xFF), byte((roomID>>8)&0xFF), byte(posX), byte(posY))

		case "HEAL":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: HEAL requer 1 argumento (hp)", p.lineNum)
			}
			hp, err := strconv.ParseUint(p.args[0], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("linha %d: valor inválido para HEAL", p.lineNum)
			}
			bytecode = append(bytecode, OpHeal, byte(hp))

		case "DAMAGE":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: DAMAGE requer 1 argumento (hp)", p.lineNum)
			}
			hp, err := strconv.ParseUint(p.args[0], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("linha %d: valor inválido para DAMAGE", p.lineNum)
			}
			bytecode = append(bytecode, OpDamage, byte(hp))

		case "PLAY_SFX":
			if len(p.args) < 1 {
				return nil, fmt.Errorf("linha %d: PLAY_SFX requer 1 argumento (sfx_id)", p.lineNum)
			}
			sfx, err := strconv.ParseUint(p.args[0], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("linha %d: valor inválido para PLAY_SFX", p.lineNum)
			}
			bytecode = append(bytecode, OpPlaySFX, byte(sfx))
		}
	}

	return bytecode, nil
}

// DisassembleScript converte bytecode binário de volta para mnemônicos textuais legíveis.
func DisassembleScript(bytecode []byte) string {
	if len(bytecode) == 0 {
		return ""
	}

	var sb strings.Builder
	idx := 0
	for idx < len(bytecode) {
		op := bytecode[idx]
		idx++

		switch op {
		case OpNop:
			sb.WriteString("NOP\n")

		case OpEnd:
			sb.WriteString("END\n")

		case OpMsg:
			if idx+2 <= len(bytecode) {
				val := binary.LittleEndian.Uint16(bytecode[idx : idx+2])
				idx += 2
				sb.WriteString(fmt.Sprintf("MSG %d\n", val))
			}

		case OpGiveItem:
			if idx+2 <= len(bytecode) {
				val := binary.LittleEndian.Uint16(bytecode[idx : idx+2])
				idx += 2
				sb.WriteString(fmt.Sprintf("GIVE_ITEM %d\n", val))
			}

		case OpTakeItem:
			if idx+2 <= len(bytecode) {
				val := binary.LittleEndian.Uint16(bytecode[idx : idx+2])
				idx += 2
				sb.WriteString(fmt.Sprintf("TAKE_ITEM %d\n", val))
			}

		case OpSetFlag:
			if idx+2 <= len(bytecode) {
				flag := bytecode[idx]
				val := bytecode[idx+1]
				idx += 2
				sb.WriteString(fmt.Sprintf("SET_FLAG %d %d\n", flag, val))
			}

		case OpCheckFlag:
			if idx+3 <= len(bytecode) {
				flag := bytecode[idx]
				offset := int16(binary.LittleEndian.Uint16(bytecode[idx+1 : idx+3]))
				idx += 3
				sb.WriteString(fmt.Sprintf("CHECK_FLAG %d %d\n", flag, offset))
			}

		case OpTeleport:
			if idx+4 <= len(bytecode) {
				roomID := binary.LittleEndian.Uint16(bytecode[idx : idx+2])
				x := bytecode[idx+2]
				y := bytecode[idx+3]
				idx += 4
				sb.WriteString(fmt.Sprintf("TELEPORT %d %d %d\n", roomID, x, y))
			}

		case OpHeal:
			if idx < len(bytecode) {
				val := bytecode[idx]
				idx++
				sb.WriteString(fmt.Sprintf("HEAL %d\n", val))
			}

		case OpDamage:
			if idx < len(bytecode) {
				val := bytecode[idx]
				idx++
				sb.WriteString(fmt.Sprintf("DAMAGE %d\n", val))
			}

		case OpPlaySFX:
			if idx < len(bytecode) {
				val := bytecode[idx]
				idx++
				sb.WriteString(fmt.Sprintf("PLAY_SFX %d\n", val))
			}

		default:
			sb.WriteString(fmt.Sprintf("// Byte Desconhecido 0x%02X\n", op))
		}
	}

	return strings.TrimSpace(sb.String())
}

// FormatBytecodeHex retorna representação em string hexadecimal dos bytes (ex: "01 0A 00 FF").
func FormatBytecodeHex(bytecode []byte) string {
	if len(bytecode) == 0 {
		return "(Vazio)"
	}
	var parts []string
	for _, b := range bytecode {
		parts = append(parts, fmt.Sprintf("%02X", b))
	}
	return strings.Join(parts, " ")
}
