package main_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/xchapter7x/autopilot"

	"github.com/cloudfoundry/cli/plugin/fakes"
	"github.com/cloudfoundry/cli/plugin/models"
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

	Context("when a version of an app already exists", func() {
		var (
			controlAppName          = "myapp"
			controlAppNameVenerable = fmt.Sprintf("%s-venerable", controlAppName)
			controlCallChain        = [][]string{
				[]string{"rename", controlAppName, controlAppNameVenerable},
				[]string{"push", controlAppName},
				[]string{"delete", controlAppNameVenerable, "-f"},
			}
			controlCallCount = len(controlCallChain)
		)

		BeforeEach(func() {
			cliConn.GetAppsReturns([]plugin_models.GetAppsModel{
				plugin_models.GetAppsModel{
					Name: controlAppName,
				},
				plugin_models.GetAppsModel{
					Name: "and-other-stuff",
				},
			}, nil)
			autopilotPlugin.Run(cliConn, []string{"push-zdd", controlAppName})
		})
		It("then it should rename the existing app to venerable", func() {
			callCount := cliConn.CliCommandCallCount()

			for i := 0; i < callCount; i++ {
				called := cliConn.CliCommandArgsForCall(i)
				fmt.Println(called)
				Ω(true).Should(BeTrue())
			}
		})

		It("then it should rename the existing app to venerable", func() {
			callCount := cliConn.CliCommandCallCount()
			Ω(callCount).Should(Equal(controlCallCount))

			for i := 0; i < callCount; i++ {
				called := cliConn.CliCommandArgsForCall(i)
				Ω(called).Should(Equal(controlCallChain[i]))
			}
		})

		It("then it should push the new version of the application", func() {
			callCount := cliConn.CliCommandCallCount()
			Ω(callCount).Should(Equal(controlCallCount))

			for i := 0; i < callCount; i++ {
				called := cliConn.CliCommandArgsForCall(i)
				Ω(called).Should(Equal(controlCallChain[i]))
			}
		})

		It("then it should remove the venerable version of the application", func() {
			callCount := cliConn.CliCommandCallCount()
			Ω(callCount).Should(Equal(controlCallCount))

			for i := 0; i < callCount; i++ {
				called := cliConn.CliCommandArgsForCall(i)
				Ω(called).Should(Equal(controlCallChain[i]))
			}
		})
	})

	Context("when an app does not yet exist", func() {
		var (
			controlAppName   = "my-new-app"
			controlCallChain = [][]string{
				[]string{"push", controlAppName},
			}
			controlCallCount = len(controlCallChain)
		)
		BeforeEach(func() {
			app1 := "myapp"
			app2 := "other-app"
			cliConn.GetAppsReturns([]plugin_models.GetAppsModel{
				plugin_models.GetAppsModel{
					Name: app1,
				},
				plugin_models.GetAppsModel{
					Name: app2,
				},
			}, nil)
			autopilotPlugin.Run(cliConn, []string{"push-zdd", controlAppName})
		})
		It("then it should only call push", func() {
			callCount := cliConn.CliCommandCallCount()

			Ω(callCount).Should(Equal(controlCallCount))

			for i := 0; i < callCount; i++ {
				called := cliConn.CliCommandArgsForCall(i)
				Ω(called).Should(Equal(controlCallChain[i]))
			}
		})
	})
})
