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

	switch driver {
	case "postgres", "godror":
		db, err := sql.Open(driver, dbstring)
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		db.SetConnMaxLifetime(0)
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported driver %s", driver)
	}
}
