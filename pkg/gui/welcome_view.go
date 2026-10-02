package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/version"
)

// NewWelcomeView constrói a tela inicial com ações rápidas para quando nenhum projeto estiver aberto.
func NewWelcomeView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("⚔️ Z-REALM (MSX 2)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabelWithStyle(
		fmt.Sprintf("O ZZT dos cRPGs para MSX-DOS 2 & V9938 • Versão %s", version.String()),
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	banner := widget.NewLabelWithStyle(
		"Crie mapas SCREEN 4 (32x18 tiles), sprites Modo 2 (16x16), gerencie segmentos de 16 KB no Memory Mapper\ne exporte diretamente para MSX-DOS 2 com um único clique!",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	btnNew := widget.NewButtonWithIcon("Criar Novo Projeto (.rpgproj)", theme.DocumentCreateIcon(), func() {
		entryName := widget.NewEntry()
		entryName.SetText("Meu Novo RPG")

		items := []*widget.FormItem{
			widget.NewFormItem("Nome do Projeto", entryName),
		}

		dialog.ShowForm("Novo Projeto Z-Realm", "Criar", "Cancelar", items, func(confirmed bool) {
			if !confirmed {
				return
			}
			fileSave := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil {
					return
				}
				path := uc.URI().Path()
				_ = uc.Close()
				if err := state.NewProject(path, entryName.Text); err != nil {
					dialog.ShowError(err, win)
				}
			}, win)
			fileSave.SetFilter(storage.NewExtensionFileFilter([]string{".rpgproj"}))
			fileSave.SetFileName("projeto.rpgproj")
			fileSave.Show()
		}, win)
	})
	btnNew.Importance = widget.HighImportance

	btnOpen := widget.NewButtonWithIcon("Abrir Projeto Existente...", theme.FolderOpenIcon(), func() {
		fileOpen := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			path := uc.URI().Path()
			_ = uc.Close()
			if err := state.OpenProject(path); err != nil {
				dialog.ShowError(err, win)
			}
		}, win)
		fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".rpgproj", ".db", ".sqlite"}))
		fileOpen.Show()
	})

	btnDemo := widget.NewButtonWithIcon("Abrir Masmorra Demo (2 Salas + Tileset)", theme.MediaPlayIcon(), func() {
		if err := state.OpenDemoProject(""); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Demo Carregada", "A Masmorra Demo de 2 salas foi gerada e carregada com sucesso!\nVocê já pode editar tiles, salas ou exportar para MSX 2.", win)
		}
	})
	btnDemo.Importance = widget.SuccessImportance

	actionsBox := container.NewVBox(
		btnNew,
		btnOpen,
		btnDemo,
	)

	contentCard := widget.NewCard(
		"Painel de Início Rápido",
		"Escolha uma ação para começar seu cRPG para MSX 2",
		container.NewPadded(actionsBox),
	)

	specsInfo := widget.NewRichTextFromMarkdown(`
### Especificações de Destino:
* **Plataforma:** MSX 2 ou superior (Z80A @ 3.58 MHz)
* **Sistema Operacional:** MSX-DOS 2 (suporte a arquivos e paginação dinâmica)
* **Memória Mínima:** 128 KB VRAM + 256 KB Memory Mapper
* **Modo de Vídeo:** V9938 SCREEN 4 (Graphic 3 - 256x192, 16 cores)
* **Dimensões da Sala:** 32 x 18 tiles (256 x 144 pixels) + 48 pixels de HUD/Diálogo
`)

	mainLayout := container.NewCenter(
		container.NewVBox(
			title,
			subtitle,
			banner,
			widget.NewSeparator(),
			contentCard,
			widget.NewSeparator(),
			specsInfo,
		),
	)

	return mainLayout
}
