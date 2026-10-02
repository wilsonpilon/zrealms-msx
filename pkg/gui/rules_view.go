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

// RulesView implementa o editor de classes de herói, atributos e itens de RPG.
type RulesView struct {
	state          *ProjectState
	win            fyne.Window
	classList      *widget.List
	itemList       *widget.List
	classes        []*models.HeroClass
	items          []*models.Item
	selectedClass  *models.HeroClass
	selectedItem   *models.Item
	classDetail    *widget.Label
	itemDetail     *widget.Label
	lblStatus      *widget.Label
}

// NewRulesView cria o componente de edição de regras e estatísticas de RPG.
func NewRulesView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &RulesView{
		state:       state,
		win:         win,
		classes:     make([]*models.HeroClass, 0),
		items:       make([]*models.Item, 0),
		classDetail: widget.NewLabel("Selecione uma classe para ver os atributos"),
		itemDetail:  widget.NewLabel("Selecione um item para ver as estatísticas"),
		lblStatus:   widget.NewLabel("Pronto."),
	}

	// 1. Lista de Classes de Heróis
	v.classList = widget.NewList(
		func() int { return len(v.classes) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Classe"))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.classes) {
				return
			}
			c := v.classes[id]
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText(
				fmt.Sprintf("#%d: %s (HP: %d, ATK: %d)", c.ID, c.Name, c.BaseHP, c.BaseAtk),
			)
		},
	)

	v.classList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.classes) {
			return
		}
		v.selectedClass = v.classes[id]
		v.updateClassDetails()
	}

	// 2. Lista de Itens e Equipamentos
	v.itemList = widget.NewList(
		func() int { return len(v.items) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.StorageIcon()), widget.NewLabel("Item"))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.items) {
				return
			}
			it := v.items[id]
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText(
				fmt.Sprintf("#%d: %s [%s] (%d GP)", it.ID, it.Name, itemTypeName(it.ItemType), it.Price),
			)
		},
	)

	v.itemList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.items) {
			return
		}
		v.selectedItem = v.items[id]
		v.updateItemDetails()
	}

	// Botões de Ação para Classes
	btnAddClass := widget.NewButtonWithIcon("Nova", theme.ContentAddIcon(), func() {
		v.showClassForm(nil)
	})
	btnEditClass := widget.NewButtonWithIcon("Editar", theme.DocumentCreateIcon(), func() {
		if v.selectedClass != nil {
			v.showClassForm(v.selectedClass)
		} else {
			dialog.ShowInformation("Aviso", "Selecione uma classe para editar.", v.win)
		}
	})
	btnDeleteClass := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedClass()
	})

	classHeader := container.NewVBox(
		widget.NewLabelWithStyle("⚔️ Classes de Heróis", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAddClass, btnEditClass, btnDeleteClass),
		widget.NewSeparator(),
	)
	classPanel := container.NewBorder(classHeader, nil, nil, nil, v.classList)

	// Botões de Ação para Itens
	btnAddItem := widget.NewButtonWithIcon("Novo", theme.ContentAddIcon(), func() {
		v.showItemForm(nil)
	})
	btnEditItem := widget.NewButtonWithIcon("Editar", theme.DocumentCreateIcon(), func() {
		if v.selectedItem != nil {
			v.showItemForm(v.selectedItem)
		} else {
			dialog.ShowInformation("Aviso", "Selecione um item para editar.", v.win)
		}
	})
	btnDeleteItem := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedItem()
	})

	itemHeader := container.NewVBox(
		widget.NewLabelWithStyle("🛡️ Itens & Equipamentos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAddItem, btnEditItem, btnDeleteItem),
		widget.NewSeparator(),
	)
	itemPanel := container.NewBorder(itemHeader, nil, nil, nil, v.itemList)

	// Painel de Detalhes
	cardClassDetail := widget.NewCard(
		"Ficha da Classe",
		"Parâmetros base no motor Z80",
		v.classDetail,
	)

	cardItemDetail := widget.NewCard(
		"Propriedades do Item",
		"Modificadores de combate e inventário",
		v.itemDetail,
	)

	detailPanel := container.NewVBox(
		cardClassDetail,
		widget.NewSeparator(),
		cardItemDetail,
		widget.NewSeparator(),
		v.lblStatus,
	)

	leftSplit := container.NewHSplit(classPanel, itemPanel)
	leftSplit.Offset = 0.50

	mainSplit := container.NewHSplit(leftSplit, detailPanel)
	mainSplit.Offset = 0.65

	state.OnProjectLoaded(func(_ *project.Project, _ string) { v.reload() })
	state.OnProjectClosed(func() {
		v.classes = make([]*models.HeroClass, 0)
		v.items = make([]*models.Item, 0)
		v.selectedClass = nil
		v.selectedItem = nil
		v.classList.Refresh()
		v.itemList.Refresh()
		v.classDetail.SetText("Nenhum projeto aberto.")
		v.itemDetail.SetText("Nenhum projeto aberto.")
	})

	if state.IsOpen() {
		v.reload()
	}

	return mainSplit
}

func (v *RulesView) reload() {
	proj := v.state.Current()
	if proj == nil {
		return
	}
	classes, err := proj.Storage().GameData.ListHeroClasses(context.Background())
	if err == nil {
		v.classes = classes
		v.classList.Refresh()
		if len(v.classes) > 0 && v.selectedClass == nil {
			v.classList.Select(0)
		}
	}
	items, err := proj.Storage().GameData.ListItems(context.Background())
	if err == nil {
		v.items = items
		v.itemList.Refresh()
		if len(v.items) > 0 && v.selectedItem == nil {
			v.itemList.Select(0)
		}
	}
}

func (v *RulesView) updateClassDetails() {
	if v.selectedClass == nil {
		v.classDetail.SetText("Nenhuma classe selecionada.")
		return
	}
	c := v.selectedClass
	info := fmt.Sprintf(
		"ID: #%d\nNome: %s\n\n"+
			"Estatísticas Base:\n"+
			"  HP Inicial:     %d pts\n"+
			"  MP Inicial:     %d pts\n"+
			"  Ataque Base:    %d pts\n"+
			"  Defesa Base:    %d pts\n\n"+
			"Armazenamento MSX: 1 registro de 8 bytes na tabela de classes da ROM/Mapper.",
		c.ID, c.Name, c.BaseHP, c.BaseMP, c.BaseAtk, c.BaseDef,
	)
	v.classDetail.SetText(info)
}

func (v *RulesView) updateItemDetails() {
	if v.selectedItem == nil {
		v.itemDetail.SetText("Nenhum item selecionado.")
		return
	}
	it := v.selectedItem
	info := fmt.Sprintf(
		"ID: #%d\nNome: %s\nTipo: %s\nPreço na Loja: %d Moedas de Ouro (GP)\n\n"+
			"Efeito no Personagem:\n"+
			"  Modificador: +%d em %s\n\n"+
			"Tamanho na Memória: 8 bytes por registro de item no cartucho.",
		it.ID, it.Name, itemTypeName(it.ItemType), it.Price,
		it.ModifierValue, statName(it.ModifierStat),
	)
	v.itemDetail.SetText(info)
}

func (v *RulesView) showClassForm(c *models.HeroClass) {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	isNew := c == nil
	if isNew {
		c = &models.HeroClass{
			Name:    "Guerreiro",
			BaseHP:  50,
			BaseMP:  10,
			BaseAtk: 12,
			BaseDef: 8,
		}
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(c.Name)
	hpEntry := widget.NewEntry()
	hpEntry.SetText(fmt.Sprintf("%d", c.BaseHP))
	mpEntry := widget.NewEntry()
	mpEntry.SetText(fmt.Sprintf("%d", c.BaseMP))
	atkEntry := widget.NewEntry()
	atkEntry.SetText(fmt.Sprintf("%d", c.BaseAtk))
	defEntry := widget.NewEntry()
	defEntry.SetText(fmt.Sprintf("%d", c.BaseDef))

	items := []*widget.FormItem{
		widget.NewFormItem("Nome da Classe", nameEntry),
		widget.NewFormItem("HP Inicial", hpEntry),
		widget.NewFormItem("MP Inicial", mpEntry),
		widget.NewFormItem("Ataque Base", atkEntry),
		widget.NewFormItem("Defesa Base", defEntry),
	}

	title := "Adicionar Nova Classe"
	if !isNew {
		title = fmt.Sprintf("Editar Classe #%d", c.ID)
	}

	dialog.ShowForm(title, "Salvar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		hp, errHP := strconv.Atoi(hpEntry.Text)
		mp, errMP := strconv.Atoi(mpEntry.Text)
		atk, errAtk := strconv.Atoi(atkEntry.Text)
		def, errDef := strconv.Atoi(defEntry.Text)
		if errHP != nil || errMP != nil || errAtk != nil || errDef != nil {
			dialog.ShowError(fmt.Errorf("atributos devem ser números inteiros"), v.win)
			return
		}

		c.Name = nameEntry.Text
		c.BaseHP = hp
		c.BaseMP = mp
		c.BaseAtk = atk
		c.BaseDef = def

		if isNew {
			if err := proj.Storage().GameData.CreateHeroClass(context.Background(), c); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText(fmt.Sprintf("Classe '%s' criada com sucesso.", c.Name))
		} else {
			if err := proj.Storage().GameData.UpdateHeroClass(context.Background(), c); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText(fmt.Sprintf("Classe '%s' atualizada com sucesso.", c.Name))
		}

		v.state.NotifyDataChanged()
		v.reload()
		v.updateClassDetails()
	}, v.win)
}

func (v *RulesView) deleteSelectedClass() {
	if v.selectedClass == nil {
		dialog.ShowInformation("Aviso", "Selecione uma classe para excluir.", v.win)
		return
	}

	proj := v.state.Current()
	if proj == nil {
		return
	}

	c := v.selectedClass
	dialog.ShowConfirm("Excluir Classe", fmt.Sprintf("Deseja realmente excluir a classe '%s' (#%d)?", c.Name, c.ID), func(ok bool) {
		if !ok {
			return
		}
		if err := proj.Storage().GameData.DeleteHeroClass(context.Background(), c.ID); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.selectedClass = nil
		v.state.NotifyDataChanged()
		v.reload()
		v.lblStatus.SetText(fmt.Sprintf("Classe '%s' excluída.", c.Name))
	}, v.win)
}

func (v *RulesView) showItemForm(item *models.Item) {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	isNew := item == nil
	if isNew {
		item = &models.Item{
			Name:          "Espada de Ferro",
			ItemType:      models.ItemWeapon,
			ModifierStat:  3, // Ataque
			ModifierValue: 5,
			Price:         100,
		}
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(item.Name)

	typeOptions := []string{
		"0 - Arma",
		"1 - Armadura",
		"2 - Consumível",
		"3 - Chave",
		"4 - Quest",
	}
	typeSelect := widget.NewSelect(typeOptions, nil)
	typeSelect.SetSelectedIndex(int(item.ItemType))

	statOptions := []string{
		"0 - Nenhum",
		"1 - HP Máximo",
		"2 - MP Máximo",
		"3 - Ataque",
		"4 - Defesa",
	}
	statSelect := widget.NewSelect(statOptions, nil)
	statSelect.SetSelectedIndex(item.ModifierStat)

	valEntry := widget.NewEntry()
	valEntry.SetText(fmt.Sprintf("%d", item.ModifierValue))

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%d", item.Price))

	items := []*widget.FormItem{
		widget.NewFormItem("Nome do Item", nameEntry),
		widget.NewFormItem("Tipo de Item", typeSelect),
		widget.NewFormItem("Atributo Modificado", statSelect),
		widget.NewFormItem("Valor do Modificador", valEntry),
		widget.NewFormItem("Preço (GP)", priceEntry),
	}

	title := "Adicionar Novo Item"
	if !isNew {
		title = fmt.Sprintf("Editar Item #%d", item.ID)
	}

	dialog.ShowForm(title, "Salvar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		val, errVal := strconv.Atoi(valEntry.Text)
		price, errPrice := strconv.Atoi(priceEntry.Text)
		if errVal != nil || errPrice != nil {
			dialog.ShowError(fmt.Errorf("valor e preço devem ser números inteiros"), v.win)
			return
		}

		item.Name = nameEntry.Text
		item.ItemType = models.ItemType(typeSelect.SelectedIndex())
		item.ModifierStat = statSelect.SelectedIndex()
		item.ModifierValue = val
		item.Price = price

		if isNew {
			if err := proj.Storage().GameData.CreateItem(context.Background(), item); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText(fmt.Sprintf("Item '%s' criado com sucesso.", item.Name))
		} else {
			if err := proj.Storage().GameData.UpdateItem(context.Background(), item); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText(fmt.Sprintf("Item '%s' atualizado com sucesso.", item.Name))
		}

		v.state.NotifyDataChanged()
		v.reload()
		v.updateItemDetails()
	}, v.win)
}

func (v *RulesView) deleteSelectedItem() {
	if v.selectedItem == nil {
		dialog.ShowInformation("Aviso", "Selecione um item para excluir.", v.win)
		return
	}

	proj := v.state.Current()
	if proj == nil {
		return
	}

	item := v.selectedItem
	dialog.ShowConfirm("Excluir Item", fmt.Sprintf("Deseja realmente excluir o item '%s' (#%d)?", item.Name, item.ID), func(ok bool) {
		if !ok {
			return
		}
		if err := proj.Storage().GameData.DeleteItem(context.Background(), item.ID); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.selectedItem = nil
		v.state.NotifyDataChanged()
		v.reload()
		v.lblStatus.SetText(fmt.Sprintf("Item '%s' excluído.", item.Name))
	}, v.win)
}

func itemTypeName(t models.ItemType) string {
	switch t {
	case models.ItemWeapon:
		return "Arma"
	case models.ItemArmor:
		return "Armadura"
	case models.ItemConsumable:
		return "Consumível"
	case models.ItemKey:
		return "Chave"
	case models.ItemQuest:
		return "Quest"
	default:
		return "Geral"
	}
}

func statName(s int) string {
	switch s {
	case 1:
		return "HP Máximo"
	case 2:
		return "MP Máximo"
	case 3:
		return "Ataque"
	case 4:
		return "Defesa"
	default:
		return "Nenhum"
	}
}
