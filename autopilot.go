package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/xchapter7x/autopilot/rewind"
)

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(&AutopilotPlugin{})
}

//AutopilotPlugin - the object implementing the plugin for zdd
type AutopilotPlugin struct{}

//Run - required command of a plugin (entry point)
func (plugin AutopilotPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "push-zdd" && len(args) == 1 {
		fmt.Println(`USAGE:
   Push a single app (with or without a manifest):
   cf push APP_NAME [-b BUILDPACK_NAME] [-c COMMAND] [-d DOMAIN] [-f MANIFEST_PATH]
   [-i NUM_INSTANCES] [-k DISK] [-m MEMORY] [-n HOST] [-p PATH] [-s STACK] [-t TIMEOUT]
   [--no-hostname] [--no-manifest] [--no-route] [--no-start]`)
		return
	}

	var err error
	appRepo := NewApplicationRepo(cliConnection)
	appName, argList := ParseArgs(args)
	venerableAppName := appName + "-venerable"

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

//ApplicationRepo - cli connection wrapper
type ApplicationRepo struct {
	conn plugin.CliConnection
}

//NewApplicationRepo - constructor function to create cli connection wrapper
func NewApplicationRepo(conn plugin.CliConnection) *ApplicationRepo {
	return &ApplicationRepo{
		conn: conn,
	}
}

//RenameApplication - rename the application given
func (repo *ApplicationRepo) RenameApplication(oldName, newName string) error {
	_, err := repo.conn.CliCommand("rename", oldName, newName)
	return err
}

//PushApplication - push the application to cf
func (repo *ApplicationRepo) PushApplication(args []string) error {
	_, err := repo.conn.CliCommand(args...)
	return err
}

//DeleteApplication - delete the application from cf
func (repo *ApplicationRepo) DeleteApplication(appName string) error {
	_, err := repo.conn.CliCommand("delete", appName, "-f")
	return err
}

//ListApplications - list applications on cf
func (repo *ApplicationRepo) ListApplications() error {
	_, err := repo.conn.CliCommand("apps")
	return err
}
