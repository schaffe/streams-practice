package main

import (
	"flag"
	"fmt"
	"os"

	"streams-practice/assignments/rover/entity"
	"streams-practice/assignments/rover/ui"
	"streams-practice/assignments/rover/world"
)

func main() {
	cmd := flag.String("cmd", "", "batch command string (L/R/M); if set, run non-interactively")
	flag.Parse()

	rover := entity.NewRover(ui.ActiveUnit, 0, 0, entity.North)
	w := world.New(rover)
	cli := ui.NewCLI(w, ui.NewLineRenderer(), os.Stdin, os.Stdout)

	if *cmd != "" {
		if err := cli.RunBatch(*cmd); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := cli.RunInteractive(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
