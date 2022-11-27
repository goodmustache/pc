package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("PipeCheck", Label("integration"), func() {
	When("-h / --help flag is passed", func() {
		It("outputs the help text", func() {
			session := StartCommand(exec.Command(pcBinary, "-h"))

			Eventually(session).Should(Say("PipeCheck will output data recieved from STDIN to STDERR\\."))
		})
	})
})
