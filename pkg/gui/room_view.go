package gui

import (
	"context"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

// RoomView implementa o editor de salas e cenários do Z-Realm (SCREEN 4 - 32x18).
type RoomView struct {
	state      *ProjectState
	win        fyne.Window
	roomList   *widget.List
	rooms      []*models.Room
	selected   *models.Room
	editor     *RoomEditorWidget
}

// NewRoomView cria o componente de visualização e edição de salas (32x18 tiles).
func NewRoomView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &RoomView{
		state: state,
		win:   win,
		rooms: make([]*models.Room, 0),
	}

	// 1. Instancia o editor completo de salas (Canvas 32x18, Paleta, Conexões, Entidades)
	v.editor = NewRoomEditorWidget(state, win, func() {
		v.roomList.Refresh()
	}, func(roomID int64) {
		v.selectRoomByID(roomID)
	})

	// 2. Lista de Salas do Projeto
	v.roomList = widget.NewList(
		func() int {
			return len(v.rooms)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Sala"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.rooms) {
				return
			}
			box := obj.(*fyne.Container)
			lbl := box.Objects[1].(*widget.Label)
			r := v.rooms[id]
			lbl.SetText(fmt.Sprintf("#%d: %s (%d, %d)", r.ID, r.Name, r.WorldX, r.WorldY))
		},
	)

	v.roomList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.rooms) {
			return
		}
		v.selected = v.rooms[id]
		v.editor.LoadRoom(v.selected, v.rooms)
	}

	btnAdd := widget.NewButtonWithIcon("Nova Sala", theme.ContentAddIcon(), func() {
		v.createNewRoom()
	})

	btnDelete := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedRoom()
	})

	btnRefresh := widget.NewButtonWithIcon("Recarregar", theme.ViewRefreshIcon(), func() {
		v.reload()
	})

	listHeader := container.NewVBox(
		widget.NewLabelWithStyle("🗺️ Salas do Mundo", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAdd, btnDelete, btnRefresh),
		widget.NewSeparator(),
	)

	sidebar := container.NewBorder(listHeader, nil, nil, nil, v.roomList)

	state.OnProjectLoaded(func(_ *project.Project, _ string) {
		v.reload()
	})
	state.OnProjectClosed(func() {
		v.rooms = make([]*models.Room, 0)
		v.selected = nil
		v.roomList.Refresh()
		v.editor.LoadRoom(nil, nil)
	})

	if state.IsOpen() {
		v.reload()
	}

	split := container.NewHSplit(sidebar, v.editor.Container())
	split.Offset = 0.22
	return split
}

func (v *RoomView) reload() {
	proj := v.state.Current()
	if proj == nil {
		v.rooms = make([]*models.Room, 0)
		v.roomList.Refresh()
		v.editor.LoadRoom(nil, nil)
		return
	}

	rooms, err := proj.Storage().Rooms.ListRooms(context.Background())
	if err != nil {
		dialog.ShowError(fmt.Errorf("erro ao carregar salas: %w", err), v.win)
		return
	}

	v.rooms = rooms
	v.roomList.Refresh()

	if len(v.rooms) > 0 {
		targetIdx := 0
		if v.selected != nil {
			for i, r := range v.rooms {
				if r.ID == v.selected.ID {
					targetIdx = i
					break
				}
			}
		}
		v.roomList.Select(targetIdx)
	} else {
		v.selected = nil
		v.editor.LoadRoom(nil, nil)
	}
}

func (v *RoomView) selectRoomByID(roomID int64) {
	for i, r := range v.rooms {
		if r.ID == roomID {
			v.roomList.Select(i)
			return
		}
	}
}

func (v *RoomView) createNewRoom() {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("Sala %d", len(v.rooms)+1))

	nextX := 0
	nextY := 0
	if len(v.rooms) > 0 {
		last := v.rooms[len(v.rooms)-1]
		nextX = last.WorldX + 1
		nextY = last.WorldY
	}

	worldXEntry := widget.NewEntry()
	worldXEntry.SetText(fmt.Sprintf("%d", nextX))

	worldYEntry := widget.NewEntry()
	worldYEntry.SetText(fmt.Sprintf("%d", nextY))

	// Busca tilesets disponíveis
	tilesets, _ := proj.Storage().Tilesets.ListTilesets(context.Background())
	tilesetOptions := []string{}
	tilesetMap := make(map[string]int64)
	for _, ts := range tilesets {
		lbl := fmt.Sprintf("#%d: %s", ts.ID, ts.Name)
		tilesetOptions = append(tilesetOptions, lbl)
		tilesetMap[lbl] = ts.ID
	}
	if len(tilesetOptions) == 0 {
		tilesetOptions = append(tilesetOptions, "#1: Tileset Padrão")
		tilesetMap["#1: Tileset Padrão"] = 1
	}

	tilesetSelect := widget.NewSelect(tilesetOptions, nil)
	tilesetSelect.SetSelectedIndex(0)

	items := []*widget.FormItem{
		widget.NewFormItem("Nome da Sala", nameEntry),
		widget.NewFormItem("Coordenada World X", worldXEntry),
		widget.NewFormItem("Coordenada World Y", worldYEntry),
		widget.NewFormItem("Tileset Vinculado", tilesetSelect),
	}

	dialog.ShowForm("Nova Sala (SCREEN 4)", "Criar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}

		wx, errX := strconv.Atoi(worldXEntry.Text)
		wy, errY := strconv.Atoi(worldYEntry.Text)
		if errX != nil || errY != nil {
			dialog.ShowError(fmt.Errorf("coordenadas World X e World Y devem ser números inteiros"), v.win)
			return
		}

		tsID := tilesetMap[tilesetSelect.Selected]
		if tsID == 0 {
			tsID = 1
		}

		emptyMatrix := make([]byte, models.RoomMatrixSize)
		r := &models.Room{
			WorldX:     wx,
			WorldY:     wy,
			Name:       nameEntry.Text,
			TilesetID:  tsID,
			TileMatrix: emptyMatrix,
		}

		if err := proj.Storage().Rooms.CreateRoom(context.Background(), r); err != nil {
			dialog.ShowError(err, v.win)
			return
		}

		v.state.NotifyDataChanged()
		v.reload()
		v.selectRoomByID(r.ID)
	}, v.win)
}

func (v *RoomView) deleteSelectedRoom() {
	if v.selected == nil {
		dialog.ShowInformation("Aviso", "Selecione uma sala para excluir.", v.win)
		return
	}

	proj := v.state.Current()
	if proj == nil {
		return
	}

	roomName := v.selected.Name
	roomID := v.selected.ID

	dialog.ShowConfirm(
		"Excluir Sala",
		fmt.Sprintf("Deseja realmente excluir a sala '%s' (#%d)?\nTodas as entidades associadas a ela serão excluídas.", roomName, roomID),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			if err := proj.Storage().Rooms.DeleteRoom(context.Background(), roomID); err != nil {
				dialog.ShowError(fmt.Errorf("falha ao excluir sala: %w", err), v.win)
				return
			}

			v.selected = nil
			v.state.NotifyDataChanged()
			v.reload()
		},
		v.win,
	)
}
