package tui

import (
	"fmt"
	"strconv"

	"edadv/compiler"
	"github.com/rivo/tview"
)

// MessagesTab gerencia as Mensagens (Capítulo 9).
type MessagesTab struct {
	Editor      *AdventureEditor
	View        *tview.Flex
	MsgList     *tview.List
	MsgForm     *tview.Form
	SelectedIdx int
}

func NewMessagesTab(editor *AdventureEditor) *MessagesTab {
	tab := &MessagesTab{
		Editor:      editor,
		SelectedIdx: 0,
	}

	tab.MsgList = tview.NewList().
		ShowSecondaryText(true)
	tab.MsgList.SetBorder(true).
		SetTitle("[ Mensagens do Jogo ]")

	tab.MsgForm = tview.NewForm()
	tab.MsgForm.SetBorder(true).
		SetTitle("[ Editar Mensagem ]")

	tab.View = tview.NewFlex().
		AddItem(tab.MsgList, 0, 1, true).
		AddItem(tab.MsgForm, 0, 2, false)

	tab.MsgList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		tab.SelectedIdx = index
		tab.loadMessageToForm(index)
	})

	tab.Refresh()
	return tab
}

func (tab *MessagesTab) GetView() tview.Primitive {
	return tab.View
}

func (tab *MessagesTab) Refresh() {
	tab.MsgList.Clear()
	for i, msg := range tab.Editor.Compiler.Game.Messages {
		shortcut := rune(0)
		if i < 9 {
			shortcut = rune('a' + i)
		}

		preview := msg.Text
		if len(preview) > 35 {
			preview = preview[:32] + "..."
		}

		category := "Autor"
		if msg.ID >= 11 && msg.ID <= 22 {
			category = "Sistema (Pró-definida)"
		}

		tab.MsgList.AddItem(
			fmt.Sprintf("MSG %d", msg.ID),
			fmt.Sprintf("[%s] %s", category, preview),
			shortcut,
			nil,
		)
	}

	if len(tab.Editor.Compiler.Game.Messages) > 0 {
		if tab.SelectedIdx >= len(tab.Editor.Compiler.Game.Messages) {
			tab.SelectedIdx = 0
		}
		tab.MsgList.SetCurrentItem(tab.SelectedIdx)
		tab.loadMessageToForm(tab.SelectedIdx)
	}
}

func (tab *MessagesTab) loadMessageToForm(index int) {
	tab.MsgForm.Clear(true)
	if index < 0 || index >= len(tab.Editor.Compiler.Game.Messages) {
		return
	}

	msg := &tab.Editor.Compiler.Game.Messages[index]

	tab.MsgForm.AddInputField("ID da Mensagem:", strconv.Itoa(msg.ID), 5, nil, func(text string) {
		if id, err := strconv.Atoi(text); err == nil {
			msg.ID = id
		}
	})

	tab.MsgForm.AddTextArea("Texto da Mensagem:", msg.Text, 45, 6, 0, func(text string) {
		msg.Text = text
	})

	tab.MsgForm.AddButton("Adicionar Nova Mensagem", func() {
		newID := 30
		if len(tab.Editor.Compiler.Game.Messages) > 0 {
			newID = tab.Editor.Compiler.Game.Messages[len(tab.Editor.Compiler.Game.Messages)-1].ID + 1
		}
		tab.Editor.Compiler.Game.Messages = append(tab.Editor.Compiler.Game.Messages, compiler.Message{
			ID:   newID,
			Text: "Nova mensagem de texto.",
		})
		tab.SelectedIdx = len(tab.Editor.Compiler.Game.Messages) - 1
		tab.Refresh()
	})
}
