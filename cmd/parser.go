package cmd

import (
	"os"

	"github.com/goodmustache/pc/cmd/internal"
	"github.com/jessevdk/go-flags"
)

func ParseAndRun(args []string) error {
	process := &Process{
		StdIn:  os.Stdin,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
	}
	process.GetUserInput = process.GetUserInputImpl

	parser := flags.NewParser(process, flags.Default)
	parser.LongDescription = internal.LongDescription

	extraArgs, err := parser.ParseArgs(args[1:])
	if err != nil {
		if flags.WroteHelp(err) {
			return nil
		}
		return err
	}
	return process.Execute(extraArgs)
}
