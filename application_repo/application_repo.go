package application_repo

import "github.com/cloudfoundry/cli/plugin"

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

//ListApplicationsWithOutput - list applications on cf with output
func (repo *ApplicationRepo) ListApplicationsWithOutput() ([]string, error) {
	output, err := repo.conn.CliCommand("apps")
	return output, err
}
