package cmd_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

// Log will write out to the Ginkgo logs. These only appear when a test fails
// or running in verbose mode.
func Log(f string, opts ...any) {
	fmt.Fprintf(GinkgoWriter, f+"\n", opts...)
}
