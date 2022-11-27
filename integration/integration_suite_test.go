package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	pcBinary string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	pcBinary, err = Build("github.com/goodmustache/pc")
	Expect(err).ToNot(HaveOccurred())
	DeferCleanup(CleanupBuildArtifacts)
	return nil
}, func(address []byte) {})

func StartCommand(command *exec.Cmd) *Session {
	session, err := Start(
		command,
		NewPrefixedWriter("OUT: ", GinkgoWriter),
		NewPrefixedWriter("ERR: ", GinkgoWriter),
	)
	Expect(err).ToNot(HaveOccurred())
	return session
}
