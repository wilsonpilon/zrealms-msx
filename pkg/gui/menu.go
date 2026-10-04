package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/version"
)

// BuildMainMenu constrói o menu principal do Z-Realm.
func BuildMainMenu(state *ProjectState, win fyne.Window, onExportTabSelect func()) *fyne.MainMenu {
	menuFile := fyne.NewMenu("Arquivo",
		fyne.NewMenuItem("Novo Projeto...", func() {
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
		}),
		fyne.NewMenuItem("Abrir Projeto...", func() {
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
		}),
		fyne.NewMenuItem("Abrir Masmorra Demo", func() {
			if err := state.OpenDemoProject(""); err != nil {
				dialog.ShowError(err, win)
			} else {
				dialog.ShowInformation("Demo Aberta", "A masmorra de teste com 2 salas e tileset padrão foi carregada!", win)
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exportar para MSX 2...", func() {
			if !state.IsOpen() {
				dialog.ShowInformation("Aviso", "Abra um projeto antes de exportar.", win)
				return
			}
			if onExportTabSelect != nil {
				onExportTabSelect()
			}
		}),
		fyne.NewMenuItem("Fechar Projeto", func() {
			state.CloseProject()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sair", func() {
			win.Close()
		}),
	)

	menuTools := fyne.NewMenu("Ferramentas",
		fyne.NewMenuItem("Testar no openMSX (One-Click Run)...", func() {
			if !state.IsOpen() {
				dialog.ShowInformation("Aviso", "Abra um projeto antes de testar.", win)
				return
			}
			if onExportTabSelect != nil {
				onExportTabSelect()
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Validar Integridade", func() {
			proj := state.Current()
			if proj == nil {
				dialog.ShowInformation("Aviso", "Nenhum projeto aberto.", win)
				return
			}
			if err := proj.ValidateIntegrity(); err != nil {
				dialog.ShowError(fmt.Errorf("falha de integridade: %w", err), win)
			} else {
				dialog.ShowInformation("Integridade OK", "O banco de dados SQLite e todas as chaves estrangeiras estão 100% íntegras!", win)
			}
		}),
	)

	menuHelp := fyne.NewMenu("Ajuda",
		fyne.NewMenuItem("Sobre o Z-Realm", func() {
			info := fmt.Sprintf(
				"Z-Realm (zrealm-msx) - %s\n\n"+
					"O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2)\n"+
					"Inspirado nos clássicos cRPGs dos anos 80/90 e na flexibilidade do ZZT de Tim Sweeney.\n\n"+
					"Tecnologias: Go 1.27 + Fyne v2.8 + SQLite ModernC + SDCC 4.6 + MSXgl\n"+
					"Licença: GNU General Public License v3 (GPL 3)\n"+
					"Autor: Barney / Wilson Pilon",
				version.String(),
			)
			dialog.ShowInformation("Sobre o Z-Realm", info, win)
		}),
		fyne.NewMenuItem("Especificações Técnicas MSX 2", func() {
			specs := "• Modo Gráfico: V9938 SCREEN 4 (Graphic 3 - 256x192)\n" +
				"• Viewport da Sala: 32x18 tiles (256x144 pixels)\n" +
				"• Sprites: Modo 2 (16x16 pixels multicolores por linha)\n" +
				"• Memória: MSX-DOS 2 Memory Mapper em páginas de 16 KB (Janela 0x8000-0xBFFF)\n" +
				"• Alvo Mínimo: MSX 2 com 128KB VRAM + 256KB Mapper"
			dialog.ShowInformation("Especificações Técnicas", specs, win)
		}),
	)

	return fyne.NewMainMenu(menuFile, menuTools, menuHelp)
}
