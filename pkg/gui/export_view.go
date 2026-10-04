package gui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/runner"
)

// ExportView implementa o painel de exportação, empacotamento de disquete e execução no openMSX.
type ExportView struct {
	state         *ProjectState
	win           fyne.Window
	entryOutDir   *widget.Entry
	checkBanks    *widget.Check
	selectMachine *widget.Select
	entryEmuPath  *widget.Entry
	logArea       *widget.Entry
}

// NewExportView cria a aba de exportação e automação One-Click Run para MSX 2.
func NewExportView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &ExportView{
		state:        state,
		win:          win,
		entryOutDir:  widget.NewEntry(),
		checkBanks:   widget.NewCheck("Gerar arquivos individuais de banco (SEGxx.BNK)", nil),
		entryEmuPath: widget.NewEntry(),
		logArea:      widget.NewMultiLineEntry(),
	}

	v.entryOutDir.SetPlaceHolder("Deixe vazio para usar subpasta ./build_msx")
	v.checkBanks.SetChecked(true)
	v.entryEmuPath.SetPlaceHolder("Deixe vazio para auto-detectar openmsx no PATH / Scoop")

	machines := []string{
		"Philips_NMS_8250 (MSX 2 + 512KB Mapper) [Padrão]",
		"Panasonic_FS-A1GT (MSX 2+ / MSX turbo R)",
	}
	v.selectMachine = widget.NewSelect(machines, nil)
	v.selectMachine.SetSelectedIndex(0)

	v.logArea.Disable()
	v.logArea.SetText("Nenhuma ação realizada nesta sessão.\nClique em 'Testar no openMSX' para o ciclo completo automatizado\nou em 'Exportar Binários' para geração manual.")

	btnRunOneClick := widget.NewButtonWithIcon("▶ Testar no openMSX (One-Click Run) [F5]", theme.MediaPlayIcon(), func() {
		v.runOneClick()
	})
	btnRunOneClick.Importance = widget.HighImportance

	btnExport := widget.NewButtonWithIcon("Exportar Binários (Manual)", theme.DocumentSaveIcon(), func() {
		v.runExport()
	})

	btnOpenFolder := widget.NewButtonWithIcon("Abrir Pasta no Explorer", theme.FolderOpenIcon(), func() {
		outDir := v.entryOutDir.Text
		if outDir == "" && state.IsOpen() {
			outDir = filepath.Join(filepath.Dir(state.FilePath()), "build_msx")
		}
		if outDir != "" {
			_ = exec.Command("explorer.exe", outDir).Start()
		}
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle("🚀 Automação 'One-Click Run' (openMSX)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Exporta o projeto SQLite, empacota o disquete MSX-DOS 2 (720 KB) e executa o jogo no openMSX com 1 clique."),
		container.NewGridWithColumns(2,
			widget.NewFormItem("Perfil da Máquina:", v.selectMachine).Widget,
			widget.NewFormItem("Caminho do openMSX (Opcional):", v.entryEmuPath).Widget,
		),
		container.NewHBox(btnRunOneClick),
		widget.NewSeparator(),

		widget.NewLabelWithStyle("📦 Exportador Binário Manual", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewFormItem("Diretório de Saída:", v.entryOutDir).Widget,
		v.checkBanks,
		container.NewHBox(btnExport, btnOpenFolder),
		widget.NewSeparator(),

		widget.NewLabelWithStyle("Registro de Operações:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	return container.NewBorder(form, nil, nil, nil, v.logArea)
}

func (v *ExportView) runOneClick() {
	if !v.state.IsOpen() {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto antes de testar.", v.win)
		return
	}

	machine := "Philips_NMS_8250"
	if strings.HasPrefix(v.selectMachine.Selected, "Panasonic") {
		machine = "Panasonic_FS-A1GT"
	}

	outDir := v.entryOutDir.Text
	emuPath := strings.TrimSpace(v.entryEmuPath.Text)

	v.logArea.SetText("Iniciando pipeline One-Click Run...\n[1/3] Exportando dados do projeto SQLite...")

	go func() {
		res, err := v.state.OneClickRun(runner.RunOptions{
			OutputDir:   outDir,
			EmulatorExe: emuPath,
			Machine:     machine,
			AutoRun:     true,
			Async:       true,
		})

		if err != nil {
			v.logArea.SetText(fmt.Sprintf("FALHA NO ONE-CLICK RUN:\n%v", err))
			dialog.ShowError(fmt.Errorf("falha ao executar no openMSX: %w", err), v.win)
			return
		}

		summary := fmt.Sprintf(
			"==========================================================\n"+
				"   ONE-CLICK RUN DISPARADO COM SUCESSO!\n"+
				"==========================================================\n\n"+
				"[1/3] Exportação Binária: OK\n"+
				"      - Tabela Mestra:    %s\n"+
				"      - Dados (GAME.DAT): %d bytes (%d segmentos de 16KB)\n"+
				"      - Recursos:         %d catalogados\n\n"+
				"[2/3] Imagem de Disquete (.DSK): OK\n"+
				"      - Arquivo DSK:      %s (%d bytes / 720 KB)\n"+
				"      - Formato:          MSX-DOS 2 FAT12 (Bootável)\n\n"+
				"[3/3] Emulador openMSX: DISPARADO!\n"+
				"      - Executável:       %s\n"+
				"      - Máquina:          %s\n"+
				"      - Linha de Comando:\n        %s\n\n"+
				"O jogo está rodando no openMSX!\n"+
				"Controles: Direcionais para andar, ESPACO para interagir, ESC para sair ao DOS.",
			res.ExportResult.HeaderPath, res.ExportResult.TotalDataBytes, res.ExportResult.TotalSegments, res.ExportResult.ResourceCount,
			res.DskPath, res.DskSize,
			res.EmulatorExe, machine,
			strings.Join(res.CommandLine, " "),
		)
		v.logArea.SetText(summary)
	}()
}

func (v *ExportView) runExport() {
	if !v.state.IsOpen() {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto antes de exportar.", v.win)
		return
	}

	outDir := v.entryOutDir.Text
	res, err := v.state.Export(outDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf("falha na exportação: %w", err), v.win)
		v.logArea.SetText(fmt.Sprintf("ERRO NA EXPORTAÇÃO:\n%v", err))
		return
	}

	summary := fmt.Sprintf(
		"==========================================================\n"+
			"   EXPORTAÇÃO CONCLUÍDA COM SUCESSO!\n"+
			"==========================================================\n\n"+
			"Tabela Mestra:    %s\n"+
			"Arquivo de Dados: %s\n"+
			"Volume de Dados:  %d bytes\n"+
			"Segmentos 16 KB:  %d segmento(s) de Memory Mapper alocados\n"+
			"Total Recursos:   %d catalogados\n\n"+
			"Os arquivos estão prontos para cópia para o disquete MSX-DOS 2 (720 KB)\n"+
			"ou montagem via emulador openMSX!",
		res.HeaderPath, res.DataPath, res.TotalDataBytes, res.TotalSegments, res.ResourceCount,
	)
	v.logArea.SetText(summary)
	dialog.ShowInformation("Exportação Concluída", "Os binários MSX 2 foram gerados com sucesso!", v.win)
}
