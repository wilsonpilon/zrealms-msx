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

// TileView implementa a visualização e edição interativa de tilesets e padrões 8x8 do MSX 2 (SCREEN 4).
type TileView struct {
	state           *ProjectState
	win             fyne.Window
	tilesetList     *widget.List
	tilesets        []*models.Tileset
	selectedTileset *models.Tileset
	tiles           []*models.Tile
	tileList        *widget.List
	selectedTile    *models.Tile
	editor          *TileEditorWidget
}

// NewTileView cria o componente de visualização e edição de tilesets.
func NewTileView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &TileView{
		state:    state,
		win:      win,
		tilesets: make([]*models.Tileset, 0),
		tiles:    make([]*models.Tile, 0),
	}

	// 1. Lista de Tilesets
	v.tilesetList = widget.NewList(
		func() int {
			return len(v.tilesets)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.GridIcon()),
				widget.NewLabel("Tileset"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.tilesets) {
				return
			}
			box := obj.(*fyne.Container)
			lbl := box.Objects[1].(*widget.Label)
			ts := v.tilesets[id]
			lbl.SetText(fmt.Sprintf("#%d: %s", ts.ID, ts.Name))
		},
	)

	v.tilesetList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.tilesets) {
			return
		}
		v.selectedTileset = v.tilesets[id]
		v.loadTilesForTileset(v.selectedTileset.ID)
	}

	// 2. Lista de Tiles do Tileset Selecionado
	v.tileList = widget.NewList(
		func() int {
			return len(v.tiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ContentPasteIcon()),
				widget.NewLabel("Tile"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.tiles) {
				return
			}
			box := obj.(*fyne.Container)
			lbl := box.Objects[1].(*widget.Label)
			t := v.tiles[id]
			lbl.SetText(fmt.Sprintf("Tile %03d: %s", t.TileIndex, collisionName(t.CollisionType)))
		},
	)

	v.tileList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.tiles) {
			return
		}
		v.selectedTile = v.tiles[id]
		v.editor.LoadTile(v.selectedTile)
	}

	// 3. Editor Interativo de Tiles (Grid 8x8 + Paleta)
	v.editor = NewTileEditorWidget(state, win, func() {
		v.tileList.Refresh()
	})

	// Botões de Ação
	btnAddTileset := widget.NewButtonWithIcon("Novo Tileset", theme.ContentAddIcon(), func() {
		v.createTileset()
	})

	btnAddTile := widget.NewButtonWithIcon("Adicionar Tile", theme.ContentAddIcon(), func() {
		v.createTile()
	})

	tilesetHeader := container.NewVBox(
		widget.NewLabelWithStyle("🧱 Tilesets", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		btnAddTileset,
		widget.NewSeparator(),
	)
	tilesetPanel := container.NewBorder(tilesetHeader, nil, nil, nil, v.tilesetList)

	tileHeader := container.NewVBox(
		widget.NewLabelWithStyle("Padrões 8x8 no Tileset", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		btnAddTile,
		widget.NewSeparator(),
	)
	tilePanel := container.NewBorder(tileHeader, nil, nil, nil, v.tileList)

	leftSplit := container.NewHSplit(tilesetPanel, tilePanel)
	leftSplit.Offset = 0.40

	mainSplit := container.NewHSplit(leftSplit, v.editor.Container())
	mainSplit.Offset = 0.35

	state.OnProjectLoaded(func(_ *project.Project, _ string) {
		v.reload()
	})
	state.OnProjectClosed(func() {
		v.tilesets = make([]*models.Tileset, 0)
		v.tiles = make([]*models.Tile, 0)
		v.selectedTileset = nil
		v.selectedTile = nil
		v.tilesetList.Refresh()
		v.tileList.Refresh()
		v.editor.LoadTile(nil)
	})

	if state.IsOpen() {
		v.reload()
	}

	return mainSplit
}

func (v *TileView) reload() {
	proj := v.state.Current()
	if proj == nil {
		return
	}
	ts, err := proj.Storage().Tilesets.ListTilesets(context.Background())
	if err != nil {
		dialog.ShowError(err, v.win)
		return
	}
	v.tilesets = ts
	v.tilesetList.Refresh()
	if len(v.tilesets) > 0 {
		v.tilesetList.Select(0)
	}
}

func (v *TileView) loadTilesForTileset(tilesetID int64) {
	proj := v.state.Current()
	if proj == nil {
		return
	}
	tiles, err := proj.Storage().Tilesets.ListTiles(context.Background(), tilesetID)
	if err != nil {
		dialog.ShowError(err, v.win)
		return
	}
	v.tiles = tiles
	v.tileList.Refresh()
	if len(v.tiles) > 0 {
		v.tileList.Select(0)
	} else {
		v.selectedTile = nil
		v.editor.LoadTile(nil)
	}
}

func (v *TileView) createTileset() {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("Tileset %d", len(v.tilesets)+1))
	descEntry := widget.NewEntry()
	descEntry.SetText("Conjunto de padrões para masmorra")

	items := []*widget.FormItem{
		widget.NewFormItem("Nome do Tileset", nameEntry),
		widget.NewFormItem("Descrição", descEntry),
	}

	dialog.ShowForm("Novo Tileset", "Criar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		ts := &models.Tileset{
			Name:        nameEntry.Text,
			Description: descEntry.Text,
		}
		if err := proj.Storage().Tilesets.CreateTileset(context.Background(), ts); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.state.NotifyDataChanged()
		v.reload()
	}, v.win)
}

func (v *TileView) createTile() {
	proj := v.state.Current()
	if proj == nil || v.selectedTileset == nil {
		dialog.ShowInformation("Aviso", "Selecione um tileset primeiro para adicionar tiles.", v.win)
		return
	}

	// Calcula próximo índice livre de 0 a 255
	usedIndices := make(map[int]bool)
	for _, t := range v.tiles {
		usedIndices[t.TileIndex] = true
	}

	nextIdx := 0
	for nextIdx < 256 {
		if !usedIndices[nextIdx] {
			break
		}
		nextIdx++
	}

	if nextIdx >= 256 {
		dialog.ShowError(fmt.Errorf("limite do VDP atingido: este tileset já possui todos os 256 tiles alocados"), v.win)
		return
	}

	// Cria o novo tile com padrão monocromático e cores padrão (Fg=15 Branco, Bg=1 Preto)
	newTile := &models.Tile{
		TilesetID:     v.selectedTileset.ID,
		TileIndex:     nextIdx,
		PatternBytes:  make([]byte, 8),
		ColorBytes:    []byte{0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1},
		CollisionType: models.CollisionPassable,
	}

	err := proj.Storage().Tilesets.SaveTile(context.Background(), newTile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("falha ao criar tile: %w", err), v.win)
		return
	}

	v.state.NotifyDataChanged()
	v.loadTilesForTileset(v.selectedTileset.ID)

	// Seleciona o novo tile criado
	for i, t := range v.tiles {
		if t.TileIndex == nextIdx {
			v.tileList.Select(i)
			break
		}
	}
}

func collisionName(c models.CollisionType) string {
	switch c {
	case models.CollisionPassable:
		return "Passável"
	case models.CollisionSolid:
		return "Sólido"
	case models.CollisionWater:
		return "Água"
	case models.CollisionDamage:
		return "Dano"
	case models.CollisionTrigger:
		return "Gatilho"
	default:
		return fmt.Sprintf("Custom (%d)", c)
	}
}
