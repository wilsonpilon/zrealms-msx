package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/version"
)

// App encapsula a aplicação gráfica do Z-Realm construída sobre o framework Fyne v2.
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	state   *ProjectState
}

// NewApp inicializa a aplicação desktop com o tema retrô escuro e gerenciamento de estado.
func NewApp() *App {
	a := app.NewWithID("com.zrealm.msx.editor")
	a.Settings().SetTheme(NewRetroDarkTheme())

	w := a.NewWindow(fmt.Sprintf("Z-Realm — MSX 2 cRPG Construction Kit (%s)", version.String()))
	w.Resize(fyne.NewSize(1100, 720))
	w.CenterOnScreen()

	state := NewProjectState()
	mainLayout := NewMainLayout(state, w)
	statusBar := NewStatusBar(state)

	mainMenu := BuildMainMenu(state, w, func() {
		mainLayout.SelectExportTab()
	})
	w.SetMainMenu(mainMenu)

	// Atualiza o título da janela conforme o projeto
	state.OnProjectLoaded(func(_ *project.Project, path string) {
		w.SetTitle(fmt.Sprintf("Z-Realm — %s (%s)", path, version.String()))
	})

	state.OnProjectClosed(func() {
		w.SetTitle(fmt.Sprintf("Z-Realm — MSX 2 cRPG Construction Kit (%s)", version.String()))
	})

	state.OnModifiedChanged(func(mod bool) {
		p := state.FilePath()
		if p == "" {
			p = "Sem Projeto"
		}
		star := ""
		if mod {
			star = " *"
		}
		w.SetTitle(fmt.Sprintf("Z-Realm — %s%s (%s)", p, star, version.String()))
	})

	root := container.NewBorder(
		nil,
		statusBar.Container(),
		nil,
		nil,
		mainLayout.Container(),
	)

	w.SetContent(root)

	return &App{
		fyneApp: a,
		window:  w,
		state:   state,
	}
}

// Run executa o loop de eventos principal da interface gráfica.
func (a *App) Run() {
	a.window.ShowAndRun()
}

// State retorna o gerenciador de estado do projeto.
func (a *App) State() *ProjectState {
	return a.state
}
