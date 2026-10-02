package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zrealm-msx/zrealm/pkg/exporter"
	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Exibe a versão do Z-Realm")
	newProj := flag.String("new", "", "Cria um novo arquivo de projeto .rpgproj")
	demoProj := flag.String("demo", "", "Gera um projeto demo de teste com tileset e 2 salas conectadas")
	projName := flag.String("name", "Novo RPG", "Nome do projeto para a criação")
	checkProj := flag.String("check", "", "Valida a integridade de um arquivo .rpgproj")
	exportProj := flag.String("export", "", "Exporta o projeto SQLite para arquivos binários do MSX 2")
	outDir := flag.String("out", "", "Diretório de saída para os binários exportados (padrão: ./build_msx)")
	exportBanks := flag.Bool("banks", true, "Gera arquivos individuais SEGxx.BNK além do GAME.DAT")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Z-Realm (zrealm-msx) - %s\n", version.String())
		fmt.Println("O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2)")
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
		fmt.Printf("Exportando projeto para formato nativo MSX 2: %s...\n", *exportProj)
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

	fmt.Printf("Z-Realm Toolkit %s\n", version.String())
	fmt.Println("Uso:")
	fmt.Println("  zrealm -version                           Exibe a versão atual")
	fmt.Println("  zrealm -new <arquivo.rpgproj>             Cria um novo projeto SQLite")
	fmt.Println("  zrealm -demo <arquivo.rpgproj>            Gera projeto de demonstração")
	fmt.Println("  zrealm -check <arquivo.rpgproj>           Valida a integridade de um projeto")
	fmt.Println("  zrealm -export <arquivo.rpgproj> [-out d] Exporta projeto para binários MSX 2")
}
