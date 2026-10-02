package gui

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/theme"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

func TestRetroDarkTheme(t *testing.T) {
	th := NewRetroDarkTheme()

	bg := th.Color(theme.ColorNameBackground, theme.VariantDark)
	if bg == nil {
		t.Fatal("ColorNameBackground não deve ser nil")
	}

	primary := th.Color(theme.ColorNamePrimary, theme.VariantDark)
	if primary == nil {
		t.Fatal("ColorNamePrimary não deve ser nil")
	}
	// Confirma que Primary é o Cyan do MSX
	nrgba, ok := primary.(color.NRGBA)
	if !ok || nrgba.R != 0x00 || nrgba.G != 0xE5 || nrgba.B != 0xFF {
		t.Fatalf("Esperado Primary = MSX Cyan (0x00, 0xE5, 0xFF), obtido: %+v", primary)
	}

	pad := th.Size(theme.SizeNamePadding)
	if pad != 6 {
		t.Fatalf("Esperado padding 6, obtido: %v", pad)
	}
}

func TestProjectStateLifecycle(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zrealm_gui_test_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	state := NewProjectState()
	if state.IsOpen() {
		t.Fatal("estado inicial não deve ter projeto aberto")
	}

	var loadedCalled, closedCalled, modifiedCalled, dataChangedCalled bool
	state.OnProjectLoaded(func(_ *project.Project, _ string) {
		loadedCalled = true
	})
	state.OnProjectClosed(func() {
		closedCalled = true
	})
	state.OnModifiedChanged(func(m bool) {
		modifiedCalled = true
	})
	state.OnDataChanged(func() {
		dataChangedCalled = true
	})

	// 1. Criação de Novo Projeto
	projPath := filepath.Join(tmpDir, "test_proj.rpgproj")
	err = state.NewProject(projPath, "Projeto GUI Teste")
	if err != nil {
		t.Fatalf("falha ao criar novo projeto: %v", err)
	}

	if !state.IsOpen() {
		t.Fatal("esperado projeto aberto")
	}
	if !loadedCalled {
		t.Fatal("callback OnProjectLoaded não foi disparado")
	}
	if state.FilePath() != projPath {
		t.Fatalf("caminho retornado incorreto: %s", state.FilePath())
	}

	// 2. Modificações e Notificações
	state.SetModified(true)
	if !state.IsModified() {
		t.Fatal("esperado IsModified == true")
	}
	if !modifiedCalled {
		t.Fatal("callback OnModifiedChanged não disparou")
	}

	state.NotifyDataChanged()
	if !dataChangedCalled {
		t.Fatal("callback OnDataChanged não disparou")
	}

	// 3. Arquivos Recentes
	recents := state.RecentFiles()
	if len(recents) == 0 || recents[0] != projPath {
		t.Fatalf("esperado arquivo em recentes: %+v", recents)
	}

	// 4. Exportação
	outDir := filepath.Join(tmpDir, "build_msx")
	res, err := state.Export(outDir)
	if err != nil {
		t.Fatalf("falha ao exportar via state: %v", err)
	}
	if res.HeaderPath == "" || res.DataPath == "" {
		t.Fatalf("resultado de exportação incompleto: %+v", res)
	}

	// 5. Fechamento
	state.CloseProject()
	if state.IsOpen() {
		t.Fatal("esperado projeto fechado")
	}
	if !closedCalled {
		t.Fatal("callback OnProjectClosed não foi disparado")
	}
}

func TestProjectStateDemo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zrealm_gui_demo_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	demoPath := filepath.Join(tmpDir, "demo.rpgproj")
	state := NewProjectState()

	err = state.OpenDemoProject(demoPath)
	if err != nil {
		t.Fatalf("falha ao abrir demo via state: %v", err)
	}
	defer state.CloseProject()

	if !state.IsOpen() {
		t.Fatal("esperado projeto demo aberto")
	}

	// Exportação da demo
	outDir := filepath.Join(tmpDir, "demo_out")
	res, err := state.Export(outDir)
	if err != nil {
		t.Fatalf("falha ao exportar demo: %v", err)
	}
	if res.ResourceCount < 3 {
		t.Fatalf("esperado ao menos 3 recursos na demo, obtido: %d", res.ResourceCount)
	}
}
