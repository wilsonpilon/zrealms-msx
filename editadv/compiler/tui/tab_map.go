package tui

import (
	"fmt"
	"strconv"

	"edadv/compiler"
	"github.com/rivo/tview"
)

// MapTab gerencia a visualização e edição de salas e saídas.
type MapTab struct {
	Editor      *AdventureEditor
	View        *tview.Flex
	RoomList    *tview.List
	RoomForm    *tview.Form
	MapVisual   *tview.TextView
	SelectedIdx int
}

func NewMapTab(editor *AdventureEditor) *MapTab {
	tab := &MapTab{
		Editor:      editor,
		SelectedIdx: 0,
	}

	tab.RoomList = tview.NewList().
		ShowSecondaryText(true)
	tab.RoomList.SetBorder(true).
		SetTitle("[ Salas do Jogo (1..99) ]")

	tab.RoomForm = tview.NewForm()
	tab.RoomForm.SetBorder(true).
		SetTitle("[ Detalhes e Saídas Cardeais ]")

	tab.MapVisual = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	tab.MapVisual.SetBorder(true).
		SetTitle("[ Conexões da Sala ]")

	// Lado direito: Form em cima + Visual em baixo
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tab.RoomForm, 0, 3, false).
		AddItem(tab.MapVisual, 8, 0, false)

	// Layout da aba: Lista à esquerda (1/3) + Detalhes à direita (2/3)
	tab.View = tview.NewFlex().
		AddItem(tab.RoomList, 0, 1, true).
		AddItem(rightFlex, 0, 2, false)

	tab.RoomList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		tab.SelectedIdx = index
		tab.loadRoomToForm(index)
		tab.renderMapVisual(index)
	})

	tab.Refresh()
	return tab
}

func (tab *MapTab) GetView() tview.Primitive {
	return tab.View
}

func (tab *MapTab) Refresh() {
	tab.RoomList.Clear()
	for i, pos := range tab.Editor.Compiler.Game.Positions {
		shortcut := rune(0)
		if i < 9 {
			shortcut = rune('a' + i)
		}
		exitsSummary := fmt.Sprintf("N:%d S:%d L:%d O:%d NE:%d NO:%d SE:%d SO:%d",
			pos.Exits.North, pos.Exits.South, pos.Exits.East, pos.Exits.West,
			pos.Exits.Northeast, pos.Exits.Northwest, pos.Exits.Southeast, pos.Exits.Southwest)
		tab.RoomList.AddItem(
			fmt.Sprintf("[%d] %s", pos.ID, pos.Name),
			exitsSummary,
			shortcut,
			nil,
		)
	}

	if len(tab.Editor.Compiler.Game.Positions) > 0 {
		if tab.SelectedIdx >= len(tab.Editor.Compiler.Game.Positions) {
			tab.SelectedIdx = 0
		}
		tab.RoomList.SetCurrentItem(tab.SelectedIdx)
		tab.loadRoomToForm(tab.SelectedIdx)
		tab.renderMapVisual(tab.SelectedIdx)
	}
}

func (tab *MapTab) loadRoomToForm(index int) {
	tab.RoomForm.Clear(true)
	if index < 0 || index >= len(tab.Editor.Compiler.Game.Positions) {
		return
	}

	pos := &tab.Editor.Compiler.Game.Positions[index]

	tab.RoomForm.AddInputField("ID da Posição:", strconv.Itoa(pos.ID), 5, nil, func(text string) {
		if id, err := strconv.Atoi(text); err == nil {
			pos.ID = id
		}
	})

	tab.RoomForm.AddInputField("Nome:", pos.Name, 30, nil, func(text string) {
		pos.Name = text
	})

	tab.RoomForm.AddTextArea("Descrição:", pos.Description, 40, 4, 0, func(text string) {
		pos.Description = text
	})

	tab.RoomForm.AddInputField("Norte (N):", strconv.Itoa(pos.Exits.North), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.North = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Sul (S):", strconv.Itoa(pos.Exits.South), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.South = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Leste (L/E):", strconv.Itoa(pos.Exits.East), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.East = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Oeste (O/W):", strconv.Itoa(pos.Exits.West), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.West = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Nordeste (NE):", strconv.Itoa(pos.Exits.Northeast), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.Northeast = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Noroeste (NO):", strconv.Itoa(pos.Exits.Northwest), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.Northwest = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Sudeste (SE):", strconv.Itoa(pos.Exits.Southeast), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.Southeast = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddInputField("Sudoeste (SO):", strconv.Itoa(pos.Exits.Southwest), 6, nil, func(text string) {
		if v, err := strconv.Atoi(text); err == nil {
			pos.Exits.Southwest = v
			tab.renderMapVisual(index)
		}
	})

	tab.RoomForm.AddButton("Adicionar Nova Sala", func() {
		newID := len(tab.Editor.Compiler.Game.Positions) + 1
		tab.Editor.Compiler.Game.Positions = append(tab.Editor.Compiler.Game.Positions, compiler.Position{
			ID:          newID,
			Name:        fmt.Sprintf("Nova Sala %d", newID),
			Description: "Descreva a nova sala aqui.",
		})
		tab.SelectedIdx = len(tab.Editor.Compiler.Game.Positions) - 1
		tab.Refresh()
	})
}

func (tab *MapTab) renderMapVisual(index int) {
	if index < 0 || index >= len(tab.Editor.Compiler.Game.Positions) {
		tab.MapVisual.SetText("")
		return
	}

	pos := tab.Editor.Compiler.Game.Positions[index]

	formatExit := func(val int) string {
		if val == 0 {
			return "[red]Bloqueado[-]"
		} else if val <= compiler.MaxPositions {
			return fmt.Sprintf("[green]Sala %d[-]", val)
		}
		return fmt.Sprintf("[yellow]Fn %d[-]", val-compiler.PositionCondOffset)
	}

	n := formatExit(pos.Exits.North)
	s := formatExit(pos.Exits.South)
	l := formatExit(pos.Exits.East)
	o := formatExit(pos.Exits.West)
	ne := formatExit(pos.Exits.Northeast)
	no := formatExit(pos.Exits.Northwest)
	se := formatExit(pos.Exits.Southeast)
	so := formatExit(pos.Exits.Southwest)

	visual := fmt.Sprintf(
		"\n        NO: %s   ▲ N: %s   ▲ NE: %s\n"+
			"                   ╲     │     ╱\n"+
			"         O: %s ◄───[yellow] Sala %d [white]───► L: %s\n"+
			"                   ╱     │     ╲\n"+
			"        SO: %s   ▼ S: %s   ▼ SE: %s\n",
		no, n, ne, o, pos.ID, l, so, s, se,
	)

	tab.MapVisual.SetText(visual)
}
