package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/xchapter7x/autopilot"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/cloudfoundry/cli/testhelpers/io"
)

var _ = Describe("Flag Parsing", func() {
	It("parses a complete set of args", func() {
		appName, args := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "manifest-path",
				"-p", "app-path",
			},
		)
		Ω(appName).Should(Equal("appname"))
		Ω(args).Should(Equal([]string{
			"push",
			"appname",
			"-f", "manifest-path",
			"-p", "app-path",
		}))
	})
})

var _ = Describe("Command Syntax", func() {

	var (
		cliConn         *fakes.FakeCliConnection
		autopilotPlugin *AutopilotPlugin
	)

	BeforeEach(func() {
		cliConn = &fakes.FakeCliConnection{}
		autopilotPlugin = &AutopilotPlugin{}
	})

	It("displays push usage when push-zdd called with no arguments", func() {
		autopilotPlugin.Run(cliConn, []string{"push-zdd"})

		Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
		args := cliConn.CliCommandArgsForCall(0)
		Ω(args).Should(Equal([]string{"push", "-h"}))
	})

	It("can push an app that already exists", func() {
		cliConn.CliCommandReturns([]string{"myapp and-other-stuff"}, nil)
		output := CaptureOutput(func() {
			autopilotPlugin.Run(cliConn, []string{"push-zdd", "myapp"})
		})

		Expect(len(ActionList)).To(Equal(3))
		Expect(output).To(ContainElement(ContainSubstring("using zero-downtime-deployment")))
	})

	It("can push an app that doesn't exist", func() {
		cliConn.CliCommandReturns([]string{"some-other-app and-other-stuff"}, nil)
		output := CaptureOutput(func() {
			autopilotPlugin.Run(cliConn, []string{"push-zdd", "my-new-app"})
		})

		Expect(len(ActionList)).To(Equal(1))
		Expect(output).ToNot(ContainElement(ContainSubstring("using zero-downtime-deployment")))
		Expect(output).To(ContainElement(ContainSubstring("new version of your application")))
	})
})
