package sqlbundle

import (
	"database/sql"
	"fmt"
)

// OpenDBWithDriver creates a connection a database, and modifies goose
// internals to be compatible with the supplied driver by calling SetDialect.
func OpenDBWithDriver(driver string, dbstring string) (*sql.DB, error) {
	if err := SetDialect(driver); err != nil {
		return nil, err
	}

	if driver == "oracle" {
		driver = "godror"
	}

	if driver == "sqlite" {
		driver = "sqlite3"
	}

	if driver == "postgresql" || driver == "pql" {
		driver = "postgres"
	}

	switch driver {
	case "postgres", "godror", "sqlite3":
		return sql.Open(driver, dbstring)
	default:
		return nil, fmt.Errorf("unsupported driver %s", driver)
	}
}
