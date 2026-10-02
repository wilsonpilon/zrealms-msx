package main

import (
	"flag"
	"fmt"
	"os"

	"edadv/compiler"
)

func main() {
	inputFile := flag.String("i", "games/demo.yaml", "Caminho do arquivo de história (.yaml ou .json)")
	outputDir := flag.String("o", "engine/src", "Diretório de saída para os arquivos .c e .h")
	flag.Parse()

	fmt.Printf("[Editor de Adventures Compiler] Compilando %s -> %s\n", *inputFile, *outputDir)

	comp := compiler.NewCompiler()
	if err := comp.LoadFile(*inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar história: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Jogo: \"%s\" por %s (v%s)\n", comp.Game.Meta.Title, comp.Game.Meta.Author, comp.Game.Meta.Version)
	fmt.Printf("Estatísticas: %d posições, %d objetos, %d comandos, %d funções, %d mensagens\n",
		len(comp.Game.Positions), len(comp.Game.Objects), len(comp.Game.Commands),
		len(comp.Game.Functions), len(comp.Game.Messages))

	if err := comp.GenerateCData(*outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gerar código C: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sucesso! Arquivos gerados:\n  - %s/game_data.h\n  - %s/game_data.c\n", *outputDir, *outputDir)
}
