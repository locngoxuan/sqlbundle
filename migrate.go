package sqlbundle

import "database/sql"

func GetDBVersion(r []int, db *sql.DB) (int64, error) {
	version, err := EnsureDBVersion(r, db)
	if err != nil {
		if err == ErrNoNextVersion{
			return 0, nil
		}
		return -1, err
	}

	return version, nil
}

func EnsureDBVersion(r []int, db *sql.DB) (int64, error) {
	rows, err := GetDialect().dbVersionQuery(r, db)
	if err != nil {
		return 0, createVersionTable(db)
	}
	defer rows.Close()

	// The most recent record for each migration specifies
	// whether it has been applied or rolled back.
	// The first version we find that has been applied is the current version.
	toSkip := make([]int64, 0)

	for rows.Next() {
		var row MigrationRecord
		if err = rows.Scan(&row.VersionID, &row.IsApplied); err != nil {
			return 0, errors.Wrap(err, "failed to scan row")
		}

		// have we already marked this version to be skipped?
		skip := false
		for _, v := range toSkip {
			if v == row.VersionID {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		// if version has been applied we're done
		if row.IsApplied {
			return row.VersionID, nil
		}

		// latest version of migration has not been applied.
		toSkip = append(toSkip, row.VersionID)
	}

	if err := rows.Err(); err != nil {
		return 0, errors.Wrap(err, "failed to get next row")
	}

	return 0, ErrNoNextVersion
}
