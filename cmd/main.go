package main

import (
	"sqlbundle"
)

func main() {
	cmd, err := sqlbundle.ReadArgument()
	if err != nil {
		panic(err)
	}
	sb, err := sqlbundle.NewSQLBundle(cmd.Argument)
	if err != nil {
		panic(err)
	}
	err = sqlbundle.Handle(cmd.Command, sb)
	if err != nil {
		panic(err)
	}
}
