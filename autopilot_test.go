package main_test

import (
	"errors"

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

var _ = Describe("ApplicationRepo", func() {
	var (
		cliConn *fakes.FakeCliConnection
		repo    *ApplicationRepo
	)

	BeforeEach(func() {
		cliConn = &fakes.FakeCliConnection{}
		repo = NewApplicationRepo(cliConn)
	})

	Describe("RenameApplication", func() {
		It("renames the application", func() {
			err := repo.RenameApplication("old-name", "new-name")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{"rename", "old-name", "new-name"}))
		})

		It("returns an error if one occurs", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("no app"))

			err := repo.RenameApplication("old-name", "new-name")
			Ω(err).Should(MatchError("no app"))
		})
	})

	Describe("PushApplication", func() {
		It("pushes an application with both a manifest and a path", func() {
			err := repo.PushApplication([]string{"push", "myapp", "-f", "/path/to/a/manifest.yml", "-p", "/path/to/the/app"})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{
				"push",
				"myapp",
				"-f", "/path/to/a/manifest.yml",
				"-p", "/path/to/the/app",
			}))
		})
		It("pushes an application with only a manifest", func() {
			err := repo.PushApplication([]string{"push", "myapp", "-f", "/path/to/a/manifest.yml"})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{
				"push",
				"myapp",
				"-f", "/path/to/a/manifest.yml",
			}))
		})

		It("returns errors from the push", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad app"))

			err := repo.PushApplication([]string{"push", "myapp", "-f", "/path/to/a/manifest.yml", "-p", "/path/to/the/app"})
			Ω(err).Should(MatchError("bad app"))
		})

		It("pushes an application if one does not exist", func() {
			err := repo.PushApplication([]string{"push", "myapp"})
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{
				"push",
				"myapp",
			}))
		})
	})

	Describe("DeleteApplication", func() {
		It("deletes all trace of an application", func() {
			err := repo.DeleteApplication("app-name")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{
				"delete", "app-name",
				"-f",
			}))
		})

		It("returns errors from the delete", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad app"))

			err := repo.DeleteApplication("app-name")
			Ω(err).Should(MatchError("bad app"))
		})
	})

	Describe("ListApplications", func() {
		It("lists all the applications", func() {
			err := repo.ListApplications()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cliConn.CliCommandCallCount()).Should(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Ω(args).Should(Equal([]string{"apps"}))
		})

		It("returns errors from the list", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad apps"))

			err := repo.ListApplications()
			Ω(err).Should(MatchError("bad apps"))
		})
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

		Expect(output).To(ContainElement(ContainSubstring("new version of your application")))
	})

	It("can push an app that doesn't exist", func() {
		cliConn.CliCommandReturns([]string{"some-other-app and-other-stuff"}, nil)
		output := CaptureOutput(func() {
			autopilotPlugin.Run(cliConn, []string{"push-zdd", "my-new-app"})
		})

		Expect(output).To(ContainElement(ContainSubstring("Using cf push")))
	})
})
