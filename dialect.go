package sqlbundle

import (
	"database/sql"
	"fmt"
	"time"
)

type DbVersion struct {
	Id        int
	Version   string
	Timestamp time.Time
}

type DbHistory struct {
	Id         int
	Version    string
	DepName    string
	DepVersion string
	File       string
	Timestamp  time.Time
}

type SQLDialect interface {
	dbVersionQuery(db *sql.DB) (*sql.Rows, error)
	dbHistoryQuery(db *sql.DB) (*sql.Rows, error)
	createTable() string
	insertVersion() string
	insertHistory() string
	deleteHistory() string
	deleteVersion() string
}

var dialect SQLDialect = &PostgresDialect{}

////////////////////////////
// Postgres
////////////////////////////

type PostgresDialect struct {}

func (pg PostgresDialect) createTable() string {
	return `CREATE TABLE db_versions (
            	id serial NOT NULL,
                version varchar(255) NOT NULL,
                timestamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );

			CREATE TABLE db_histories (
            	id serial NOT NULL,
                version varchar(255) NOT NULL,
				dep_name text,
				dep_version varchar(255),
				file_name text NOT NULL,
                timestamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`
}

func (pg PostgresDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT id, version from db_versions ORDER BY id DESC`))
	if err != nil {
		return nil, err
	}
	return rows, err
}

func (pg PostgresDialect) dbHistoryQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT id, version, dep_name, dep_version, file_name from db_histories ORDER BY id DESC`))
	if err != nil {
		return nil, err
	}
	return rows, err
}

func (pg PostgresDialect) insertVersion() string {
	return fmt.Sprintf("INSERT INTO db_versions (version) VALUES ($1)")
}

func (pg PostgresDialect) insertHistory() string {
	return fmt.Sprintf("INSERT INTO db_histories (version, dep_name, dep_version, file_name) VALUES ($1, $2, $3, $4)")
}

func (pg PostgresDialect) deleteVersion() string {
	return fmt.Sprintf("DELETE FROM db_versions WHERE version = $1")
}

func (pg PostgresDialect) deleteHistory() string {
	return fmt.Sprintf("DELETE FROM db_histories WHERE dep_name = $1 AND dep_version = $2 AND file_name = $3")
}

func GetDialect() SQLDialect {
	return dialect
}

// SetDialect sets the SQLDialect
func SetDialect(d string) error {
	switch d {
	case "postgres":
		dialect = &PostgresDialect{}
		break
	case "oracle":
		//dialect = &OracleDialect{}
	default:
		return fmt.Errorf("%q: unknown dialect", d)
	}

	return nil
}
