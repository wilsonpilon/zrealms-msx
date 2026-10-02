package tui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rivo/tview"
)

// BuildTab gerencia a compilação para matrizes C e o build da ROM do MSX.
type BuildTab struct {
	Editor   *AdventureEditor
	View     *tview.Flex
	Controls *tview.Form
	LogView  *tview.TextView
}

func NewBuildTab(editor *AdventureEditor) *BuildTab {
	tab := &BuildTab{
		Editor: editor,
	}

	tab.Controls = tview.NewForm()
	tab.Controls.SetBorder(true).
		SetTitle("[ Opções de Compilação & Build ]")

	tab.LogView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	tab.LogView.SetBorder(true).
		SetTitle("[ Log de Compilação & MSXgl ]")

	tab.Controls.AddButton("1. Compilar para C (game_data.h/c) [F5]", func() {
		tab.CompileC()
	})

	tab.Controls.AddButton("2. Construir ROM MSX 32K (build.bat) [F9]", func() {
		tab.BuildROM()
	})

	tab.View = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tab.Controls, 7, 0, false).
		AddItem(tab.LogView, 0, 1, true)

	tab.LogView.SetText("[yellow]Pronto.[-] Clique nos botões acima ou use [yellow]F5[-] para compilar C e [yellow]F9[-] para construir a ROM.")
	return tab
}

func (tab *BuildTab) GetView() tview.Primitive {
	return tab.View
}

// CompileC gera o código C estático a partir da aventura atual.
func (tab *BuildTab) CompileC() {
	tab.LogView.Clear()
	tab.log("[cyan]>>> Compilando aventura para matrizes C (game_data.h / game_data.c)...[-]")

	if err := tab.Editor.Compiler.Game.Validate(); err != nil {
		tab.log(fmt.Sprintf("[red]Erro de validação: %v[-]", err))
		return
	}

	// 1. Gera na pasta engine/src
	outDir := "engine/src"
	if err := tab.Editor.Compiler.GenerateCData(outDir); err != nil {
		tab.log(fmt.Sprintf("[red]Erro ao gerar C: %v[-]", err))
		return
	}
	tab.log(fmt.Sprintf("[green]Sucesso! Matrizes geradas em %s/game_data.h e .c[-]", outDir))

	// 2. Copia para a pasta do projeto MSXgl se existir
	adventDir := "MSXgl/projects/advent"
	if _, err := os.Stat(adventDir); err == nil {
		copyFile(filepath.Join(outDir, "game_data.h"), filepath.Join(adventDir, "game_data.h"))
		copyFile(filepath.Join(outDir, "game_data.c"), filepath.Join(adventDir, "game_data.c"))
		tab.log("[green]Matrizes copiadas para o projeto MSXgl/projects/advent/.[-]")
	}

	tab.Editor.SetStatus("[green]Compilação C concluída com sucesso![-]")
}

// BuildROM executa o build.bat do MSXgl para produzir a ROM de 32KB.
func (tab *BuildTab) BuildROM() {
	tab.CompileC()

	tab.log("[cyan]>>> Iniciando build da ROM MSX (MSXgl / SDCC)...[-]")
	adventDir := "MSXgl/projects/advent"

	// Executa build.bat dentro da pasta MSXgl/projects/advent
	cmd := exec.Command("cmd", "/c", "build.bat")
	cmd.Dir = adventDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		tab.log(fmt.Sprintf("[red]Erro no build do MSXgl: %v[-]\n%s", err, string(output)))
		tab.Editor.SetStatus("[red]Falha no build da ROM! Verifique o log.[-]")
		return
	}

	tab.log(fmt.Sprintf("[green]Build concluído com sucesso![-]\n%s", string(output)))

	romPath := filepath.Join(adventDir, "out", "advent.rom")
	if fi, err := os.Stat(romPath); err == nil {
		tab.log(fmt.Sprintf("[yellow]Artefato ROM gerado: %s (%d bytes)[-]", romPath, fi.Size()))
		tab.Editor.SetStatus(fmt.Sprintf("[green]ROM gerada: advent.rom (%d bytes)[-]", fi.Size()))
	}
}

func (tab *BuildTab) log(msg string) {
	current := tab.LogView.GetText(false)
	if current != "" {
		tab.LogView.SetText(current + "\n" + msg)
	} else {
		tab.LogView.SetText(msg)
	}
	tab.LogView.ScrollToEnd()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
