package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const runnerScript = `#!/usr/bin/env bash
	_int_signal() {
	  kill -INT "$child" 2>/dev/null
	}
	_term_signal() {
	  kill -TERM "$child" 2>/dev/null
	}

	trap _int_signal SIGINT
	trap _term_signal SIGTERM

	%s &

	child=$!
	wait "$child"`

const runnerScript2 = `#!/usr/bin/env bash
	trap 'kill -INT $(jobs -p)' SIGINT
	trap 'kill -TERM $(jobs -p)' SIGTERM

	(%s &)

	wait`

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

// CreateCommand will return a command with `{{.pcCommand}}` replaced with the
// generated binary.
func CreateCommand(tt string, args ...any) *exec.Cmd {
	f, err := ioutil.TempFile("", "pc-int-test-")
	Expect(err).ToNot(HaveOccurred())
	defer f.Close()
	defer f.Chmod(0755)
	DeferCleanup(func() {
		os.RemoveAll(f.Name())
	})

	pipedCommands := fmt.Sprintf(tt, args...)
	t := template.Must(template.New("test-script").Parse(fmt.Sprintf(runnerScript2, pipedCommands)))
	err = t.Execute(f, map[string]string{
		"pcBinary": pcBinary,
	})
	Expect(err).ToNot(HaveOccurred())
	return exec.Command(f.Name())
}

// StartCommand will start the provided command, wrapping in it a *gexec.Session.
func StartCommand(command *exec.Cmd) *Session {
	session, err := Start(
		command,
		NewPrefixedWriter("OUT: ", GinkgoWriter),
		NewPrefixedWriter("ERR: ", GinkgoWriter),
	)
	Expect(err).ToNot(HaveOccurred())
	return session
}
