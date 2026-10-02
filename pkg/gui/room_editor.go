package gui

import (
	"context"
	"fmt"
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// RoomEditorWidget reúne o viewport 32x18 da sala, paleta de tiles, conexões cardeais e gestor de entidades.
type RoomEditorWidget struct {
	state            *ProjectState
	win              fyne.Window
	room             *models.Room
	allRooms         []*models.Room
	tiles            []*models.Tile
	entities         []*models.Entity
	sprites          []*models.Sprite
	canvas           *RoomCanvas
	container        *fyne.Container

	// Controles do Cabeçalho da Sala
	entryName        *widget.Entry
	entryWorldX      *widget.Entry
	entryWorldY      *widget.Entry
	lblStatus        *widget.Label
	lblHover         *widget.Label

	// Ferramentas e Tile Ativo
	toolRadio        *widget.RadioGroup
	lblSelectedTile  *widget.Label
	imgSelectedTile  *canvas.Image
	paletteBox       *fyne.Container

	// Conexões Cardeais
	lblNorth         *widget.Label
	btnGoNorth       *widget.Button
	btnSetNorth      *widget.Button

	lblSouth         *widget.Label
	btnGoSouth       *widget.Button
	btnSetSouth      *widget.Button

	lblEast          *widget.Label
	btnGoEast        *widget.Button
	btnSetEast       *widget.Button

	lblWest          *widget.Label
	btnGoWest        *widget.Button
	btnSetWest       *widget.Button

	btnAutoConnect   *widget.Button

	// Gestor de Entidades
	entityList       *widget.List
	selectedEntity   *models.Entity
	btnAddEntity     *widget.Button
	btnEditEntity    *widget.Button
	btnDeleteEntity  *widget.Button

	// Callbacks
	onRoomSaved      func()
	onNavigateToRoom func(roomID int64)
}

// NewRoomEditorWidget instancia o editor interativo de salas (32x18 tiles SCREEN 4).
func NewRoomEditorWidget(state *ProjectState, win fyne.Window, onSaved func(), onNavigate func(roomID int64)) *RoomEditorWidget {
	re := &RoomEditorWidget{
		state:            state,
		win:              win,
		allRooms:         make([]*models.Room, 0),
		tiles:            make([]*models.Tile, 0),
		entities:         make([]*models.Entity, 0),
		sprites:          make([]*models.Sprite, 0),
		entryName:        widget.NewEntry(),
		entryWorldX:      widget.NewEntry(),
		entryWorldY:      widget.NewEntry(),
		lblStatus:        widget.NewLabel("Selecione uma sala para iniciar a edição."),
		lblHover:         widget.NewLabel("Cursor: X=--, Y=-- | Tile: #---"),
		lblSelectedTile:  widget.NewLabel("Tile Ativo: #000"),
		onRoomSaved:      onSaved,
		onNavigateToRoom: onNavigate,
	}

	re.entryName.SetPlaceHolder("Nome da Sala")
	re.entryWorldX.SetPlaceHolder("0")
	re.entryWorldY.SetPlaceHolder("0")

	// 1. Canvas 32x18 (256x144 pixels)
	re.canvas = NewRoomCanvas(nil, nil, nil)
	re.canvas.OnChanged(func() {
		re.lblStatus.SetText("Modificações não salvas na matriz de tiles.")
	})
	re.canvas.OnTilePicked(func(tileIdx byte) {
		re.selectTile(tileIdx)
		re.lblStatus.SetText(fmt.Sprintf("Tile #%03d selecionado com o Conta-Gotas.", tileIdx))
	})
	re.canvas.OnHoverCoords(func(x, y int, tileIdx byte) {
		re.lblHover.SetText(fmt.Sprintf("Cursor: X=%02d, Y=%02d | Tile: #%03d", x, y, tileIdx))
	})

	// 2. Pré-visualização do Tile Selecionado
	re.imgSelectedTile = canvas.NewImageFromImage(RenderTileThumbnail(nil))
	re.imgSelectedTile.SetMinSize(fyne.NewSize(32, 32))
	re.imgSelectedTile.ScaleMode = canvas.ImageScalePixels

	// 3. Ferramentas de Pintura
	re.toolRadio = widget.NewRadioGroup([]string{
		"Pincel (Carimbo)",
		"Balde (Preencher)",
		"Borracha (Tile 0)",
		"Conta-Gotas",
	}, func(selected string) {
		switch selected {
		case "Pincel (Carimbo)":
			re.canvas.SetTool(ToolBrush)
		case "Balde (Preencher)":
			re.canvas.SetTool(ToolFill)
		case "Borracha (Tile 0)":
			re.canvas.SetTool(ToolEraser)
		case "Conta-Gotas":
			re.canvas.SetTool(ToolEyedropper)
		}
	})
	re.toolRadio.Horizontal = true
	re.toolRadio.SetSelected("Pincel (Carimbo)")

	// Botões utilitários de preenchimento rápido
	btnClearAll := widget.NewButton("Limpar (Tile 0)", func() {
		dialog.ShowConfirm("Limpar Sala", "Deseja preencher toda a sala com o tile 0?", func(ok bool) {
			if ok {
				re.canvas.Clear(0)
				re.lblStatus.SetText("Sala limpa com o tile 0.")
			}
		}, re.win)
	})

	btnFillAll := widget.NewButton("Preencher Tudo", func() {
		t := re.canvas.GetSelectedTile()
		dialog.ShowConfirm("Preencher Tudo", fmt.Sprintf("Preencher toda a sala com o Tile #%03d?", t), func(ok bool) {
			if ok {
				re.canvas.Clear(t)
				re.lblStatus.SetText(fmt.Sprintf("Sala totalmente preenchida com Tile #%03d.", t))
			}
		}, re.win)
	})

	btnFillBorder := widget.NewButton("Preencher Bordas", func() {
		t := re.canvas.GetSelectedTile()
		re.canvas.FillBorder(t)
		re.lblStatus.SetText(fmt.Sprintf("Bordas da sala preenchidas com Tile #%03d.", t))
	})

	// 4. Paleta de Carimbo (Tiles do Tileset)
	re.paletteBox = container.NewHBox()
	paletteScroll := container.NewHScroll(re.paletteBox)
	paletteScroll.SetMinSize(fyne.NewSize(512, 54))

	paletteCard := widget.NewCard(
		"Paleta de Carimbo",
		"Clique em um tile para selecioná-lo para desenho",
		paletteScroll,
	)

	// 5. Conexões Cardeais (Norte, Sul, Leste, Oeste)
	re.lblNorth = widget.NewLabel("Norte: (Nenhum)")
	re.btnGoNorth = widget.NewButtonWithIcon("Ir", theme.NavigateNextIcon(), func() {
		if re.room != nil && re.room.NorthRoomID != nil && re.onNavigateToRoom != nil {
			re.onNavigateToRoom(*re.room.NorthRoomID)
		}
	})
	re.btnSetNorth = widget.NewButton("Alterar", func() { re.pickConnectionDialog("Norte", &re.room.NorthRoomID) })

	re.lblSouth = widget.NewLabel("Sul: (Nenhum)")
	re.btnGoSouth = widget.NewButtonWithIcon("Ir", theme.NavigateNextIcon(), func() {
		if re.room != nil && re.room.SouthRoomID != nil && re.onNavigateToRoom != nil {
			re.onNavigateToRoom(*re.room.SouthRoomID)
		}
	})
	re.btnSetSouth = widget.NewButton("Alterar", func() { re.pickConnectionDialog("Sul", &re.room.SouthRoomID) })

	re.lblEast = widget.NewLabel("Leste: (Nenhum)")
	re.btnGoEast = widget.NewButtonWithIcon("Ir", theme.NavigateNextIcon(), func() {
		if re.room != nil && re.room.EastRoomID != nil && re.onNavigateToRoom != nil {
			re.onNavigateToRoom(*re.room.EastRoomID)
		}
	})
	re.btnSetEast = widget.NewButton("Alterar", func() { re.pickConnectionDialog("Leste", &re.room.EastRoomID) })

	re.lblWest = widget.NewLabel("Oeste: (Nenhum)")
	re.btnGoWest = widget.NewButtonWithIcon("Ir", theme.NavigateNextIcon(), func() {
		if re.room != nil && re.room.WestRoomID != nil && re.onNavigateToRoom != nil {
			re.onNavigateToRoom(*re.room.WestRoomID)
		}
	})
	re.btnSetWest = widget.NewButton("Alterar", func() { re.pickConnectionDialog("Oeste", &re.room.WestRoomID) })

	re.btnAutoConnect = widget.NewButtonWithIcon("Auto-Conectar por Coordenadas", theme.MediaFastForwardIcon(), func() {
		re.autoConnectAllRooms()
	})

	cardinalGrid := container.NewVBox(
		container.NewHBox(re.lblNorth, re.btnGoNorth, re.btnSetNorth),
		container.NewHBox(re.lblSouth, re.btnGoSouth, re.btnSetSouth),
		container.NewHBox(re.lblEast, re.btnGoEast, re.btnSetEast),
		container.NewHBox(re.lblWest, re.btnGoWest, re.btnSetWest),
		widget.NewSeparator(),
		re.btnAutoConnect,
	)

	cardinalCard := widget.NewCard(
		"🧭 Conexões Cardeais",
		"Transição contígua de telas no V9938 (SCREEN 4)",
		cardinalGrid,
	)

	// 6. Gestor de Entidades da Sala
	re.entityList = widget.NewList(
		func() int {
			return len(re.entities)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.AccountIcon()),
				widget.NewLabel("Entidade"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(re.entities) {
				return
			}
			box := obj.(*fyne.Container)
			lbl := box.Objects[1].(*widget.Label)
			ent := re.entities[id]
			lbl.SetText(fmt.Sprintf("#%d: %s (%02d,%02d) - %s", ent.ID, ent.Name, ent.PosX, ent.PosY, behaviorName(ent.BehaviorType)))
		},
	)
	re.entityList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(re.entities) {
			re.selectedEntity = re.entities[id]
		}
	}

	re.btnAddEntity = widget.NewButtonWithIcon("Adicionar", theme.ContentAddIcon(), func() {
		re.showEntityForm(nil)
	})
	re.btnEditEntity = widget.NewButtonWithIcon("Editar", theme.DocumentCreateIcon(), func() {
		if re.selectedEntity != nil {
			re.showEntityForm(re.selectedEntity)
		} else {
			dialog.ShowInformation("Aviso", "Selecione uma entidade na lista para editar.", re.win)
		}
	})
	re.btnDeleteEntity = widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		re.deleteSelectedEntity()
	})

	entityActions := container.NewHBox(re.btnAddEntity, re.btnEditEntity, re.btnDeleteEntity)
	// Envolve a lista com tamanho adequado
	entityContainer := container.NewBorder(nil, entityActions, nil, nil, re.entityList)

	entityCard := widget.NewCard(
		"👾 Entidades e Atores",
		"NPCs, baús, portas, gatilhos e inimigos na sala",
		container.NewGridWrap(fyne.NewSize(320, 200), entityContainer),
	)

	// 7. Botão Principal de Salvamento
	btnSave := widget.NewButtonWithIcon("Salvar Sala no Banco SQLite", theme.DocumentSaveIcon(), func() {
		re.SaveCurrentRoom()
	})
	btnSave.Importance = widget.HighImportance

	// 8. Cabeçalho de Propriedades da Sala
	headerForm := container.NewHBox(
		widget.NewLabelWithStyle("Nome:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWrap(fyne.NewSize(200, 36), re.entryName),
		widget.NewLabelWithStyle("World X:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWrap(fyne.NewSize(60, 36), re.entryWorldX),
		widget.NewLabelWithStyle("World Y:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWrap(fyne.NewSize(60, 36), re.entryWorldY),
		btnSave,
	)

	// Coluna da Esquerda (Viewport 32x18 + Paleta + Ferramentas)
	toolsBar := container.NewHBox(
		re.toolRadio,
		widget.NewSeparator(),
		re.imgSelectedTile,
		re.lblSelectedTile,
		widget.NewSeparator(),
		btnClearAll,
		btnFillAll,
		btnFillBorder,
	)

	leftCol := container.NewVBox(
		headerForm,
		widget.NewSeparator(),
		toolsBar,
		re.lblHover,
		container.NewCenter(re.canvas),
		paletteCard,
	)

	// Coluna da Direita (Conexões Cardeais + Entidades)
	rightCol := container.NewVBox(
		cardinalCard,
		widget.NewSeparator(),
		entityCard,
		widget.NewSeparator(),
		re.lblStatus,
	)

	split := container.NewHSplit(leftCol, rightCol)
	split.Offset = 0.65

	re.container = container.NewPadded(split)
	return re
}

// Container retorna a visualização principal do editor de salas.
func (re *RoomEditorWidget) Container() fyne.CanvasObject {
	return re.container
}

// LoadRoom carrega uma sala e seus dados associados para o editor.
func (re *RoomEditorWidget) LoadRoom(room *models.Room, allRooms []*models.Room) {
	re.room = room
	re.allRooms = allRooms
	if room == nil {
		re.canvas.SetRoom(nil, nil)
		re.entryName.SetText("")
		re.entryWorldX.SetText("")
		re.entryWorldY.SetText("")
		re.lblStatus.SetText("Nenhuma sala selecionada.")
		re.updateConnectionsUI()
		return
	}

	re.entryName.SetText(room.Name)
	re.entryWorldX.SetText(fmt.Sprintf("%d", room.WorldX))
	re.entryWorldY.SetText(fmt.Sprintf("%d", room.WorldY))

	// Carrega tileset e entidades do banco de dados
	proj := re.state.Current()
	if proj != nil {
		tiles, err := proj.Storage().Tilesets.ListTiles(context.Background(), room.TilesetID)
		if err == nil {
			re.tiles = tiles
			re.canvas.SetTiles(tiles)
		}

		entities, err := proj.Storage().Rooms.ListEntitiesByRoom(context.Background(), room.ID)
		if err == nil {
			re.entities = entities
		} else {
			re.entities = make([]*models.Entity, 0)
		}

		sprites, err := proj.Storage().Sprites.ListSprites(context.Background())
		if err == nil {
			re.sprites = sprites
		}
	}

	re.canvas.SetRoom(room, re.entities)
	re.selectedEntity = nil
	re.entityList.Refresh()
	re.rebuildPalette()
	re.updateConnectionsUI()
	re.lblStatus.SetText(fmt.Sprintf("Editando Sala #%d: %s (%d, %d)", room.ID, room.Name, room.WorldX, room.WorldY))
}

func (re *RoomEditorWidget) selectTile(idx byte) {
	re.canvas.SetSelectedTile(idx)
	re.lblSelectedTile.SetText(fmt.Sprintf("Tile Ativo: #%03d", idx))

	var targetTile *models.Tile
	for _, t := range re.tiles {
		if t.TileIndex == int(idx) {
			targetTile = t
			break
		}
	}
	re.imgSelectedTile.Image = RenderTileThumbnail(targetTile)
	re.imgSelectedTile.Refresh()
}

func (re *RoomEditorWidget) rebuildPalette() {
	re.paletteBox.Objects = nil

	// Se houver tiles definidos no tileset, exibe cada um
	if len(re.tiles) > 0 {
		for _, t := range re.tiles {
			tileIdx := byte(t.TileIndex)
			btn := widget.NewButton(fmt.Sprintf("#%03d", tileIdx), func() {
				re.selectTile(tileIdx)
			})
			re.paletteBox.Add(btn)
		}
	} else {
		// Tiles padrão 0 a 15
		for i := 0; i < 16; i++ {
			tileIdx := byte(i)
			btn := widget.NewButton(fmt.Sprintf("#%03d", tileIdx), func() {
				re.selectTile(tileIdx)
			})
			re.paletteBox.Add(btn)
		}
	}

	// Campo direto de seleção numérica
	tileInput := widget.NewEntry()
	tileInput.SetPlaceHolder("0..255")
	btnPickCustom := widget.NewButton("Ir", func() {
		val, err := strconv.Atoi(tileInput.Text)
		if err == nil && val >= 0 && val <= 255 {
			re.selectTile(byte(val))
		}
	})
	re.paletteBox.Add(widget.NewLabel("Outro:"))
	re.paletteBox.Add(container.NewGridWrap(fyne.NewSize(70, 36), tileInput))
	re.paletteBox.Add(btnPickCustom)

	re.paletteBox.Refresh()
}

func (re *RoomEditorWidget) updateConnectionsUI() {
	if re.room == nil {
		re.lblNorth.SetText("Norte: (Nenhum)")
		re.btnGoNorth.Disable()
		re.lblSouth.SetText("Sul: (Nenhum)")
		re.btnGoSouth.Disable()
		re.lblEast.SetText("Leste: (Nenhum)")
		re.btnGoEast.Disable()
		re.lblWest.SetText("Oeste: (Nenhum)")
		re.btnGoWest.Disable()
		return
	}

	re.updateConnectionLabel(re.room.NorthRoomID, re.lblNorth, re.btnGoNorth, "Norte")
	re.updateConnectionLabel(re.room.SouthRoomID, re.lblSouth, re.btnGoSouth, "Sul")
	re.updateConnectionLabel(re.room.EastRoomID, re.lblEast, re.btnGoEast, "Leste")
	re.updateConnectionLabel(re.room.WestRoomID, re.lblWest, re.btnGoWest, "Oeste")
}

func (re *RoomEditorWidget) updateConnectionLabel(roomID *int64, lbl *widget.Label, btn *widget.Button, dir string) {
	if roomID == nil {
		lbl.SetText(fmt.Sprintf("%s: (Nenhum)", dir))
		btn.Disable()
		return
	}

	var target *models.Room
	for _, r := range re.allRooms {
		if r.ID == *roomID {
			target = r
			break
		}
	}

	if target != nil {
		lbl.SetText(fmt.Sprintf("%s: #%d (%s)", dir, target.ID, target.Name))
		btn.Enable()
	} else {
		lbl.SetText(fmt.Sprintf("%s: #%d (ID Inválido)", dir, *roomID))
		btn.Disable()
	}
}

func (re *RoomEditorWidget) pickConnectionDialog(dir string, targetID **int64) {
	if re.room == nil {
		return
	}

	options := []string{"(Nenhum - Desconectar)"}
	idMap := make(map[string]*int64)
	idMap["(Nenhum - Desconectar)"] = nil

	for _, r := range re.allRooms {
		if r.ID == re.room.ID {
			continue // Não conectar consigo mesma
		}
		label := fmt.Sprintf("#%d: %s (%d, %d)", r.ID, r.Name, r.WorldX, r.WorldY)
		options = append(options, label)
		rid := r.ID
		idMap[label] = &rid
	}

	sel := widget.NewSelect(options, nil)
	if *targetID == nil {
		sel.SetSelected("(Nenhum - Desconectar)")
	} else {
		for opt, id := range idMap {
			if id != nil && *id == **targetID {
				sel.SetSelected(opt)
				break
			}
		}
	}

	dialog.ShowCustomConfirm(
		fmt.Sprintf("Vincular Conexão %s", dir),
		"Confirmar",
		"Cancelar",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Selecione a sala conectada ao lado %s:", dir)),
			sel,
		),
		func(confirmed bool) {
			if confirmed {
				chosen := idMap[sel.Selected]
				*targetID = chosen
				re.updateConnectionsUI()
				re.lblStatus.SetText(fmt.Sprintf("Conexão %s atualizada.", dir))
			}
		},
		re.win,
	)
}

func (re *RoomEditorWidget) autoConnectAllRooms() {
	if len(re.allRooms) == 0 {
		dialog.ShowInformation("Aviso", "Não há salas cadastradas para conectar.", re.win)
		return
	}

	changes := AutoConnectRooms(re.allRooms)
	proj := re.state.Current()
	if proj != nil {
		for _, r := range re.allRooms {
			_ = proj.Storage().Rooms.UpdateRoom(context.Background(), r)
		}
	}

	re.updateConnectionsUI()
	re.lblStatus.SetText(fmt.Sprintf("Auto-conexão concluída: %d alterações realizadas.", changes))
	dialog.ShowInformation("Auto-Conexão Concluída", fmt.Sprintf("Foram atualizadas %d conexões cardeais entre as salas baseadas em suas posições (WorldX, WorldY).", changes), re.win)
}

func (re *RoomEditorWidget) showEntityForm(ent *models.Entity) {
	if re.room == nil {
		return
	}

	isNew := ent == nil
	if isNew {
		ent = &models.Entity{
			RoomID:       re.room.ID,
			Name:         "Novo NPC",
			PosX:         models.RoomWidth / 2,
			PosY:         models.RoomHeight / 2,
			BehaviorType: models.BehaviorStaticNPC,
		}
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(ent.Name)

	posXEntry := widget.NewEntry()
	posXEntry.SetText(fmt.Sprintf("%d", ent.PosX))

	posYEntry := widget.NewEntry()
	posYEntry.SetText(fmt.Sprintf("%d", ent.PosY))

	behaviorOptions := []string{
		"0 - NPC Estacionário",
		"1 - NPC Andarilho",
		"2 - Patrulha Predefinida",
		"3 - Baú de Itens",
		"4 - Porta / Passagem",
		"5 - Gatilho de Evento",
		"6 - Inimigo Hostil",
	}
	behaviorSelect := widget.NewSelect(behaviorOptions, nil)
	behaviorSelect.SetSelectedIndex(int(ent.BehaviorType))

	// Opções de Sprite
	spriteOptions := []string{"(Nenhum)"}
	spriteMap := make(map[string]*int64)
	spriteMap["(Nenhum)"] = nil
	for _, s := range re.sprites {
		lbl := fmt.Sprintf("#%d: %s", s.ID, s.Name)
		spriteOptions = append(spriteOptions, lbl)
		sid := s.ID
		spriteMap[lbl] = &sid
	}
	spriteSelect := widget.NewSelect(spriteOptions, nil)
	if ent.SpriteID != nil {
		for lbl, sid := range spriteMap {
			if sid != nil && *sid == *ent.SpriteID {
				spriteSelect.SetSelected(lbl)
				break
			}
		}
	} else {
		spriteSelect.SetSelected("(Nenhum)")
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Nome da Entidade", nameEntry),
		widget.NewFormItem("Posição X (0..31)", posXEntry),
		widget.NewFormItem("Posição Y (0..17)", posYEntry),
		widget.NewFormItem("Tipo de Comportamento", behaviorSelect),
		widget.NewFormItem("Sprite (Modo 2)", spriteSelect),
	}

	title := "Adicionar Nova Entidade"
	if !isNew {
		title = fmt.Sprintf("Editar Entidade #%d", ent.ID)
	}

	dialog.ShowForm(title, "Salvar", "Cancelar", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}

		px, errX := strconv.Atoi(posXEntry.Text)
		py, errY := strconv.Atoi(posYEntry.Text)
		if errX != nil || px < 0 || px >= models.RoomWidth || errY != nil || py < 0 || py >= models.RoomHeight {
			dialog.ShowError(fmt.Errorf("coordenadas inválidas: X deve ser 0..31 e Y deve ser 0..17"), re.win)
			return
		}

		ent.Name = nameEntry.Text
		ent.PosX = px
		ent.PosY = py
		ent.BehaviorType = models.BehaviorType(behaviorSelect.SelectedIndex())
		ent.SpriteID = spriteMap[spriteSelect.Selected]

		proj := re.state.Current()
		if proj == nil {
			return
		}

		if isNew {
			if err := proj.Storage().Rooms.CreateEntity(context.Background(), ent); err != nil {
				dialog.ShowError(err, re.win)
				return
			}
			re.entities = append(re.entities, ent)
			re.lblStatus.SetText(fmt.Sprintf("Entidade '%s' adicionada na posição (%d, %d).", ent.Name, px, py))
		} else {
			if err := proj.Storage().Rooms.UpdateEntity(context.Background(), ent); err != nil {
				dialog.ShowError(err, re.win)
				return
			}
			re.lblStatus.SetText(fmt.Sprintf("Entidade '%s' atualizada.", ent.Name))
		}

		re.canvas.SetEntities(re.entities)
		re.entityList.Refresh()
	}, re.win)
}

func (re *RoomEditorWidget) deleteSelectedEntity() {
	if re.selectedEntity == nil {
		dialog.ShowInformation("Aviso", "Selecione uma entidade para excluir.", re.win)
		return
	}

	ent := re.selectedEntity
	dialog.ShowConfirm("Excluir Entidade", fmt.Sprintf("Deseja realmente excluir a entidade '%s'?", ent.Name), func(ok bool) {
		if !ok {
			return
		}

		proj := re.state.Current()
		if proj != nil {
			if err := proj.Storage().Rooms.DeleteEntity(context.Background(), ent.ID); err != nil {
				dialog.ShowError(err, re.win)
				return
			}
		}

		// Remove da lista em memória
		newList := make([]*models.Entity, 0, len(re.entities)-1)
		for _, e := range re.entities {
			if e.ID != ent.ID {
				newList = append(newList, e)
			}
		}
		re.entities = newList
		re.selectedEntity = nil
		re.canvas.SetEntities(re.entities)
		re.entityList.Refresh()
		re.lblStatus.SetText(fmt.Sprintf("Entidade '%s' removida.", ent.Name))
	}, re.win)
}

// SaveCurrentRoom persiste as alterações do nome, coordenadas e matriz de tiles no banco SQLite.
func (re *RoomEditorWidget) SaveCurrentRoom() {
	if re.room == nil {
		return
	}

	proj := re.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Nenhum projeto ativo para salvar.", re.win)
		return
	}

	wx, errX := strconv.Atoi(re.entryWorldX.Text)
	wy, errY := strconv.Atoi(re.entryWorldY.Text)
	if errX != nil || errY != nil {
		dialog.ShowError(fmt.Errorf("coordenadas WorldX e WorldY devem ser números inteiros"), re.win)
		return
	}

	re.room.Name = re.entryName.Text
	re.room.WorldX = wx
	re.room.WorldY = wy

	if err := proj.Storage().Rooms.UpdateRoom(context.Background(), re.room); err != nil {
		dialog.ShowError(fmt.Errorf("falha ao salvar sala: %w", err), re.win)
		return
	}

	re.state.NotifyDataChanged()
	if re.onRoomSaved != nil {
		re.onRoomSaved()
	}
	re.lblStatus.SetText(fmt.Sprintf("Sala #%d ('%s') salva com sucesso no SQLite!", re.room.ID, re.room.Name))
}

// RenderTileThumbnail gera uma imagem NRGBA 8x8 do tile para renderização em escala.
func RenderTileThumbnail(t *models.Tile) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	var pat, col []byte
	if t != nil && len(t.PatternBytes) == 8 && len(t.ColorBytes) == 8 {
		pat = t.PatternBytes
		col = t.ColorBytes
	} else {
		pat = make([]byte, 8)
		col = []byte{0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1}
	}

	for y := 0; y < 8; y++ {
		colByte := col[y]
		fg := GetMSXColor((colByte >> 4) & 0x0F)
		bg := GetMSXColor(colByte & 0x0F)
		patByte := pat[y]

		for x := 0; x < 8; x++ {
			bit := (patByte >> uint(7-x)) & 1
			if bit == 1 {
				img.SetNRGBA(x, y, fg)
			} else {
				img.SetNRGBA(x, y, bg)
			}
		}
	}
	return img
}

func behaviorName(b models.BehaviorType) string {
	switch b {
	case models.BehaviorStaticNPC:
		return "NPC Estacionário (0)"
	case models.BehaviorWanderingNPC:
		return "NPC Andarilho (1)"
	case models.BehaviorPatrolNPC:
		return "Patrulha (2)"
	case models.BehaviorChest:
		return "Baú (3)"
	case models.BehaviorDoor:
		return "Porta (4)"
	case models.BehaviorTrigger:
		return "Gatilho (5)"
	case models.BehaviorHostile:
		return "Inimigo Hostil (6)"
	default:
		return fmt.Sprintf("Tipo %d", b)
	}
}
