package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zrealm-msx/zrealm/pkg/exporter"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

func TestExtractSystemFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_extract_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := ExtractSystemFiles(tmpDir); err != nil {
		t.Fatalf("ExtractSystemFiles falhou: %v", err)
	}

	expectedFiles := []string{"autoexec.bat", "COMMAND2.COM", "MSXDOS2.SYS", "zrealm.com"}
	for _, f := range expectedFiles {
		p := filepath.Join(tmpDir, f)
		fi, err := os.Stat(p)
		if err != nil {
			t.Errorf("arquivo esperado %s não existe: %v", f, err)
			continue
		}
		if fi.Size() == 0 {
			t.Errorf("arquivo %s está vazio (0 bytes)", f)
		}
	}
}

func TestFindDskTool(t *testing.T) {
	toolPath, err := FindDskTool("")
	if err != nil {
		t.Fatalf("FindDskTool falhou: %v", err)
	}
	if toolPath == "" {
		t.Fatal("caminho de msxtar retornado está vazio")
	}
	if _, err := os.Stat(toolPath); err != nil {
		t.Fatalf("ferramenta dsk encontrada não existe fisicamente em %s: %v", toolPath, err)
	}
}

func TestFindOpenMSX(t *testing.T) {
	emuPath, err := FindOpenMSX("")
	if err != nil {
		t.Fatalf("FindOpenMSX falhou: %v", err)
	}
	if emuPath == "" {
		t.Fatal("caminho de openmsx retornado está vazio")
	}
	if _, err := os.Stat(emuPath); err != nil {
		t.Fatalf("openmsx encontrado não existe fisicamente em %s: %v", emuPath, err)
	}
}

func TestBuildOpenMSXArgs(t *testing.T) {
	dsk := "c:/caminho/teste.dsk"

	// Caso 1: Padrão Philips NMS 8250
	args1 := BuildOpenMSXArgs(dsk, RunOptions{
		Machine: "Philips_NMS_8250",
	})
	expected1 := []string{"-machine", "Philips_NMS_8250", "-ext", "msxdos2", "-ext", "ram512k", "-diska", filepath.Clean(dsk)}
	if len(args1) != len(expected1) {
		t.Errorf("tamanho inesperado de args1: obtido %d, esperado %d (%v)", len(args1), len(expected1), args1)
	}

	// Caso 2: Máquina Panasonic com script TCL
	args2 := BuildOpenMSXArgs(dsk, RunOptions{
		Machine:   "Panasonic_FS-A1GT",
		ExtraArgs: []string{},
		TclScript: "c:/scripts/teste.tcl",
	})
	foundDiska := false
	foundScript := false
	for i, a := range args2 {
		if a == "-diska" && i+1 < len(args2) {
			foundDiska = true
		}
		if a == "-script" && i+1 < len(args2) {
			foundScript = true
		}
	}
	if !foundDiska {
		t.Errorf("argumento -diska ausente em args2: %v", args2)
	}
	if !foundScript {
		t.Errorf("argumento -script ausente em args2: %v", args2)
	}
}

func TestPackDskAndVerify(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_pack_dsk_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Cria projeto demo temporário
	projPath := filepath.Join(tmpDir, "demo.rpgproj")
	proj, err := project.CreateDemoProject(projPath)
	if err != nil {
		t.Fatalf("erro ao criar projeto demo: %v", err)
	}
	defer proj.Close()

	// Exporta os dados
	expOutDir := filepath.Join(tmpDir, "build_msx")
	expRes, err := exporter.Export(proj, exporter.ExportOptions{
		OutputDir:       expOutDir,
		ExportBankFiles: true,
	})
	if err != nil {
		t.Fatalf("erro ao exportar dados: %v", err)
	}

	// Empacota o DSK
	targetDsk := filepath.Join(tmpDir, "meu_jogo.dsk")
	createdDsk, err := PackDsk(PackOptions{
		DskPath:       targetDsk,
		HeaderBinPath: expRes.HeaderPath,
		GameDatPath:   expRes.DataPath,
	})
	if err != nil {
		t.Fatalf("PackDsk falhou: %v", err)
	}

	// Valida tamanho
	fi, err := os.Stat(createdDsk)
	if err != nil {
		t.Fatalf("não foi possível ler o arquivo DSK criado: %v", err)
	}
	const expectedSize = 737280
	if fi.Size() != expectedSize {
		t.Fatalf("tamanho do DSK incorreto: obtido %d, esperado %d", fi.Size(), expectedSize)
	}

	// Inspeciona conteúdo com VerifyDskContents
	entries, err := VerifyDskContents(createdDsk, "")
	if err != nil {
		t.Fatalf("VerifyDskContents falhou: %v", err)
	}

	contentStr := strings.ToUpper(strings.Join(entries, "\n"))
	expectedNames := []string{"AUTOEXEC.BAT", "ZREALM.COM", "COMMAND2.COM", "MSXDOS2.SYS", "HEADER.BIN", "GAME.DAT"}
	for _, en := range expectedNames {
		if !strings.Contains(contentStr, en) {
			t.Errorf("arquivo %s não encontrado no disquete DSK:\n%s", en, contentStr)
		}
	}
}

func TestOneClickRunPipeline(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_oneclick_*")
	if err != nil {
		t.Fatalf("falha ao criar pasta temporária: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projPath := filepath.Join(tmpDir, "demo.rpgproj")
	proj, err := project.CreateDemoProject(projPath)
	if err != nil {
		t.Fatalf("erro ao criar projeto demo: %v", err)
	}
	defer proj.Close()

	// Testa execução com AutoRun=false (apenas pipeline de build/export/dsk)
	res, err := OneClickRun(proj, RunOptions{
		OutputDir: filepath.Join(tmpDir, "msx_out"),
		AutoRun:   false,
	})
	if err != nil {
		t.Fatalf("OneClickRun falhou: %v", err)
	}

	if res.ExportResult == nil {
		t.Error("ExportResult ausente no resultado do OneClickRun")
	}
	if res.DskPath == "" {
		t.Error("DskPath vazio no resultado do OneClickRun")
	}
	if res.DskSize != 737280 {
		t.Errorf("DskSize incorreto: obtido %d, esperado 737280", res.DskSize)
	}

	if _, err := os.Stat(res.DskPath); err != nil {
		t.Errorf("arquivo DSK não encontrado no caminho indicado: %v", err)
	}
}
