package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/goodmustache/pc/cmd/internal"
	"github.com/goodmustache/pc/version"
)

// ErrExit is returned when the user denies the provided input.
var ErrExit = errors.New("user negatively confirmed input")

type Process struct {
	SkipValidation bool  `long:"skip-validation" description:"skips waiting for human validation" env:"PC_SKIP_VALIDATION"`
	BufferSize     int64 `long:"buffer-size" short:"b" default:"52428800" description:"size of the STDIN temporary buffer in bytes"`
	Version        bool  `long:"version" short:"v" description:"display version"`

	StdIn        io.Reader
	StdOut       io.Writer
	StdErr       io.Writer
	GetUserInput func(stream io.Writer) (byte, error)

	// OutputCompleteMessage is used for testing only. It outputs an additional
	// "complete" prior to exiting when set to true.
	OutputCompleteMessage bool
}

func (p Process) Execute(_ []string) error {
	// Output version an exit when version flag passed.
	if p.Version {
		return p.printVersion()
	}

	// Output Bannor
	fmt.Fprintln(p.StdErr, internal.Bannor)

	// Read STDIN into STDOUT and buffer
	inputBuffer := bytes.NewBuffer(make([]byte, 0, p.BufferSize))
	combinedWriter := io.MultiWriter(p.StdErr, inputBuffer)
	_, err := io.CopyBuffer(combinedWriter, p.StdIn, make([]byte, internal.CopyBufferSize))
	if err != nil {
		return err
	}

	if p.SkipValidation {
		// Output warning to STDOUT if skip is enabled
		fmt.Fprintln(p.StdErr, internal.SkipValidationMessage)
	} else {
		// Send validation to STDOUT and wait on user input
		err = p.validate()
		if err != nil {
			return err
		}
	}

	_, err = io.CopyBuffer(p.StdOut, inputBuffer, make([]byte, internal.CopyBufferSize))
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
	rawInput, err := p.GetUserInput(p.StdErr)
	if err != nil {
		return err
	}
	if input := rune(rawInput); input != 'y' && input != 'Y' {
		return ErrExit
	}
	return nil
}
