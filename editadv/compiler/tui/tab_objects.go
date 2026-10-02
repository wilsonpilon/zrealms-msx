package tui

import (
	"fmt"
	"strconv"
	"strings"

	"edadv/compiler"
	"github.com/rivo/tview"
)

// ObjectsTab gerencia os objetos e o Byte de Consistência (Capítulo 5).
type ObjectsTab struct {
	Editor      *AdventureEditor
	View        *tview.Flex
	ObjectList  *tview.List
	ObjectForm  *tview.Form
	SelectedIdx int
}

func NewObjectsTab(editor *AdventureEditor) *ObjectsTab {
	tab := &ObjectsTab{
		Editor:      editor,
		SelectedIdx: 0,
	}

	tab.ObjectList = tview.NewList().
		ShowSecondaryText(true)
	tab.ObjectList.SetBorder(true).
		SetTitle("[ Objetos do Jogo (1..99) ]")

	tab.ObjectForm = tview.NewForm()
	tab.ObjectForm.SetBorder(true).
		SetTitle("[ Detalhes e Byte de Consistência ]")

	tab.View = tview.NewFlex().
		AddItem(tab.ObjectList, 0, 1, true).
		AddItem(tab.ObjectForm, 0, 2, false)

	tab.ObjectList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		tab.SelectedIdx = index
		tab.loadObjectToForm(index)
	})

	tab.Refresh()
	return tab
}

func (tab *ObjectsTab) GetView() tview.Primitive {
	return tab.View
}

func (tab *ObjectsTab) Refresh() {
	tab.ObjectList.Clear()
	for i, obj := range tab.Editor.Compiler.Game.Objects {
		shortcut := rune(0)
		if i < 9 {
			shortcut = rune('a' + i)
		}

		sitDesc := "Não existe"
		if obj.InitialSituation > 0 && obj.InitialSituation <= compiler.MaxPositions {
			sitDesc = fmt.Sprintf("Na Sala %d", obj.InitialSituation)
		} else if obj.InitialSituation == compiler.ObjSitCarregado {
			sitDesc = "Na mão (inventário)"
		} else if obj.InitialSituation == compiler.ObjSitEmObj3Aberto {
			sitDesc = "No Objeto 3 (aberto)"
		} else if obj.InitialSituation == compiler.ObjSitEmObj3Fechado {
			sitDesc = "No Objeto 3 (fechado)"
		}

		tab.ObjectList.AddItem(
			fmt.Sprintf("[%d] %s", obj.ID, obj.Name),
			fmt.Sprintf("Situação: %s | Byte Consistência: 0x%02X", sitDesc, obj.Consistency.ToByte()),
			shortcut,
			nil,
		)
	}

	if len(tab.Editor.Compiler.Game.Objects) > 0 {
		if tab.SelectedIdx >= len(tab.Editor.Compiler.Game.Objects) {
			tab.SelectedIdx = 0
		}
		tab.ObjectList.SetCurrentItem(tab.SelectedIdx)
		tab.loadObjectToForm(tab.SelectedIdx)
	}
}

func (tab *ObjectsTab) loadObjectToForm(index int) {
	tab.ObjectForm.Clear(true)
	if index < 0 || index >= len(tab.Editor.Compiler.Game.Objects) {
		return
	}

	obj := &tab.Editor.Compiler.Game.Objects[index]

	tab.ObjectForm.AddInputField("ID do Objeto:", strconv.Itoa(obj.ID), 5, nil, func(text string) {
		if id, err := strconv.Atoi(text); err == nil {
			obj.ID = id
		}
	})

	tab.ObjectForm.AddInputField("Nome Principal:", obj.Name, 25, nil, func(text string) {
		obj.Name = text
	})

	synonymsStr := strings.Join(obj.Synonyms, "/")
	tab.ObjectForm.AddInputField("Sinônimos (sep. por /):", synonymsStr, 30, nil, func(text string) {
		if text == "" {
			obj.Synonyms = nil
		} else {
			obj.Synonyms = strings.Split(text, "/")
		}
	})

	tab.ObjectForm.AddInputField("Situação Inicial (0..253):", strconv.Itoa(int(obj.InitialSituation)), 5, nil, func(text string) {
		if sit, err := strconv.Atoi(text); err == nil {
			obj.InitialSituation = uint8(sit)
		}
	})

	tab.ObjectForm.AddTextArea("Descrição:", obj.Description, 40, 3, 0, func(text string) {
		obj.Description = text
	})

	// Checkboxes do Byte de Consistência (Capítulo 5.3)
	tab.ObjectForm.AddCheckbox("Bit 0: Pode Pegar (Função 6)", obj.Consistency.CanTake, func(checked bool) {
		obj.Consistency.CanTake = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 1: Pode Guardar no Obj 3 (Função 7)", obj.Consistency.CanPutInObj3, func(checked bool) {
		obj.Consistency.CanPutInObj3 = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 2: Pode Trocar (Função 8)", obj.Consistency.CanTrade, func(checked bool) {
		obj.Consistency.CanTrade = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 3: Pode Comprar (Função 9)", obj.Consistency.CanBuy, func(checked bool) {
		obj.Consistency.CanBuy = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 4: Pode Roubar (Função 10)", obj.Consistency.CanSteal, func(checked bool) {
		obj.Consistency.CanSteal = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 5: Pode Tirar (Função 11)", obj.Consistency.CanRemove, func(checked bool) {
		obj.Consistency.CanRemove = checked
	})

	tab.ObjectForm.AddCheckbox("Bit 6: Pode Quebrar (Função 12)", obj.Consistency.CanBreak, func(checked bool) {
		obj.Consistency.CanBreak = checked
	})

	tab.ObjectForm.AddButton("Adicionar Novo Objeto", func() {
		newID := len(tab.Editor.Compiler.Game.Objects) + 1
		tab.Editor.Compiler.Game.Objects = append(tab.Editor.Compiler.Game.Objects, compiler.Object{
			ID:               newID,
			Name:             fmt.Sprintf("NOVO_ITEM_%d", newID),
			InitialSituation: 1,
			Consistency: compiler.ConsistencyConfig{
				CanTake: true,
			},
			Description: "Descrição do novo objeto.",
		})
		tab.SelectedIdx = len(tab.Editor.Compiler.Game.Objects) - 1
		tab.Refresh()
	})
}
