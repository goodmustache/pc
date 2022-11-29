package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/goodmustache/pc/cmd/internal"
	"github.com/goodmustache/pc/version"
)

// ErrExit is returned when the user denies the provided input.
var ErrExit = errors.New("user negatively confirmed input")

type Process struct {
	SkipValidation bool `long:"skip-validation" description:"skips waiting for human validation" env:"PC_SKIP_VALIDATION"`
	Version        bool `long:"version" short:"v" description:"display version"`

	StdIn        io.Reader
	StdOut       io.Writer
	StdErr       io.Writer
	GetUserInput func(stream io.Writer) (byte, error)

	// OutputCompleteMessage is used for testing only. It outputs an additional
	// "complete" prior to exiting when set to true.
	OutputCompleteMessage bool
}

func (p Process) Execute([]string) error {
	// Output version an exit when version flag passed.
	if p.Version {
		return p.printVersion()
	}

	// Read STDIN into buffer
	inputBuffer, err := io.ReadAll(p.StdIn)
	if err != nil {
		return err
	}

	// Output Bannor and buffer into STDERR
	fmt.Fprintln(p.StdErr, internal.Bannor)
	_, err = p.StdErr.Write(inputBuffer)
	if err != nil {
		return err
	}

	// Validation check
	err = p.validate()
	if err != nil {
		return err
	}

	// Output buffer to STDOUT
	_, err = p.StdOut.Write(inputBuffer)
	if err != nil {
		return err
	}

	if p.OutputCompleteMessage {
		fmt.Fprint(p.StdErr, "complete")
	}
	return nil
}

func (p Process) printVersion() error {
	fmt.Fprintf(p.StdOut, "PipeCheck %s", version.Version)
	return nil
}

func (p Process) validate() error {
	// Output warning to STDOUT if skip is enabled
	if p.SkipValidation {
		fmt.Fprintln(p.StdErr, internal.SkipValidationMessage)
		return nil
	}

	// Otherwise send validation to STDOUT and wait on user input
	rawInput, err := p.GetUserInput(p.StdErr)
	if err != nil {
		return err
	}
	if input := rune(rawInput); input != 'y' && input != 'Y' {
		return ErrExit
	}
	return nil
}
