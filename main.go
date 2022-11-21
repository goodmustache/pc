package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const BannerWidth = 70

func main() {
	prevCommandOutput, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		handleError(fmt.Errorf("PC ERROR - Unable to read from STDIN: %w", err))
	}

	printHeader()

	_, err = os.Stderr.Write(prevCommandOutput)
	if err != nil {
		handleError(fmt.Errorf("PC ERROR - Unable to write to STDERR: %w", err))
	}

	_, err = os.Stdout.Write(prevCommandOutput)
	if err != nil {
		handleError(fmt.Errorf("PC ERROR - Unable to write to STDOUT: %w", err))
	}
}

func printHeader() {
	output(strings.Repeat("=", BannerWidth) + "\n")
	output("PipeCheck: The following was read in and will be passed through\n")
	output(strings.Repeat("=", BannerWidth) + "\n")
}

func output(f string, opts ...any) {
	fmt.Fprintf(os.Stderr, f, opts...)
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
	os.Exit(1)
}
