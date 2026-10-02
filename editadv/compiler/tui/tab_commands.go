package tui

import (
	"fmt"
	"strconv"

	"edadv/compiler"
	"github.com/rivo/tview"
)

// CommandsTab gerencia os Comandos e Funções em Bytecode (Capítulos 7, 8 e 10).
type CommandsTab struct {
	Editor      *AdventureEditor
	View        *tview.Flex
	ProcList    *tview.List
	InstTable   *tview.Table
	HelpView    *tview.TextView
	IsCommand   bool // true = Comandos, false = Funções
	SelectedIdx int
}

func NewCommandsTab(editor *AdventureEditor) *CommandsTab {
	tab := &CommandsTab{
		Editor:      editor,
		IsCommand:   true,
		SelectedIdx: 0,
	}

	tab.ProcList = tview.NewList().
		ShowSecondaryText(true)
	tab.ProcList.SetBorder(true).
		SetTitle("[ Procedimentos ]")

	tab.InstTable = tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	tab.InstTable.SetBorder(true).
		SetTitle("[ Instruções Bytecode (Capítulo 8) ]")

	tab.HelpView = tview.NewTextView().
		SetDynamicColors(true)
	tab.HelpView.SetBorder(true).
		SetTitle("[ Significado da Instrução ]")

	// Lado direito: Tabela de instruções em cima + Explicação do bytecode em baixo
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tab.InstTable, 0, 2, true).
		AddItem(tab.HelpView, 6, 0, false)

	tab.View = tview.NewFlex().
		AddItem(tab.ProcList, 0, 1, false).
		AddItem(rightFlex, 0, 2, true)

	tab.ProcList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		tab.SelectedIdx = index
		tab.loadInstructions(index)
	})

	tab.InstTable.SetSelectionChangedFunc(func(row, column int) {
		tab.updateInstructionHelp(row)
	})

	tab.Refresh()
	return tab
}

func (tab *CommandsTab) GetView() tview.Primitive {
	return tab.View
}

func (tab *CommandsTab) Refresh() {
	tab.ProcList.Clear()

	// Alternador de modo no topo da lista
	modeTitle := "[yellow]Comandos[-] (Frases do Jogador)"
	if !tab.IsCommand {
		modeTitle = "[yellow]Funções[-] (Eventos do Sistema)"
	}
	tab.ProcList.AddItem(">>> Alternar Comandos / Funções <<<", modeTitle, 't', func() {
		tab.IsCommand = !tab.IsCommand
		tab.SelectedIdx = 0
		tab.Refresh()
	})

	if tab.IsCommand {
		for i, cmd := range tab.Editor.Compiler.Game.Commands {
			shortcut := rune(0)
			if i < 9 {
				shortcut = rune('1' + i)
			}
			comment := cmd.Comment
			if comment == "" {
				comment = fmt.Sprintf("%d instruções", len(cmd.Instructions))
			}
			tab.ProcList.AddItem(
				fmt.Sprintf("Cmd %d: V:%d O1:%d O2:%d", i+1, cmd.Verb, cmd.Object1, cmd.Object2),
				comment,
				shortcut,
				nil,
			)
		}
	} else {
		for i, fn := range tab.Editor.Compiler.Game.Functions {
			shortcut := rune(0)
			if i < 9 {
				shortcut = rune('1' + i)
			}
			desc := fn.Name
			if desc == "" {
				desc = fmt.Sprintf("%d instruções", len(fn.Instructions))
			}
			tab.ProcList.AddItem(
				fmt.Sprintf("Função %d", fn.ID),
				desc,
				shortcut,
				nil,
			)
		}
	}

	if tab.SelectedIdx > 0 {
		tab.ProcList.SetCurrentItem(tab.SelectedIdx)
		tab.loadInstructions(tab.SelectedIdx)
	} else {
		tab.loadInstructions(1)
	}
}

func (tab *CommandsTab) loadInstructions(index int) {
	tab.InstTable.Clear()

	// Cabeçalho da tabela
	headers := []string{"Passo", "Label", "Opcode", "P1", "P2", "P3"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tview.Styles.TitleColor).
			SetSelectable(false).
			SetAlign(tview.AlignCenter)
		tab.InstTable.SetCell(0, col, cell)
	}

	var instructions []compiler.Instruction
	if tab.IsCommand {
		cmdIdx := index - 1
		if cmdIdx >= 0 && cmdIdx < len(tab.Editor.Compiler.Game.Commands) {
			instructions = tab.Editor.Compiler.Game.Commands[cmdIdx].Instructions
		}
	} else {
		fnIdx := index - 1
		if fnIdx >= 0 && fnIdx < len(tab.Editor.Compiler.Game.Functions) {
			instructions = tab.Editor.Compiler.Game.Functions[fnIdx].Instructions
		}
	}

	for row, inst := range instructions {
		tab.InstTable.SetCell(row+1, 0, tview.NewTableCell(strconv.Itoa(row)).SetAlign(tview.AlignCenter))
		tab.InstTable.SetCell(row+1, 1, tview.NewTableCell(inst.Label).SetAlign(tview.AlignCenter))
		tab.InstTable.SetCell(row+1, 2, tview.NewTableCell(inst.Op).SetTextColor(tview.Styles.SecondaryTextColor))
		tab.InstTable.SetCell(row+1, 3, tview.NewTableCell(strconv.Itoa(inst.P1)).SetAlign(tview.AlignRight))
		tab.InstTable.SetCell(row+1, 4, tview.NewTableCell(strconv.Itoa(inst.P2)).SetAlign(tview.AlignRight))
		tab.InstTable.SetCell(row+1, 5, tview.NewTableCell(strconv.Itoa(inst.P3)).SetAlign(tview.AlignRight))
	}

	tab.updateInstructionHelp(1)
}

func (tab *CommandsTab) updateInstructionHelp(row int) {
	if row <= 0 || row >= tab.InstTable.GetRowCount() {
		tab.HelpView.SetText("Selecione uma instrução para ver a explicação de funcionamento.")
		return
	}

	opCell := tab.InstTable.GetCell(row, 2)
	if opCell == nil {
		return
	}

	op := opCell.Text
	p1Cell := tab.InstTable.GetCell(row, 3)
	p2Cell := tab.InstTable.GetCell(row, 4)
	p3Cell := tab.InstTable.GetCell(row, 5)

	p1, _ := strconv.Atoi(p1Cell.Text)
	p2, _ := strconv.Atoi(p2Cell.Text)
	p3, _ := strconv.Atoi(p3Cell.Text)

	var desc string
	switch op {
	case "MSG":
		desc = fmt.Sprintf("Imprime a mensagem número %d na tela central.", p1)
	case "PEGA":
		desc = fmt.Sprintf("Coloca o objeto %d na mão (inventário). Se 0, usa o objeto em evidência.", p1)
	case "SOLTA":
		desc = fmt.Sprintf("Solta o objeto %d na sala atual. Se 100, solta todos os objetos carregados.", p1)
	case "POE":
		desc = fmt.Sprintf("Guarda o objeto %d dentro do Objeto 3 (recipiente).", p1)
	case "NVC":
		desc = "Encerra o comando atual e aguarda nova entrada de texto do jogador."
	case "LDR":
		desc = fmt.Sprintf("Carrega o registrador %d com o valor %d.", p1, p2)
	case "TEMOS":
		desc = fmt.Sprintf("Se o jogador possui o objeto %d, salta para o passo %d.", p1, p2)
	case "AQUI":
		desc = fmt.Sprintf("Se o objeto %d está presente na sala atual, salta para o passo %d.", p1, p2)
	case "LOCAL":
		desc = fmt.Sprintf("Se a posição atual do jogador for %d, salta para o passo %d.", p1, p2)
	case "REG=":
		desc = fmt.Sprintf("Se Reg[%d] == %d, salta para o passo %d; senão continua.", p1, p2, p3)
	case "REG>":
		desc = fmt.Sprintf("Se Reg[%d] > %d, salta para o passo %d; senão continua.", p1, p2, p3)
	case "REG<":
		desc = fmt.Sprintf("Se Reg[%d] < %d, salta para o passo %d; senão continua.", p1, p2, p3)
	case "GOTO":
		desc = fmt.Sprintf("Desvio incondicional para o passo %d.", p1)
	case "GOSUB":
		desc = fmt.Sprintf("Chama a Função %d como subrotina (1 nível).", p1)
	case "RET":
		desc = "Retorna da subrotina GOSUB."
	case "OK":
		desc = "Imprime 'Ok.' na tela e aguarda novo comando."
	case "DESC":
		desc = "Descreve a posição atual e aguarda novo comando."
	case "FIM":
		desc = "Finaliza a partida."
	case "NEU":
		desc = "Reinicia a partida do zero (reset)."
	case "CLS":
		desc = "Limpa o Campo Central da tela."
	default:
		desc = fmt.Sprintf("Instrução %s (P1:%d, P2:%d, P3:%d)", op, p1, p2, p3)
	}

	tab.HelpView.SetText(fmt.Sprintf("[yellow]%s[-] %d %d %d:\n%s", op, p1, p2, p3, desc))
}
