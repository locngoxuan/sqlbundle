package main

import (
	"fmt"
	"os"
	"github.com/locngoxuan/sqlbundle"
)

var version = "1.3.0"

func main() {
	sqlbundle.SetVersion(version)
	sqlbundle.SetLogWriter(os.Stdout)
	cmd, err := sqlbundle.ReadArgument()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	sb, err := sqlbundle.NewSQLBundle(cmd.Argument)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = sqlbundle.Handle(cmd.Command, sb)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
