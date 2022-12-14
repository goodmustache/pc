package integration_test

import (
	"os/exec"
	"time"

	"gopkg.in/alessio/shellescape.v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("PipeCheck", Label("integration"), func() {
	When("version flag is passed", func() {
		When("-v flag is passed", func() {
			It("outputs the version", func() {
				session := StartCommand(exec.Command(pcBinary, "-v"))

				Eventually(session).Should(Say(`PipeCheck v\d+\.\d+\.\d+(-pre)?`))
			})
		})

		When("--version flag is passed", func() {
			It("outputs the version", func() {
				session := StartCommand(exec.Command(pcBinary, "--version"))

				Eventually(session).Should(Say(`PipeCheck v\d+\.\d+\.\d+(-pre)?`))
			})
		})
	})

	When("help flag is passed", func() {
		When("-h flag is passed", func() {
			It("outputs the help text", func() {
				session := StartCommand(exec.Command(pcBinary, "-h"))

				Eventually(session).Should(Say("PipeCheck will output data received from STDIN to STDERR\\."))
			})
		})

		When("--help flag is passed", func() {
			It("outputs the help text", func() {
				session := StartCommand(exec.Command(pcBinary, "--help"))

				Eventually(session).Should(Say("PipeCheck will output data received from STDIN to STDERR\\."))
			})
		})
	})

	When("passed a piped command as input", func() {
		It("outputs and blocks successfully", Label("tty-required"), func() {
			cmd := CreateCommand(`echo %s | {{.pcBinary}} | xargs -n 1 echo -`, shellescape.Quote("foo\nbar\nbaz"))
			session := StartCommand(cmd)
			defer session.Kill() // shutdown the process

			Eventually(session.Err).Should(Say("PipeCheck: The following was read in and will be passed through:"))

			Eventually(session.Err).Should(Say(`=\nfoo\n`))
			Eventually(session.Err).Should(Say(`bar\n`))
			Eventually(session.Err).Should(Say(`baz\n=`))

			Eventually(session.Err).Should(Say(`Proceed with this data \(y/N\):`))

			// validate blocking
			Consistently(session).ShouldNot(Say(`- foo`))
			// Since we're unable to test the input, of "y" end test at blocking
		})

		It("terminates when passed a 'ctrl+c' at input", Pending, func() {
			cmd := CreateCommand(`echo %s | {{.pcBinary}} | xargs -n 1 echo -`, shellescape.Quote("foo\nbar\nbaz"))
			session := StartCommand(cmd)

			Eventually(session.Err).Should(Say("PipeCheck: The following was read in and will be passed through:"))

			Eventually(session.Err).Should(Say(`=\nfoo\n`))
			Eventually(session.Err).Should(Say(`bar\n`))
			Eventually(session.Err).Should(Say(`baz\n=`))

			Eventually(session.Err).Should(Say(`Proceed with this data \(y/N\):`))
			session.Interrupt()

			Consistently(session).ShouldNot(Say(`- foo`))
			Eventually(session).WithTimeout(3 * time.Second).Should(Exit(1))
		})

		It("does not output banner until previous commands complete", func() {
			cmd := CreateCommand(`sleep 3 && echo %s | {{.pcBinary}} --skip-validation | xargs -n 1 echo -`, shellescape.Quote("foo\nbar\nbaz"))
			session := StartCommand(cmd)

			Consistently(session.Err).WithTimeout(2 * time.Second).ShouldNot(Say("PipeCheck: The following was read in and will be passed through:"))
			Eventually(session.Err).WithTimeout(2 * time.Second).Should(Say("PipeCheck: The following was read in and will be passed through:"))
			Eventually(session).Should(Exit(0))
		})

		When("skip-validation flag is passed", func() {
			It("outputs and continues successfully", func() {
				cmd := CreateCommand(`echo %s | {{.pcBinary}} --skip-validation | xargs -n 1 echo -`, shellescape.Quote("foo\nbar\nbaz"))
				session := StartCommand(cmd)

				Eventually(session.Err).Should(Say("PipeCheck: The following was read in and will be passed through:"))

				Eventually(session.Err).Should(Say(`foo\n`))
				Eventually(session.Err).Should(Say(`bar\n`))
				Eventually(session.Err).Should(Say(`baz\n=`))

				Eventually(session.Err).Should(Say(`Skipping Validation, will proceed to next command`))

				Eventually(session).Should(Say(`- foo\n`))
				Eventually(session).Should(Say(`- bar\n`))
				Eventually(session).Should(Say(`- baz\n`))

				Eventually(session).Should(Exit(0))
			})
		})
	})
})
