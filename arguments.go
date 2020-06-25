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
	Workdir    string
	Filename   string
	Force      bool
	Repository string
	Username   string
	Password   string
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
		f.StringVar(&cmd.Argument.Version, "version", "", "")
		f.StringVar(&cmd.Argument.Group, "group", "", "")
		f.StringVar(&cmd.Argument.Artifact, "artifact", "", "")
		f.StringVar(&cmd.Argument.Workdir, "workdir", "", "")
		f.StringVar(&cmd.Argument.Filename, "file", "", "")
		f.BoolVar(&cmd.Argument.Force, "force", false, "")
		f.StringVar(&cmd.Argument.Repository, "repo", "", "")
		f.StringVar(&cmd.Argument.Username, "user", "", "")
		f.StringVar(&cmd.Argument.Password, "pass", "", "")
		err = f.Parse(os.Args[2:])
	}
	return
}
