package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/zrealm-msx/zrealm/pkg/project"
)

// MainLayout coordena a alternância entre a tela de boas-vindas e as abas do projeto ativo.
type MainLayout struct {
	state     *ProjectState
	win       fyne.Window
	container *fyne.Container
	tabs      *container.AppTabs
	tabExport *container.TabItem
}

// NewMainLayout constrói o gerenciador de layout principal com navegação por abas.
func NewMainLayout(state *ProjectState, win fyne.Window) *MainLayout {
	ml := &MainLayout{
		state:     state,
		win:       win,
		container: container.NewStack(),
	}

	welcomeView := NewWelcomeView(state, win)

	// Criação das abas do editor
	roomTab := container.NewTabItemWithIcon("Salas (32x18)", theme.FolderIcon(), NewRoomView(state, win))
	tileTab := container.NewTabItemWithIcon("Tilesets (8x8)", theme.GridIcon(), NewTileView(state, win))
	spriteTab := container.NewTabItemWithIcon("Sprites (16x16)", theme.RadioButtonCheckedIcon(), NewSpriteView(state, win))
	scriptTab := container.NewTabItemWithIcon("Diálogos & Roteiros", theme.DocumentIcon(), NewScriptView(state, win))
	rulesTab := container.NewTabItemWithIcon("Regras & RPG Stats", theme.AccountIcon(), NewRulesView(state, win))
	exportTab := container.NewTabItemWithIcon("Exportador MSX 2", theme.MediaPlayIcon(), NewExportView(state, win))

	ml.tabExport = exportTab
	ml.tabs = container.NewAppTabs(
		roomTab,
		tileTab,
		spriteTab,
		scriptTab,
		rulesTab,
		exportTab,
	)
	ml.tabs.SetTabLocation(container.TabLocationTop)

	// Registra callbacks de alternância de estado
	state.OnProjectLoaded(func(_ *project.Project, _ string) {
		ml.showEditor()
	})

	state.OnProjectClosed(func() {
		ml.showWelcome(welcomeView)
	})

	if state.IsOpen() {
		ml.showEditor()
	} else {
		ml.showWelcome(welcomeView)
	}

	return ml
}

func (ml *MainLayout) showWelcome(welcomeView fyne.CanvasObject) {
	ml.container.Objects = []fyne.CanvasObject{welcomeView}
	ml.container.Refresh()
}

func (ml *MainLayout) showEditor() {
	ml.container.Objects = []fyne.CanvasObject{ml.tabs}
	ml.container.Refresh()
}

// SelectExportTab seleciona programaticamente a aba de exportação.
func (ml *MainLayout) SelectExportTab() {
	if ml.tabs != nil && ml.tabExport != nil {
		ml.tabs.Select(ml.tabExport)
	}
}

// Container retorna o container visual principal.
func (ml *MainLayout) Container() fyne.CanvasObject {
	return ml.container
}
