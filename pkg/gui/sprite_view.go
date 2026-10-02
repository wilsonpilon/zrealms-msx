package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

// SpriteView implementa a aba de edição interativa de sprites Modo 2 do MSX 2 (16x16 pixels).
type SpriteView struct {
	state          *ProjectState
	win            fyne.Window
	spriteList     *widget.List
	sprites        []*models.Sprite
	selectedSprite *models.Sprite
	editor         *SpriteEditorWidget
}

// NewSpriteView cria o componente de visualização e edição de sprites.
func NewSpriteView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &SpriteView{
		state:   state,
		win:     win,
		sprites: make([]*models.Sprite, 0),
	}

	// Lista de Sprites Cadastrados
	v.spriteList = widget.NewList(
		func() int {
			return len(v.sprites)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.RadioButtonCheckedIcon()),
				widget.NewLabel("Sprite"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.sprites) {
				return
			}
			box := obj.(*fyne.Container)
			lbl := box.Objects[1].(*widget.Label)
			s := v.sprites[id]
			lbl.SetText(fmt.Sprintf("#%d: %s", s.ID, s.Name))
		},
	)

	v.spriteList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.sprites) {
			return
		}
		v.selectedSprite = v.sprites[id]
		v.editor.LoadSprite(v.selectedSprite)
	}

	// Editor Interativo
	v.editor = NewSpriteEditorWidget(state, win, func() {
		v.reload()
	})

	// Botões de Ação na Lista
	btnAdd := widget.NewButtonWithIcon("Novo Sprite", theme.ContentAddIcon(), func() {
		v.createSprite()
	})

	btnDelete := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedSprite()
	})

	header := container.NewVBox(
		widget.NewLabelWithStyle("👾 Sprites Modo 2 (16x16)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAdd, btnDelete),
		widget.NewSeparator(),
	)
	listPanel := container.NewBorder(header, nil, nil, nil, v.spriteList)

	split := container.NewHSplit(listPanel, v.editor.Container())
	split.Offset = 0.32

	state.OnProjectLoaded(func(_ *project.Project, _ string) {
		v.reload()
	})
	state.OnProjectClosed(func() {
		v.sprites = make([]*models.Sprite, 0)
		v.selectedSprite = nil
		v.spriteList.Refresh()
		v.editor.LoadSprite(nil)
	})

	if state.IsOpen() {
		v.reload()
	}

	return split
}

func (v *SpriteView) reload() {
	proj := v.state.Current()
	if proj == nil {
		return
	}
	sprites, err := proj.Storage().Sprites.ListSprites(context.Background())
	if err != nil {
		dialog.ShowError(err, v.win)
		return
	}
	v.sprites = sprites
	v.spriteList.Refresh()

	if len(v.sprites) > 0 {
		// Mantém seleção se ainda existir
		selectedIndex := 0
		if v.selectedSprite != nil {
			for i, s := range v.sprites {
				if s.ID == v.selectedSprite.ID {
					selectedIndex = i
					break
				}
			}
		}
		v.spriteList.Select(selectedIndex)
	} else {
		v.selectedSprite = nil
		v.editor.LoadSprite(nil)
	}
}

func (v *SpriteView) createSprite() {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("Hero_Walk_%d", len(v.sprites)))

	items := []*widget.FormItem{
		widget.NewFormItem("Nome do Sprite", nameEntry),
	}

	dialog.ShowForm("Novo Sprite 16x16 (Modo 2)", "Criar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		s := models.NewSprite(nameEntry.Text)
		if err := proj.Storage().Sprites.CreateSprite(context.Background(), s); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.state.NotifyDataChanged()
		v.reload()

		// Seleciona o novo sprite
		for i, sp := range v.sprites {
			if sp.ID == s.ID {
				v.spriteList.Select(i)
				break
			}
		}
	}, v.win)
}

func (v *SpriteView) deleteSelectedSprite() {
	if v.selectedSprite == nil {
		dialog.ShowInformation("Aviso", "Selecione um sprite para excluir.", v.win)
		return
	}
	proj := v.state.Current()
	if proj == nil {
		return
	}

	msg := fmt.Sprintf("Tem certeza que deseja excluir o sprite '%s' (#%d)?", v.selectedSprite.Name, v.selectedSprite.ID)
	dialog.ShowConfirm("Excluir Sprite", msg, func(confirmed bool) {
		if !confirmed {
			return
		}
		if err := proj.Storage().Sprites.DeleteSprite(context.Background(), v.selectedSprite.ID); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.state.NotifyDataChanged()
		v.selectedSprite = nil
		v.reload()
	}, v.win)
}
