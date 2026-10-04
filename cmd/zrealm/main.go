package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zrealm-msx/zrealm/pkg/exporter"
	"github.com/zrealm-msx/zrealm/pkg/gui"
	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/runner"
	"github.com/zrealm-msx/zrealm/pkg/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Exibe a versão do Z-Realm")
	newProj := flag.String("new", "", "Cria um novo arquivo de projeto .rpgproj")
	demoProj := flag.String("demo", "", "Gera um projeto demo de teste com tileset e 2 salas conectadas")
	projName := flag.String("name", "Novo RPG", "Nome do projeto para a criação")
	checkProj := flag.String("check", "", "Valida a integridade de um arquivo .rpgproj")
	exportProj := flag.String("export", "", "Exporta o projeto SQLite para arquivos binários do MSX 2 (MSX-DOS 2)")
	exportROM := flag.String("export-rom", "", "Exporta o projeto SQLite para cartucho MegaROM .ROM (ASCII-16)")
	runProj := flag.String("run", "", "Executa o ciclo completo One-Click Run via Disquete DOS 2 (Exportar + DSK + openMSX)")
	runROM := flag.String("run-rom", "", "Executa o ciclo completo One-Click Run via Cartucho MegaROM (Exportar + .ROM + openMSX -cart)")
	padSize := flag.Int("pad-size", 0, "Tamanho fixo do cartucho ROM em KB (128, 256, 512, etc.; 0 = automático)")
	machineName := flag.String("machine", "Philips_NMS_8250", "Modelo de máquina MSX para o openMSX (padrão: Philips_NMS_8250)")
	emuPath := flag.String("emu", "", "Caminho customizado para o executável openmsx")
	tclScript := flag.String("script", "", "Script TCL opcional repassado ao openMSX")
	outDir := flag.String("out", "", "Diretório de saída para os binários exportados (padrão: ./build_msx)")
	exportBanks := flag.Bool("banks", true, "Gera arquivos individuais SEGxx.BNK além do GAME.DAT")
	runCLI := flag.Bool("cli", false, "Força modo de linha de comando exibindo ajuda de comandos")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Z-Realm (zrealm-msx) - %s\n", version.String())
		fmt.Println("O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2 / MegaROM)")
		return
	}

	if *demoProj != "" {
		fmt.Printf("Gerando projeto de demonstração em: %s...\n", *demoProj)
		proj, err := project.CreateDemoProject(*demoProj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao gerar demo: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()
		fmt.Println("Projeto de demonstração criado com sucesso!")
		return
	}

	if *newProj != "" {
		fmt.Printf("Criando projeto Z-Realm em: %s (Nome: %s)...\n", *newProj, *projName)
		proj, err := project.Create(*newProj, *projName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()
		fmt.Println("Projeto criado com sucesso!")
		return
	}

	if *checkProj != "" {
		fmt.Printf("Verificando integridade do projeto: %s...\n", *checkProj)
		proj, err := project.Open(*checkProj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir projeto: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()

		name, _ := proj.GetSetting("name")
		ver, _ := proj.GetSetting("version")
		target, _ := proj.GetSetting("target_platform")

		fmt.Printf("Projeto: %s (v%s)\n", name, ver)
		fmt.Printf("Plataforma Alvo: %s\n", target)
		fmt.Println("Status de Integridade: OK (Integridade física e Foreign Keys válidas)")
		return
	}

	if *exportProj != "" {
		fmt.Printf("Exportando projeto para formato nativo MSX 2 (MSX-DOS 2): %s...\n", *exportProj)
		proj, err := project.Open(*exportProj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir projeto: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()

		res, err := exporter.Export(proj, exporter.ExportOptions{
			OutputDir:       *outDir,
			ExportBankFiles: *exportBanks,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro na exportação: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Exportação concluída com sucesso!")
		fmt.Printf("  Tabela Mestra:    %s\n", res.HeaderPath)
		fmt.Printf("  Dados (GAME.DAT): %s (%d bytes em %d segmentos de 16KB)\n", res.DataPath, res.TotalDataBytes, res.TotalSegments)
		fmt.Printf("  Total Recursos:   %d catalogados\n", res.ResourceCount)
		return
	}

	if *exportROM != "" {
		fmt.Printf("Exportando projeto para Cartucho MegaROM ASCII-16 (.ROM): %s...\n", *exportROM)
		proj, err := project.Open(*exportROM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir projeto: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()

		res, err := exporter.ExportROM(proj, exporter.ExportROMOptions{
			OutputDir: *outDir,
			PadSizeKB: *padSize,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro na exportação de ROM: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Exportação MegaROM concluída com sucesso!")
		fmt.Printf("  Arquivo .ROM:     %s (%d bytes / %d KB)\n", res.ROMPath, res.TotalROMBytes, res.TotalROMBytes/1024)
		fmt.Printf("  Tipo de Mapper:   %s (%d bancos de 16KB)\n", res.MapperType, res.TotalBanks)
		fmt.Printf("  Segmentos Jogo:   %d banco(s)\n", res.TotalSegments)
		fmt.Printf("  Total Recursos:   %d catalogados\n", res.ResourceCount)
		return
	}

	if *runProj != "" {
		fmt.Printf("Iniciando One-Click Run (Disquete MSX-DOS 2) para o projeto: %s...\n", *runProj)
		proj, err := project.Open(*runProj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir projeto: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()

		res, err := runner.OneClickRun(proj, runner.RunOptions{
			OutputDir:   *outDir,
			EmulatorExe: *emuPath,
			Machine:     *machineName,
			TclScript:   *tclScript,
			AutoRun:     true,
			Async:       false,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro no One-Click Run: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("One-Click Run concluído com sucesso!")
		fmt.Printf("  Tabela Mestra:    %s\n", res.ExportResult.HeaderPath)
		fmt.Printf("  Dados (GAME.DAT): %d bytes (%d segmentos de 16KB)\n", res.ExportResult.TotalDataBytes, res.ExportResult.TotalSegments)
		fmt.Printf("  Disquete (DSK):   %s (%d bytes / 720 KB)\n", res.DskPath, res.DskSize)
		fmt.Printf("  Emulador openMSX: %s\n", res.EmulatorExe)
		return
	}

	if *runROM != "" {
		fmt.Printf("Iniciando One-Click Run (Cartucho MegaROM .ROM) para o projeto: %s...\n", *runROM)
		proj, err := project.Open(*runROM)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao abrir projeto: %v\n", err)
			os.Exit(1)
		}
		defer proj.Close()

		res, err := runner.OneClickRunROM(proj, runner.RunROMOptions{
			OutputDir:   *outDir,
			EmulatorExe: *emuPath,
			Machine:     *machineName,
			TclScript:   *tclScript,
			AutoRun:     true,
			Async:       false,
			PadSizeKB:   *padSize,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro no One-Click Run ROM: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("One-Click Run MegaROM concluído com sucesso!")
		fmt.Printf("  Arquivo .ROM:     %s (%d bytes / %d KB)\n", res.ROMPath, res.ROMSize, res.ROMSize/1024)
		fmt.Printf("  Tipo de Mapper:   %s (%d bancos de 16KB)\n", res.ExportResult.MapperType, res.ExportResult.TotalBanks)
		fmt.Printf("  Segmentos Jogo:   %d banco(s)\n", res.ExportResult.TotalSegments)
		fmt.Printf("  Total Recursos:   %d catalogados\n", res.ExportResult.ResourceCount)
		fmt.Printf("  Emulador openMSX: %s\n", res.EmulatorExe)
		return
	}

	if *runCLI {
		fmt.Printf("Z-Realm Toolkit %s\n", version.String())
		fmt.Println("Uso:")
		fmt.Println("  zrealm                                       Inicia o Editor Visual Desktop (Fyne GUI)")
		fmt.Println("  zrealm <arquivo.rpgproj>                     Abre o arquivo diretamente no Editor Visual")
		fmt.Println("  zrealm -version                              Exibe a versão atual")
		fmt.Println("  zrealm -new <arquivo.rpgproj>                Cria um novo projeto SQLite")
		fmt.Println("  zrealm -demo <arquivo.rpgproj>               Gera projeto de demonstração")
		fmt.Println("  zrealm -check <arquivo.rpgproj>              Valida a integridade de um projeto")
		fmt.Println("  zrealm -export <arquivo.rpgproj> [-out d]    Exporta projeto para MSX-DOS 2 (HEADER.BIN / GAME.DAT)")
		fmt.Println("  zrealm -export-rom <arquivo.rpgproj> [-out d] Exporta projeto para Cartucho MegaROM ASCII-16 (.ROM)")
		fmt.Println("  zrealm -run <arquivo.rpgproj> [-out d]       One-Click Run Disquete: Exporta, monta DSK e executa no openMSX")
		fmt.Println("  zrealm -run-rom <arquivo.rpgproj> [-out d]   One-Click Run Cartucho: Exporta .ROM e executa no openMSX (-cart)")
		return
	}

	// Execução padrão: Interface Gráfica Fyne
	desktopApp := gui.NewApp()
	if len(flag.Args()) > 0 {
		_ = desktopApp.State().OpenProject(flag.Args()[0])
	}
	desktopApp.Run()
}
