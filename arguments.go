package sqlbundle

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Command struct {
	Command string
	Argument
}

type Argument struct {
	Version    string
	Group      string
	Artifact   string
	WorkDir    string
	Filename   string
	Force      bool
	Repository string
	Username   string
	Password   string
	DBDriver   string
	DBString   string
}

func ReadArgument() (cmd Command, err error) {
	if len(os.Args) == 1 {
		err = errors.New("missing command")
		return
	}

	cmd.Command = strings.TrimSpace(os.Args[1])
	cmd.Argument = Argument{}
	if len(os.Args) > 2 {
		f := flag.NewFlagSet("sqlbundle", flag.ExitOnError)
		f.StringVar(&cmd.Argument.Version, "version", "", "version of database")
		f.StringVar(&cmd.Argument.Group, "group", "", "group name")
		f.StringVar(&cmd.Argument.Artifact, "artifact", "", "artifact name")
		f.StringVar(&cmd.Argument.WorkDir, "workdir", "", "working directory")
		f.StringVar(&cmd.Argument.Filename, "file", "", "name of sql file")
		f.BoolVar(&cmd.Argument.Force, "force", false, "force to delete before install")
		f.StringVar(&cmd.Argument.Repository, "repo", "", "address of repository")
		f.StringVar(&cmd.Argument.Username, "user", "", "repository credential")
		f.StringVar(&cmd.Argument.Password, "pass", "", "repository credential")
		f.StringVar(&cmd.Argument.DBDriver, "db-driver", "", "database driver")
		f.StringVar(&cmd.Argument.DBString, "db-connection", "", "connection string of database")
		err = f.Parse(os.Args[2:])
	}
	return
}
