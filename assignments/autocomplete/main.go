package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"streams-practice/assignments/autocomplete/cli"
	"streams-practice/assignments/autocomplete/corpus"
	"streams-practice/assignments/autocomplete/index"
)

func dataPath() string {
	_, src, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(src), "corpus", "frequencies.json")
}

func main() {
	prefix := flag.String("prefix", "", "search prefix (batch mode)")
	flag.Parse()

	entries, err := corpus.Load(dataPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	t := index.New(entries)
	c := cli.NewCLI(t, os.Stdin, os.Stdout)

	if *prefix != "" || flag.NFlag() > 0 {
		if err := c.RunBatch(*prefix); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := c.RunInteractive(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
