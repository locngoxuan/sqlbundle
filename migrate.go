package sqlbundle

import (
	"database/sql"
	"fmt"
)

func createVersionTable(db *sql.DB) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	d := GetDialect()

	stmts := d.createTable()
	for _, stmt := range stmts {
		if _, err = txn.Exec(stmt); err != nil {
			printInfo(fmt.Sprintf("Fail to execute create table statements %s", stmt))
			_ = txn.Rollback()
			return err
		}
	}
	if err = txn.Commit(); err != nil {
		_ = txn.Rollback()
		return err
	}
	return nil
}

func QueryDatabaseVersions(db *sql.DB) ([]DbVersion, error) {
	rows, err := GetDialect().dbVersionQuery(db)
	if err != nil {
		return []DbVersion{}, createVersionTable(db)
	}
	defer func() {
		_ = rows.Close()
	}()

	versions := make([]DbVersion, 0)

	for rows.Next() {
		var row DbVersion
		if err = rows.Scan(&row.Id, &row.Version); err != nil {
			return nil, err
		}
		versions = append(versions, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func QueryDatabaseHistories(db *sql.DB) ([]DbHistory, error) {
	rows, err := GetDialect().dbHistoryQuery(db)
	if err != nil {
		return []DbHistory{}, createVersionTable(db)
	}
	defer func() {
		_ = rows.Close()
	}()

	histories := make([]DbHistory, 0)

	for rows.Next() {
		var row DbHistory
		if err = rows.Scan(&row.Id, &row.Version, &row.DepName, &row.DepVersion, &row.File, &row.CheckSum); err != nil {
			return nil, err
		}
		histories = append(histories, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}
