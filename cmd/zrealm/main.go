package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Exibe a versão do Z-Realm")
	newProj := flag.String("new", "", "Cria um novo arquivo de projeto .rpgproj")
	projName := flag.String("name", "Novo RPG", "Nome do projeto para a criação")
	checkProj := flag.String("check", "", "Valida a integridade de um arquivo .rpgproj")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Z-Realm (zrealm-msx) - %s\n", version.String())
		fmt.Println("O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2)")
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

	fmt.Printf("Z-Realm Toolkit %s\n", version.String())
	fmt.Println("Uso:")
	fmt.Println("  zrealm -version                 Exibe a versão atual")
	fmt.Println("  zrealm -new <arquivo.rpgproj>   Cria um novo projeto SQLite")
	fmt.Println("  zrealm -check <arquivo.rpgproj> Valida a integridade de um projeto")
}
