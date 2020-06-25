package sqlbundle

import (
	"fmt"
	"strings"
)

func (sb *SQLBundle) Downgrade() error {
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
		AppVersion: sb.ReadVersion(),
		Version:    sb.ReadVersion(),
		Group:      sb.Config.GroupId,
		Artifact:   sb.Config.ArtifactId,
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

	if versions == nil || len(versions) == 0 {
		printInfo("not found any version for downgrading")
		return nil
	}

	downgrades := make([]string, 0)
	set := make(map[string]struct{})
	if isEmpty(sb.Argument.Version) {
		downgrades = append(downgrades, versions[0].Version)
		set[versions[0].Version] = struct{}{}
	} else {
		downGradeVer := sb.Argument.Version
		for _, v := range versions {
			if v.Version == downGradeVer {
				break
			}
			downgrades = append(downgrades, v.Version)
			set[v.Version] = struct{}{}
		}
	}

	if len(downgrades) == 0{
		printInfo("not found any version for downgrading")
		return nil
	}

	histories, err := QueryDatabaseHistories(db)
	if err != nil {
		return err
	}

	kept := make([]DbHistory, 0)
	applied := make(map[string]struct{})
	for _, h := range histories {
		applied[fmt.Sprintf("%s-%s", h.DepName, h.File)] = struct{}{}
		_, ok := set[h.Version]
		if !ok {
			kept = append(kept, h)
		}
	}

	for _, h := range kept {
		script.ignore(h.DepName, h.File)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sqlFiles := script.notIgnored()
	for i, j := 0, len(sqlFiles)-1; i < j; i, j = i+1, j-1 {
		sqlFiles[i], sqlFiles[j] = sqlFiles[j], sqlFiles[i]
	}
	d := GetDialect()
	historyStatement, err := tx.Prepare(d.deleteHistory())
	if err != nil {
		return err
	}

	for _, sql := range sqlFiles {
		_, ok := applied[fmt.Sprintf("%s.%s-%s", sql.Group, sql.Artifact, sql.FileName)]
		if !ok {
			continue
		}
		statements, err := parseStatements(sql.FilePath, false)
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
		printInfo(fmt.Sprintf("Redo%s%s", strings.Repeat(" ", 10), sql.FileName))
		_, err = historyStatement.Exec(fmt.Sprintf("%s.%s", sql.Group, sql.Artifact), sql.Version, sql.FileName)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	versionStatement, err := tx.Prepare(d.deleteVersion())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, ver := range downgrades {
		_, err = versionStatement.Exec(ver)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
