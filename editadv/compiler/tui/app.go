package tui

import (
	"fmt"
	"os"

	"edadv/compiler"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

// AdventureEditor é o aplicativo TUI estilo Turbo Vision.
type AdventureEditor struct {
	App         *tview.Application
	Pages       *tview.Pages
	TopMenu     *tview.TextView
	StatusBar   *tview.TextView
	CurrentFile string
	Compiler    *compiler.Compiler

	// Abas
	TabMap      *MapTab
	TabObjects  *ObjectsTab
	TabCommands *CommandsTab
	TabMessages *MessagesTab
	TabBuild    *BuildTab
}

// NewAdventureEditor inicializa o editor com o tema clássico Turbo Vision (Azul / Branco / Ciano / Amarelo).
func NewAdventureEditor(filePath string) *AdventureEditor {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Configuração do Tema Turbo Vision
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDarkBlue
	tview.Styles.ContrastBackgroundColor = tcell.ColorNavy
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlue
	tview.Styles.BorderColor = tcell.ColorLightCyan
	tview.Styles.TitleColor = tcell.ColorYellow
	tview.Styles.GraphicsColor = tcell.ColorLightCyan
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorYellow
	tview.Styles.TertiaryTextColor = tcell.ColorGreen

	comp := compiler.NewCompiler()
	if filePath == "" {
		filePath = "games/demo.yaml"
	}

	if _, err := os.Stat(filePath); err == nil {
		_ = comp.LoadFile(filePath)
	} else {
		// Inicializa jogo vazio se não existir
		comp.Game = compiler.AdventureGame{
			Meta: compiler.AdventureMetadata{
				Title:   "Nova Aventura MSX",
				Author:  "Autor",
				Version: "1.0",
			},
			Config: compiler.AdventureConfig{
				InitialPosition: 1,
				MaxCarried:      compiler.DefaultMaxCarried,
				MaxInContainer:  compiler.DefaultMaxInContainer,
			},
		}
	}

	editor := &AdventureEditor{
		App:         app,
		Pages:       pages,
		CurrentFile: filePath,
		Compiler:    comp,
	}

	// Menu Superior (Barra de Funções estilo Turbo Vision)
	editor.TopMenu = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	editor.updateTopMenu("map")

	// Barra de Status Inferior
	editor.StatusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	editor.SetStatus("Pronto. Use [yellow]1-5[white] para alternar abas, [yellow]F2[white] Salvar, [yellow]F5[white] Compilar C, [yellow]F9[white] Gerar ROM.")

	// Criação das Abas
	editor.TabMap = NewMapTab(editor)
	editor.TabObjects = NewObjectsTab(editor)
	editor.TabCommands = NewCommandsTab(editor)
	editor.TabMessages = NewMessagesTab(editor)
	editor.TabBuild = NewBuildTab(editor)

	// Adiciona as abas no gerenciador de páginas
	pages.AddPage("map", editor.TabMap.GetView(), true, true)
	pages.AddPage("objects", editor.TabObjects.GetView(), true, false)
	pages.AddPage("commands", editor.TabCommands.GetView(), true, false)
	pages.AddPage("messages", editor.TabMessages.GetView(), true, false)
	pages.AddPage("build", editor.TabBuild.GetView(), true, false)

	// Layout Principal Vertical: TopMenu (1 linha) + Pages (flex) + StatusBar (1 linha)
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(editor.TopMenu, 1, 0, false).
		AddItem(pages, 0, 1, true).
		AddItem(editor.StatusBar, 1, 0, false)

	// Teclas Globais de Atalho
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			editor.ShowHelp()
			return nil
		case tcell.KeyF2:
			editor.SaveGame()
			return nil
		case tcell.KeyF5:
			editor.SwitchTab("build")
			editor.TabBuild.CompileC()
			return nil
		case tcell.KeyF9:
			editor.SwitchTab("build")
			editor.TabBuild.BuildROM()
			return nil
		case tcell.KeyF10:
			app.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case '1':
				if !editor.isInputFocused() {
					editor.SwitchTab("map")
					return nil
				}
			case '2':
				if !editor.isInputFocused() {
					editor.SwitchTab("objects")
					return nil
				}
			case '3':
				if !editor.isInputFocused() {
					editor.SwitchTab("commands")
					return nil
				}
			case '4':
				if !editor.isInputFocused() {
					editor.SwitchTab("messages")
					return nil
				}
			case '5':
				if !editor.isInputFocused() {
					editor.SwitchTab("build")
					return nil
				}
			}
		}
		return event
	})

	app.SetRoot(mainLayout, true)
	return editor
}

func (e *AdventureEditor) isInputFocused() bool {
	focus := e.App.GetFocus()
	switch focus.(type) {
	case *tview.InputField, *tview.TextArea:
		return true
	default:
		return false
	}
}

func (e *AdventureEditor) updateTopMenu(currentTab string) {
	tabStyle := func(name, key, tabID string) string {
		if tabID == currentTab {
			return fmt.Sprintf("[black:yellow] %s:%s [-:-]", key, name)
		}
		return fmt.Sprintf("[white:blue] %s:%s [-:-]", key, name)
	}

	e.TopMenu.SetText(fmt.Sprintf(
		" %s %s %s %s %s   [yellow]F1[-]:Ajuda [yellow]F2[-]:Salvar [yellow]F5[-]:C [yellow]F9[-]:ROM [yellow]F10[-]:Sair",
		tabStyle("Salas", "1", "map"),
		tabStyle("Objetos", "2", "objects"),
		tabStyle("Comandos", "3", "commands"),
		tabStyle("Mensagens", "4", "messages"),
		tabStyle("Build", "5", "build"),
	))
}

// SwitchTab alterna entre as abas principais.
func (e *AdventureEditor) SwitchTab(tabName string) {
	e.Pages.SwitchToPage(tabName)
	e.updateTopMenu(tabName)
	switch tabName {
	case "map":
		e.App.SetFocus(e.TabMap.GetView())
		e.TabMap.Refresh()
	case "objects":
		e.App.SetFocus(e.TabObjects.GetView())
		e.TabObjects.Refresh()
	case "commands":
		e.App.SetFocus(e.TabCommands.GetView())
		e.TabCommands.Refresh()
	case "messages":
		e.App.SetFocus(e.TabMessages.GetView())
		e.TabMessages.Refresh()
	case "build":
		e.App.SetFocus(e.TabBuild.GetView())
	}
}

// SetStatus define a mensagem de rodapé na barra de status.
func (e *AdventureEditor) SetStatus(msg string) {
	e.StatusBar.SetText(" " + msg)
}

// SaveGame salva a aventura no arquivo YAML.
func (e *AdventureEditor) SaveGame() {
	if err := e.Compiler.Game.Validate(); err != nil {
		e.ShowError("Erro de Validação: " + err.Error())
		return
	}

	if err := e.saveYAMLFile(); err != nil {
		e.ShowError("Erro ao gravar: " + err.Error())
		return
	}

	e.SetStatus(fmt.Sprintf("[green]Aventura salva com sucesso em: %s[-]", e.CurrentFile))
}

func (e *AdventureEditor) saveYAMLFile() error {
	data, err := yaml.Marshal(e.Compiler.Game)
	if err != nil {
		return err
	}
	return os.WriteFile(e.CurrentFile, data, 0644)
}

// ShowHelp exibe modal de ajuda estilo Turbo Vision.
func (e *AdventureEditor) ShowHelp() {
	modal := tview.NewModal().
		SetText("Editor de Adventures - MSX Clean-Room\n\n"+
			"Atalhos de Teclado:\n"+
			"1 a 5   : Alternar entre as abas\n"+
			"Tab     : Navegar entre campos/painéis\n"+
			"Enter   : Editar/Confirmar seleção\n"+
			"F2      : Salvar história em YAML\n"+
			"F5      : Compilar matrizes C (game_data.h/c)\n"+
			"F9      : Compilar ROM do MSX via MSXgl\n"+
			"F10     : Sair do editor").
		AddButtons([]string{"Entendido"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.Pages.HidePage("help_modal")
			e.Pages.RemovePage("help_modal")
		})

	e.Pages.AddPage("help_modal", modal, true, true)
}

// ShowError exibe diálogo modal com mensagem de erro.
func (e *AdventureEditor) ShowError(msg string) {
	modal := tview.NewModal().
		SetText("[red]Erro:[-] " + msg).
		AddButtons([]string{"Fechar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.Pages.HidePage("error_modal")
			e.Pages.RemovePage("error_modal")
		})
	e.Pages.AddPage("error_modal", modal, true, true)
}

// Run inicia o loop do TUI.
func (e *AdventureEditor) Run() error {
	return e.App.Run()
}
