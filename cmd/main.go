package main

import "sqlbundle"

func main() {
	sb, err := sqlbundle.NewSQLBundle("/home/deadpool/workspace/goproj/sqlbundle/example")
	if err != nil {
		panic(err)
	}
	err = sb.Create("1_example.sql")
	if err != nil{
		panic(err)
	}
}
