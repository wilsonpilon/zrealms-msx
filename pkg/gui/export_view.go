package gui

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ExportView implementa o painel de exportação e compilação binária para o MSX 2.
type ExportView struct {
	state       *ProjectState
	win         fyne.Window
	entryOutDir *widget.Entry
	checkBanks  *widget.Check
	logArea     *widget.Entry
}

// NewExportView cria a aba de exportação binária para o MSX 2.
func NewExportView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &ExportView{
		state:       state,
		win:         win,
		entryOutDir: widget.NewEntry(),
		checkBanks:  widget.NewCheck("Gerar arquivos individuais de banco (SEGxx.BNK)", nil),
		logArea:     widget.NewMultiLineEntry(),
	}

	v.entryOutDir.SetPlaceHolder("Deixe vazio para usar subpasta ./build_msx")
	v.checkBanks.SetChecked(true)
	v.logArea.Disable()
	v.logArea.SetText("Nenhuma exportação realizada nesta sessão.\nClique no botão abaixo para gerar os binários MSX 2.")

	btnExport := widget.NewButtonWithIcon("Exportar Masmorra para MSX 2 (V9938/Z80)", theme.MediaPlayIcon(), func() {
		v.runExport()
	})
	btnExport.Importance = widget.HighImportance

	btnOpenFolder := widget.NewButtonWithIcon("Abrir Pasta de Saída no Explorer", theme.FolderOpenIcon(), func() {
		outDir := v.entryOutDir.Text
		if outDir == "" && state.IsOpen() {
			outDir = filepath.Join(filepath.Dir(state.FilePath()), "build_msx")
		}
		if outDir != "" {
			_ = exec.Command("explorer.exe", outDir).Start()
		}
	})

	form := container.NewVBox(
		widget.NewLabelWithStyle("📦 Exportador Binário MSX-DOS 2", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Converte todas as salas, tilesets, padrões V9938 e tabelas do SQLite em arquivos binários planos otimizados."),
		widget.NewFormItem("Diretório de Saída:", v.entryOutDir).Widget,
		v.checkBanks,
		container.NewHBox(btnExport, btnOpenFolder),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Registro de Exportação:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	return container.NewBorder(form, nil, nil, nil, v.logArea)
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
