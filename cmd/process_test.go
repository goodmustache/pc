package cmd_test

import (
	"io"
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/goodmustache/pc/cmd"
	"github.com/goodmustache/pc/cmd/internal"
)

var _ = Describe("Process", func() {
	var (
		process cmd.Process

		fakeIn  *Buffer
		fakeOut *Buffer
		fakeErr *Buffer

		pipedInput []byte
		executeErr error

		testInput func(c rune) error
	)

	BeforeEach(func() {
		fakeIn, fakeOut, fakeErr = NewBuffer(), NewBuffer(), NewBuffer()
		userInput := make(chan byte)

		process = cmd.Process{
			SkipValidation: true,
			StdIn:          fakeIn,
			StdOut:         fakeOut,
			StdErr:         fakeErr,

			GetUserInput: func(stream io.Writer) (byte, error) {
				_, _ = stream.Write([]byte(internal.ValidationMessage))
				return <-userInput, nil
			},
			OutputCompleteMessage: true,
		}

		pipedInput = []byte("I\nam\nsome-test\ninput")
		_, _ = fakeIn.Write(pipedInput)

		testInput = func(c rune) error {
			complete := make(chan any)

			var executeErr error
			go func() {
				defer GinkgoRecover()
				defer close(complete)

				executeErr = process.Execute(nil) // blocking
			}()

			Eventually(fakeErr).Should(Say(regexp.QuoteMeta(internal.ValidationMessage)))
			Consistently(fakeErr).ShouldNot(Say("complete"))
			userInput <- byte(c)
			Eventually(complete).Should(BeClosed())

			return executeErr
		}
	})

	When("version flag is passed", func() {
		BeforeEach(func() {
			process.Version = true
		})

		It("returns the version and exits", func() {
			executeErr = process.Execute(nil)
			Expect(executeErr).ToNot(HaveOccurred())
			Expect(fakeOut).To(Say(`PipeCheck \d+\.\d+\.\d+`))
			Expect(fakeErr.Contents()).To(BeEmpty())
		})
	})

	It("outputs banner to StdErr", func() {
		executeErr = process.Execute(nil)
		Expect(executeErr).ToNot(HaveOccurred())

		Expect(fakeErr).To(Say(internal.Bannor))
	})

	It("outputs StdIn to StdErr", func() {
		executeErr = process.Execute(nil)
		Expect(executeErr).ToNot(HaveOccurred())

		Expect(fakeErr).To(Say(string(pipedInput)))
	})

	It("outputs skipping validation message to StdErr", func() {
		process.GetUserInput = nil
		executeErr = process.Execute(nil)
		Expect(executeErr).ToNot(HaveOccurred())

		Eventually(fakeErr).Should(Say(regexp.QuoteMeta(internal.SkipValidationMessage)))
		Eventually(fakeErr).Should(Say("complete"))
	})

	When("skip validation is false (default)", func() {
		BeforeEach(func() {
			process.SkipValidation = false
		})

		It("outputs confirmation to continue to StdErr and blocks until user input", func() {
			executeErr = testInput('y')
			Expect(executeErr).ToNot(HaveOccurred())
		})

		When("when reading user input", func() {
			When("given an explicit 'y'", func() {
				It("passes the input to StdOut and returns no error", func() {
					executeErr = testInput('y')
					Expect(executeErr).ToNot(HaveOccurred())
					Expect(fakeOut.Contents()).To(Equal(pipedInput))
				})
			})

			When("given an explicit 'Y'", func() {
				It("passes the input to StdOut and returns no error", func() {
					executeErr = testInput('Y')
					Expect(executeErr).ToNot(HaveOccurred())
					Expect(fakeOut.Contents()).To(Equal(pipedInput))
				})
			})

			When("given any character that's not 'y/Y'", func() {
				It("breaks and returns an error", func() {
					executeErr = testInput('l')
					Expect(executeErr).To(MatchError(cmd.ErrExit))
					Expect(fakeOut.Contents()).To(BeEmpty())
				})
			})
		})
	})
})
