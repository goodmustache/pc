package cmd

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/goodmustache/pc/cmd/internal"
	"github.com/jessevdk/go-flags"
)

func ParseAndRun(args []string) error {
	process := &Process{
		StdIn:  os.Stdin,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		GetUserInput: func(stream io.Writer) (byte, error) {
			in, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
			if err != nil {
				return 0, fmt.Errorf("unable to reopen TTY: %w", err)
			}

			fmt.Fprint(stream, internal.ValidationMessage)
			input := make([]byte, 1)
			_, err = syscall.Read(in, input)
			if err != nil {
				return 0, fmt.Errorf("unable to read from user STDIN: %w", err)
			}
			fmt.Fprintln(stream, internal.ValidationMessageBannor)

			return input[0], nil
		},
	}
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
