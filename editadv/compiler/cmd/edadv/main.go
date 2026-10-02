package main

import (
	"flag"
	"fmt"
	"os"

	"edadv/compiler/tui"
)

func main() {
	filePath := flag.String("f", "games/demo.yaml", "Arquivo de aventura (.yaml ou .json)")
	flag.Parse()

	// Se passou argumento posicional
	if flag.NArg() > 0 {
		*filePath = flag.Arg(0)
	}

	app := tui.NewAdventureEditor(*filePath)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na execução da TUI: %v\n", err)
		os.Exit(1)
	}
}
