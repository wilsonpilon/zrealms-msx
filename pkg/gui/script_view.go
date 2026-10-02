package gui

import (
	"context"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

// ScriptView implementa a visualização e edição de scripts de eventos (Bytecode VM) e diálogos.
type ScriptView struct {
	state          *ProjectState
	win            fyne.Window

	// Scripts de Eventos
	scriptList     *widget.List
	scripts        []*models.Script
	selectedScript *models.Script
	entryName      *widget.Entry
	entrySource    *widget.Entry
	lblBytecode    *widget.Label
	lblCompileInfo *widget.Label

	// Diálogos e Textos
	stringList     *widget.List
	strings        []*models.StringEntry
	selectedString *models.StringEntry
	lblStringInfo  *widget.Label
	lblPreviewBox  *widget.Label

	lblStatus      *widget.Label
}

// NewScriptView cria a aba de roteiros, eventos e diálogos do jogo.
func NewScriptView(state *ProjectState, win fyne.Window) fyne.CanvasObject {
	v := &ScriptView{
		state:          state,
		win:            win,
		scripts:        make([]*models.Script, 0),
		strings:        make([]*models.StringEntry, 0),
		entryName:      widget.NewEntry(),
		entrySource:    widget.NewMultiLineEntry(),
		lblBytecode:    widget.NewLabel("(Nenhum bytecode gerado)"),
		lblCompileInfo: widget.NewLabel("Escreva o script e clique em Compilar."),
		lblStringInfo:  widget.NewLabel("Selecione um diálogo para editar."),
		lblPreviewBox:  widget.NewLabel(""),
		lblStatus:      widget.NewLabel("Pronto."),
	}

	// =========================================================
	// 1. ABA DE SCRIPTS DE EVENTOS (BYTECODE VM)
	// =========================================================
	v.scriptList = widget.NewList(
		func() int { return len(v.scripts) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.MediaPlayIcon()), widget.NewLabel("Script"))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.scripts) {
				return
			}
			s := v.scripts[id]
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText(
				fmt.Sprintf("#%d: %s (%d bytes)", s.ID, s.Name, len(s.Bytecode)),
			)
		},
	)

	v.scriptList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.scripts) {
			return
		}
		v.selectedScript = v.scripts[id]
		v.loadScriptDetails()
	}

	btnAddScript := widget.NewButtonWithIcon("Novo Script", theme.ContentAddIcon(), func() {
		v.createNewScript()
	})
	btnDeleteScript := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedScript()
	})

	scriptListHeader := container.NewVBox(
		widget.NewLabelWithStyle("⚡ Scripts de Eventos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAddScript, btnDeleteScript),
		widget.NewSeparator(),
	)
	scriptSidebar := container.NewBorder(scriptListHeader, nil, nil, nil, v.scriptList)

	// Painel de Edição de Código do Script
	v.entryName.SetPlaceHolder("Nome do Script (ex: bau_tesouro)")
	v.entrySource.SetPlaceHolder("// Digite aqui as instruções da VM:\nMSG 1\nGIVE_ITEM 2\nEND")
	v.entrySource.Wrapping = fyne.TextWrapWord

	btnCompile := widget.NewButtonWithIcon("Compilar Bytecode", theme.MediaPlayIcon(), func() {
		v.compileCurrentScript()
	})
	btnSaveScript := widget.NewButtonWithIcon("Salvar Script no SQLite", theme.DocumentSaveIcon(), func() {
		v.saveCurrentScript()
	})
	btnSaveScript.Importance = widget.HighImportance

	// Botões de Inserção Rápida de Snippets
	btnSnippetMsg := widget.NewButton("Inserir MSG", func() { v.insertSnippet("MSG 1\n") })
	btnSnippetItem := widget.NewButton("Inserir GIVE_ITEM", func() { v.insertSnippet("GIVE_ITEM 1\n") })
	btnSnippetFlag := widget.NewButton("Inserir CHECK_FLAG", func() { v.insertSnippet("CHECK_FLAG 1 label_fim\n") })
	btnSnippetTele := widget.NewButton("Inserir TELEPORT", func() { v.insertSnippet("TELEPORT 1 16 9\n") })

	snippetsBar := container.NewHBox(
		widget.NewLabel("Snippets:"),
		btnSnippetMsg, btnSnippetItem, btnSnippetFlag, btnSnippetTele,
	)

	scriptEditorPanel := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Nome:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewGridWrap(fyne.NewSize(220, 36), v.entryName),
			btnCompile,
			btnSaveScript,
		),
		snippetsBar,
		widget.NewSeparator(),
		container.NewGridWrap(fyne.NewSize(500, 220), v.entrySource),
		widget.NewSeparator(),
		widget.NewCard("Bytecode VM Compilado (Hexadecimal)", "Formato binário para execução direta no Z80 do MSX 2", container.NewVBox(
			v.lblCompileInfo,
			v.lblBytecode,
		)),
	)

	scriptSplit := container.NewHSplit(scriptSidebar, container.NewPadded(scriptEditorPanel))
	scriptSplit.Offset = 0.28

	// =========================================================
	// 2. ABA DE TABELA DE DIÁLOGOS & TEXTOS DO HUD
	// =========================================================
	v.stringList = widget.NewList(
		func() int { return len(v.strings) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Texto"))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(v.strings) {
				return
			}
			s := v.strings[id]
			obj.(*fyne.Container).Objects[1].(*widget.Label).SetText(
				fmt.Sprintf("[%s] #%d: %s", s.ContextTag, s.ID, truncate(s.TextContent, 28)),
			)
		},
	)

	v.stringList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(v.strings) {
			return
		}
		v.selectedString = v.strings[id]
		v.updateStringDetails()
	}

	btnAddString := widget.NewButtonWithIcon("Novo Diálogo", theme.ContentAddIcon(), func() {
		v.showStringForm(nil)
	})
	btnEditString := widget.NewButtonWithIcon("Editar", theme.DocumentCreateIcon(), func() {
		if v.selectedString != nil {
			v.showStringForm(v.selectedString)
		} else {
			dialog.ShowInformation("Aviso", "Selecione um texto para editar.", v.win)
		}
	})
	btnDeleteString := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), func() {
		v.deleteSelectedString()
	})

	stringListHeader := container.NewVBox(
		widget.NewLabelWithStyle("📜 Textos & Diálogos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(btnAddString, btnEditString, btnDeleteString),
		widget.NewSeparator(),
	)
	stringSidebar := container.NewBorder(stringListHeader, nil, nil, nil, v.stringList)

	cardDialoguePreview := widget.NewCard(
		"Simulador de Caixa de Diálogo (MSX SCREEN 4: 32x4 caracteres)",
		"Área das linhas 20 a 23 da tela no chip V9938",
		container.NewVBox(
			v.lblStringInfo,
			widget.NewSeparator(),
			v.lblPreviewBox,
		),
	)

	stringSplit := container.NewHSplit(stringSidebar, container.NewPadded(cardDialoguePreview))
	stringSplit.Offset = 0.32

	// Tabs principais: Scripts e Diálogos
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Scripts de Eventos (VM)", theme.MediaPlayIcon(), scriptSplit),
		container.NewTabItemWithIcon("Tabela de Diálogos", theme.DocumentIcon(), stringSplit),
	)

	state.OnProjectLoaded(func(_ *project.Project, _ string) { v.reload() })
	state.OnProjectClosed(func() {
		v.scripts = make([]*models.Script, 0)
		v.strings = make([]*models.StringEntry, 0)
		v.selectedScript = nil
		v.selectedString = nil
		v.scriptList.Refresh()
		v.stringList.Refresh()
		v.entryName.SetText("")
		v.entrySource.SetText("")
		v.lblBytecode.SetText("(Nenhum projeto aberto)")
		v.lblCompileInfo.SetText("")
		v.lblStringInfo.SetText("Nenhum projeto aberto.")
		v.lblPreviewBox.SetText("")
	})

	if state.IsOpen() {
		v.reload()
	}

	return container.NewBorder(nil, v.lblStatus, nil, nil, tabs)
}

func (v *ScriptView) reload() {
	proj := v.state.Current()
	if proj == nil {
		return
	}

	// Carrega Scripts
	scripts, errS := proj.Storage().GameData.ListScripts(context.Background())
	if errS == nil {
		v.scripts = scripts
		v.scriptList.Refresh()
		if len(v.scripts) > 0 && v.selectedScript == nil {
			v.scriptList.Select(0)
		}
	}

	// Carrega Diálogos
	strs, errD := proj.Storage().GameData.ListStrings(context.Background())
	if errD == nil {
		v.strings = strs
		v.stringList.Refresh()
		if len(v.strings) > 0 && v.selectedString == nil {
			v.stringList.Select(0)
		}
	}
}

func (v *ScriptView) loadScriptDetails() {
	if v.selectedScript == nil {
		return
	}
	s := v.selectedScript
	v.entryName.SetText(s.Name)
	v.entrySource.SetText(s.SourceCode)

	if len(s.Bytecode) > 0 {
		v.lblBytecode.SetText(FormatBytecodeHex(s.Bytecode))
		v.lblCompileInfo.SetText(fmt.Sprintf("Bytecode Válido: %d bytes alocados.", len(s.Bytecode)))
	} else {
		v.lblBytecode.SetText("(Bytecode não compilado)")
		v.lblCompileInfo.SetText("Clique em 'Compilar Bytecode' para montar.")
	}
}

func (v *ScriptView) insertSnippet(snippet string) {
	cur := v.entrySource.Text
	if cur != "" && !strings.HasSuffix(cur, "\n") {
		cur += "\n"
	}
	v.entrySource.SetText(cur + snippet)
}

func (v *ScriptView) compileCurrentScript() ([]byte, error) {
	src := v.entrySource.Text
	bytecode, err := CompileScript(src)
	if err != nil {
		v.lblCompileInfo.SetText(fmt.Sprintf("❌ Erro de Compilação: %s", err.Error()))
		v.lblBytecode.SetText("(Falha no build)")
		return nil, err
	}

	v.lblCompileInfo.SetText(fmt.Sprintf("✅ Sucesso! Tamanho: %d bytes.", len(bytecode)))
	v.lblBytecode.SetText(FormatBytecodeHex(bytecode))
	if v.selectedScript != nil {
		v.selectedScript.Bytecode = bytecode
		v.selectedScript.SourceCode = src
	}
	return bytecode, nil
}

func (v *ScriptView) saveCurrentScript() {
	if v.selectedScript == nil {
		dialog.ShowInformation("Aviso", "Selecione um script para salvar.", v.win)
		return
	}
	proj := v.state.Current()
	if proj == nil {
		return
	}

	bytecode, err := v.compileCurrentScript()
	if err != nil {
		dialog.ShowError(fmt.Errorf("não é possível salvar script com erro de sintaxe: %w", err), v.win)
		return
	}

	v.selectedScript.Name = v.entryName.Text
	v.selectedScript.SourceCode = v.entrySource.Text
	v.selectedScript.Bytecode = bytecode

	if err := proj.Storage().GameData.UpdateScript(context.Background(), v.selectedScript); err != nil {
		dialog.ShowError(err, v.win)
		return
	}

	v.state.NotifyDataChanged()
	v.scriptList.Refresh()
	v.lblStatus.SetText(fmt.Sprintf("Script '%s' salvo com sucesso no banco SQLite!", v.selectedScript.Name))
}

func (v *ScriptView) createNewScript() {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("script_%d", len(v.scripts)+1))

	items := []*widget.FormItem{
		widget.NewFormItem("Nome do Script", nameEntry),
	}

	dialog.ShowForm("Novo Script de Evento", "Criar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}

		defaultSrc := "// Evento de Interação\nMSG 1\nEND\n"
		bc, _ := CompileScript(defaultSrc)

		s := &models.Script{
			Name:       nameEntry.Text,
			SourceCode: defaultSrc,
			Bytecode:   bc,
		}

		if err := proj.Storage().GameData.CreateScript(context.Background(), s); err != nil {
			dialog.ShowError(err, v.win)
			return
		}

		v.state.NotifyDataChanged()
		v.reload()
		for i, sc := range v.scripts {
			if sc.ID == s.ID {
				v.scriptList.Select(i)
				break
			}
		}
	}, v.win)
}

func (v *ScriptView) deleteSelectedScript() {
	if v.selectedScript == nil {
		dialog.ShowInformation("Aviso", "Selecione um script para excluir.", v.win)
		return
	}

	proj := v.state.Current()
	if proj == nil {
		return
	}

	s := v.selectedScript
	dialog.ShowConfirm("Excluir Script", fmt.Sprintf("Deseja realmente excluir o script '%s' (#%d)?", s.Name, s.ID), func(ok bool) {
		if !ok {
			return
		}
		if err := proj.Storage().GameData.DeleteScript(context.Background(), s.ID); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.selectedScript = nil
		v.state.NotifyDataChanged()
		v.reload()
	}, v.win)
}

func (v *ScriptView) updateStringDetails() {
	if v.selectedString == nil {
		v.lblStringInfo.SetText("Nenhum texto selecionado.")
		v.lblPreviewBox.SetText("")
		return
	}
	s := v.selectedString
	info := fmt.Sprintf(
		"ID: #%d\nTag de Contexto: %s\nTamanho: %d caracteres (%d bytes)",
		s.ID, s.ContextTag, len([]rune(s.TextContent)), len(s.TextContent),
	)
	v.lblStringInfo.SetText(info)

	// Simula a moldura da caixa de diálogo do MSX (32 colunas)
	var sb strings.Builder
	border := "+--------------------------------+"
	sb.WriteString(border + "\n")
	lines := wrapText(s.TextContent, 30)
	for i := 0; i < 4; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		padding := 30 - len([]rune(line))
		if padding < 0 {
			padding = 0
		}
		sb.WriteString("| " + line + strings.Repeat(" ", padding) + " |\n")
	}
	sb.WriteString(border)
	v.lblPreviewBox.SetText(sb.String())
}

func (v *ScriptView) showStringForm(s *models.StringEntry) {
	proj := v.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Abra ou crie um projeto primeiro.", v.win)
		return
	}

	isNew := s == nil
	if isNew {
		s = &models.StringEntry{
			ContextTag:  "dialog_novo",
			TextContent: "Olá, forasteiro!",
		}
	}

	tagEntry := widget.NewEntry()
	tagEntry.SetText(s.ContextTag)

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetText(s.TextContent)
	textEntry.Wrapping = fyne.TextWrapWord

	items := []*widget.FormItem{
		widget.NewFormItem("Tag de Contexto", tagEntry),
		widget.NewFormItem("Conteúdo do Diálogo", textEntry),
	}

	title := "Novo Diálogo"
	if !isNew {
		title = fmt.Sprintf("Editar Diálogo #%d", s.ID)
	}

	dialog.ShowForm(title, "Salvar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		s.ContextTag = tagEntry.Text
		s.TextContent = textEntry.Text

		if isNew {
			if err := proj.Storage().GameData.CreateString(context.Background(), s); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText("Diálogo criado com sucesso.")
		} else {
			if err := proj.Storage().GameData.UpdateString(context.Background(), s); err != nil {
				dialog.ShowError(err, v.win)
				return
			}
			v.lblStatus.SetText("Diálogo atualizado com sucesso.")
		}

		v.state.NotifyDataChanged()
		v.reload()
		v.updateStringDetails()
	}, v.win)
}

func (v *ScriptView) deleteSelectedString() {
	if v.selectedString == nil {
		dialog.ShowInformation("Aviso", "Selecione um diálogo para excluir.", v.win)
		return
	}

	proj := v.state.Current()
	if proj == nil {
		return
	}

	s := v.selectedString
	dialog.ShowConfirm("Excluir Diálogo", fmt.Sprintf("Deseja realmente excluir o diálogo #%d ('%s')?", s.ID, s.ContextTag), func(ok bool) {
		if !ok {
			return
		}
		if err := proj.Storage().GameData.DeleteString(context.Background(), s.ID); err != nil {
			dialog.ShowError(err, v.win)
			return
		}
		v.selectedString = nil
		v.state.NotifyDataChanged()
		v.reload()
	}, v.win)
}

func wrapText(text string, maxLen int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	curLine := ""

	for _, w := range words {
		if len(curLine)+len(w)+1 <= maxLen {
			if curLine == "" {
				curLine = w
			} else {
				curLine += " " + w
			}
		} else {
			lines = append(lines, curLine)
			curLine = w
		}
	}
	if curLine != "" {
		lines = append(lines, curLine)
	}
	return lines
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
