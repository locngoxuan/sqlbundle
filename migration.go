package sqlbundle

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MigrationScript struct {
	Group      string
	Artifact   string
	Version    string
	FileName   string
	NextScript *MigrationScript
}

func (ms *MigrationScript) ListAll() {
	s := ms
	for s != nil {
		printInfo(s.Group, s.Artifact, s.Version, s.FileName)
		s = s.NextScript
	}
}

func (ms *MigrationScript) append(script *MigrationScript) {
	s := ms
	for {
		next := s.NextScript
		if next == nil {
			s.NextScript = script
			break
		}
		s = next
	}
}

func collectSql(script *MigrationScript, dir, group, artifact, version string) (err error) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		fileExt := filepath.Ext(path)
		if fileExt != ".sql" {
			return nil
		}
		_, fileName := filepath.Split(path)
		next := &MigrationScript{
			Version:  version,
			Group:    group,
			Artifact: artifact,
			FileName: fileName,
		}
		script.append(next)
		return nil
	}
	if exists(dir) {
		err = filepath.Walk(dir, walkFunc)
		if err != nil {
			return
		}
	}
	return
}

func (sb *SQLBundle) collectMigrations(script *MigrationScript) (err error) {
	if sb.Config == nil {
		err = errors.New("can not read config")
		return
	}
	if sb.Config.Dependencies != nil && len(sb.Config.Dependencies) > 0 && exists(sb.DepsDir) {
		for _, depLink := range sb.Config.Dependencies {
			depName := path.Base(depLink)
			depName = strings.TrimSuffix(depName, filepath.Ext(depName))
			depPath := filepath.Join(sb.DepsDir, depName)
			if !exists(depPath) {
				err = errors.New("not found dependence " + depLink)
				break
			}

			bundle, err := NewSQLBundle(Argument{
				Workdir: depPath,
			})
			if err != nil {
				break
			}
			err = bundle.readConfig()
			if err != nil {
				break
			}
			err = bundle.collectMigrations(script)
			if err != nil {
				break
			}
		}
		if err != nil {
			return
		}
	}
	err = collectSql(script, sb.SourceDir, sb.Config.GroupId, sb.Config.GroupId, sb.ReadVersion())
	return nil
}
