package sqlbundle

import "database/sql"

type DbVersion struct {
	Id       int
	Group    string
	Artifact string
	Version  string
}

type DbHistory struct {
	Id       int
	Group    string
	Artifact string
	Version  string
	File     string
}

type SQLDialect interface {
	dbVersionQuery(db *sql.DB) (*sql.Rows, error)
	createDbVersion(db *sql.DB) error
}

//type PostgresDialect struct {
//}
//
//type OracleDialect struct {
//}
//
//var dialect SQLDialect = &PostgresDialect{}
var dialect SQLDialect

func GetDialect() SQLDialect {
	return dialect
}
