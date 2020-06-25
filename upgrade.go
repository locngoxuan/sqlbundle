package sqlbundle

import (
	"fmt"
	"strings"
)

func (sb *SQLBundle) Upgrade() error {
	db, err := OpenDBWithDriver(sb.Argument.DBDriver, sb.Argument.DBString)
	if err != nil {
		return err
	}

	defer func() {
		_ = db.Close()
	}()

	err = sb.readConfig()
	if err != nil {
		return err
	}

	script := &MigrationScript{
		Version:  sb.ReadVersion(),
		Group:    sb.Config.GroupId,
		Artifact: sb.Config.ArtifactId,
	}

	err = collectMigrations(*sb, script)
	if err != nil {
		return err
	}
	//script.ListAll()

	versions, err := QueryDatabaseVersions(db)
	if err != nil {
		return err
	}

	currentVersion := sb.ReadVersion()
	for _, v := range versions {
		if v.Version == currentVersion {
			return fmt.Errorf("version %s was already installed", currentVersion)
		}
	}

	histories, err := QueryDatabaseHistories(db)
	if err != nil {
		return err
	}

	for _, h := range histories {
		script.markAsApplied(h.DepName, h.File)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sqlFiles := script.notAppliedYet()

	d := GetDialect()
	historyStatement, err := tx.Prepare(d.insertHistory())
	if err != nil{
		return err
	}

	for _, sql := range sqlFiles {
		statements, err := parseStatements(sql.FilePath, true)
		if err != nil {
			return err
		}

		for _, statement := range statements {
			if _, err = tx.Exec(statement); err != nil {
				printInfo(fmt.Sprintf("Fail to execute query %s", statements), err)
				printInfo("Rollback transaction")
				_ = tx.Rollback()
				return err
			}
		}
		printInfo(fmt.Sprintf("Apply%s%s", strings.Repeat(" ", 10), sql.FileName))
		_, err = historyStatement.Exec(sb.ReadVersion(), fmt.Sprintf("%s.%s", sql.Group, sql.Artifact), sql.Version, sql.FileName)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	versionStatement, err := tx.Prepare(d.insertVersion())
	if err != nil{
		_ = tx.Rollback()
		return err
	}

	_, err = versionStatement.Exec(sb.ReadVersion())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
