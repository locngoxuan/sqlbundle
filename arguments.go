package sqlbundle

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	usagePrefix = `Usage: sqlbundle COMMAND [OPTIONS]

COMMAND:
  init          Init new sql project		
  create        Create a new sql file
  install       Download and install dependencies into deps directory
  pack          Packing
  clean         Remove build directory
  publish       Deploy package to repository
  version       Print version of sql-bundle
  upgrade       Upgrade database to latest version 
  downgrade     Downgrade database to previous version or any specific version
  help          Print usage info

Examples:

Options:
`
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

var f *flag.FlagSet

func ReadArgument() (cmd Command, err error) {
	f = flag.NewFlagSet("sqlbundle", flag.ContinueOnError)
	f.SetOutput(os.Stdout)
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
	f.Usage = func() {
		/**
		Do nothing
		 */
		_, _ = fmt.Fprint(f.Output(), usagePrefix)
		f.PrintDefaults()
		os.Exit(1)
	}
	if len(os.Args) == 1 {
		f.Usage()
		return
	}

	cmd.Command = strings.TrimSpace(os.Args[1])
	cmd.Argument = Argument{}
	if len(os.Args) > 2 {
		err = f.Parse(os.Args[2:])
	}
	return
}
