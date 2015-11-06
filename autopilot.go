package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/xchapter7x/autopilot/application_repo"
	"github.com/xchapter7x/autopilot/rewind"
)

var ActionList []rewind.Action

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func appExists(output []string, appName string) bool {
	for _, app := range output {
		if strings.Contains(app, appName) {
			return true
		}
	}
	return false
}

func main() {
	plugin.Start(&AutopilotPlugin{})
}

//AutopilotPlugin - the object implementing the plugin for zdd
type AutopilotPlugin struct{}

//Run - required command of a plugin (entry point)
func (plugin AutopilotPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	var err error
	appRepo := application_repo.NewApplicationRepo(cliConnection)

	if args[0] == "push-zdd" && len(args) == 1 {
		err = appRepo.PushApplication([]string{"push", "-h"})
		fatalIf(err)
		return
	}

	appName, argList := ParseArgs(args)
	venerableAppName := appName + "-venerable"

	output, err := appRepo.ListApplicationsWithOutput()
	fatalIf(err)

	if appNotFound(output, appName) {
		fmt.Printf("\n%s not found! Using cf push.\n\n", appName)
		err := appRepo.PushApplication(argList)
		fatalIf(err)

		fmt.Printf("\nYour application has successfully been pushed!\n\n")
		return
	}

	actions := rewind.Actions{
		Actions: []rewind.Action{
			// rename
			{
				Forward: func() error {
					return appRepo.RenameApplication(appName, venerableAppName)
				},
			},

			// push
			{
				Forward: func() error {
					return appRepo.PushApplication(argList)
				},
				ReversePrevious: func() error {
					appRepo.DeleteApplication(appName)
					return appRepo.RenameApplication(venerableAppName, appName)
				},
			},

			// delete
			{
				Forward: func() error {
					return appRepo.DeleteApplication(venerableAppName)
				},
			},
		},
		RewindFailureMessage: "Oh no. Something's gone wrong. I've tried to roll back but you should check to see if everything is OK.",
	}

	err = actions.Execute()
	fatalIf(err)

	fmt.Printf("\nA new version of your application has successfully been pushed!\n\n")

	err = appRepo.ListApplications()
	fatalIf(err)
}

//GetMetadata - required command of plugin (returns meta data about plugin)
func (AutopilotPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "push-zero-downtime-deployment",
		Commands: []plugin.Command{
			{
				Name:     "push-zdd",
				HelpText: "Perform a zero-downtime push of an application over the top of an old one",
			},
		},
	}
}

//ParseArgs - parse given cli arguments
func ParseArgs(args []string) (string, []string) {
	args[0] = "push"
	appName := args[1]
	return appName, args
}

//ErrNoManifest - error to return when there is no manifest if required
var ErrNoManifest = errors.New("a manifest is required to push this application")

func getAppList(appRepo *application_repo.ApplicationRepo) []string {
	output, err := appRepo.ListApplicationsWithOutput()
	fatalIf(err)
	return output
}
