package application_repo_test

import (
	"errors"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/autopilot/application_repo"
)

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
		Context("when we called", func() {
			It("then it should call the get apps list api", func() {
				err := repo.ListApplications()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cliConn.GetAppsCallCount()).Should(Equal(1))
			})
		})
		Context("when the get apps call yields an error", func() {
			It("then it should return the errors ", func() {
				cliConn.GetAppsReturns(nil, errors.New("bad apps"))
				err := repo.ListApplications()
				Ω(err).Should(MatchError("bad apps"))
			})
		})
	})
})
