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

//AutopilotPlugin - the object implementing the plugin for zdd
type AutopilotPlugin struct {
	appRepo          *application_repo.ApplicationRepo
	appName          string
	venerableAppName string
}

func main() {
	plugin.Start(&AutopilotPlugin{})
}

//Run - required command of a plugin (entry point)
func (plugin AutopilotPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	var err error
	plugin.appRepo = application_repo.NewApplicationRepo(cliConnection)

	if args[0] == "push-zdd" && len(args) == 1 {
		err = plugin.appRepo.PushApplication([]string{"push", "-h"})
		fatalIf(err)
		return
	}

	appName, argList := ParseArgs(args)
	plugin.appName = appName
	plugin.venerableAppName = appName + "-venerable"

	actions := rewind.Actions{
		Actions:              plugin.getActions(argList),
		RewindFailureMessage: "Oh no. Something's gone wrong. I've tried to roll back but you should check to see if everything is OK.",
	}

	err = actions.Execute()
	fatalIf(err)

	fmt.Printf("\nA new version of your application has successfully been pushed!\n\n")

	err = plugin.appRepo.ListApplications()
	fatalIf(err)
}

func (plugin AutopilotPlugin) getActions(argList []string) (actionList []rewind.Action) {
	actionList = []rewind.Action{plugin.getPushAction(argList)}

	if appExists(getAppList(plugin.appRepo), plugin.appName) {
		fmt.Printf("\n%s was found, using zero-downtime-deployment\n\n", plugin.appName)
		actionList = []rewind.Action{
			plugin.getRenameAction(),
			plugin.getPushAction(argList),
			plugin.getDeleteAction(),
		}

		plugin.addReversePrevious(&actionList[1])

	}
	return
}

func (plugin AutopilotPlugin) getPushAction(argList []string) rewind.Action {
	return rewind.Action{
		Forward: func() error {
			return plugin.appRepo.PushApplication(argList)
		},
	}
}

func (plugin AutopilotPlugin) addReversePrevious(action *rewind.Action) {
	action.ReversePrevious = func() error {
		plugin.appRepo.DeleteApplication(plugin.appName)

		return plugin.appRepo.RenameApplication(plugin.venerableAppName, plugin.appName)
	}
}

func (plugin AutopilotPlugin) getRenameAction() rewind.Action {
	return rewind.Action{
		Forward: func() error {
			return plugin.appRepo.RenameApplication(plugin.appName, plugin.venerableAppName)
		},
	}
}

func (plugin AutopilotPlugin) getDeleteAction() rewind.Action {
	return rewind.Action{
		Forward: func() error {
			return plugin.appRepo.DeleteApplication(plugin.venerableAppName)
		},
	}
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

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

//appExists - check if appName is in output
func appExists(output []string, appName string) bool {
	for _, app := range output {
		if strings.Contains(app, appName) {
			return true
		}
	}
	return false
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
