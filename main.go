package main

import (
	"fmt"
	"os"

	"github.com/goodmustache/pc/cmd"
)

func main() {
	err := cmd.ParseAndRun(os.Args)
	if err != nil {
		if err != cmd.ErrExit {
			fmt.Fprintf(os.Stderr, "PipeCheck Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}
