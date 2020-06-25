package main

import (
	"fmt"
	"os"
	"sqlbundle"
)

func main() {
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
