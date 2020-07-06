package sqlbundle

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MigrationScript struct {
	AppVersion string
	Group      string
	Artifact   string
	Version    string
	FilePath   string
	FileName   string
	NextScript *MigrationScript
	Ignore     bool
}

type MigrationSQL struct {
	AppVersion string
	Group      string
	Artifact   string
	Version    string
	FilePath   string
	FileName   string
}

func (ms *MigrationScript) ListAll() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Group", "Artifact", "Version", "File"})
	s := ms
	for s != nil {
		t.AppendRow(table.Row{s.Group, s.Artifact, s.Version, s.FileName})
		s = s.NextScript
	}
	t.Render()
}

func (ms *MigrationScript) notIgnored() []MigrationSQL {
	s := ms
	paths := make([]MigrationSQL, 0)
	for s != nil {
		if s.Ignore || isEmpty(s.FilePath) {
			s = s.NextScript
			continue
		}
		paths = append(paths, MigrationSQL{
			AppVersion: s.AppVersion,
			Group:      s.Group,
			Artifact:   s.Artifact,
			Version:    s.Version,
			FilePath:   s.FilePath,
			FileName:   s.FileName,
		})
		s = s.NextScript
	}
	return paths
}

func (ms *MigrationScript) ignore(depName, fileName string) {
	s := ms
	for s != nil {
		if fmt.Sprintf("%s.%s", s.Group, s.Artifact) == depName && s.FileName == fileName {
			s.Ignore = true
			break
		}
		s = s.NextScript
	}
}

func (ms *MigrationScript) append(script *MigrationScript) {
	s := ms
	script.Ignore = false
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
			AppVersion: script.AppVersion,
			Version:    version,
			Group:      group,
			Artifact:   artifact,
			FilePath:   path,
			FileName:   fileName,
			Ignore:     false,
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

func collectMigrations(sb SQLBundle, script *MigrationScript) (err error) {
	if sb.Config == nil {
		err = errors.New("can not read config")
		return
	}
	if sb.Config.Dependencies != nil && len(sb.Config.Dependencies) > 0 && exists(sb.DepsDir) {
		for _, depLink := range sb.Config.Dependencies {
			depName := path.Base(depLink)
			depName = strings.TrimSuffix(depName, filepath.Ext(depName))
			depName = fmt.Sprintf("%s_%s", depName, getMD5Hash(depLink))
			depPath := filepath.Join(sb.DepsDir, depName)
			if !exists(depPath) {
				err = errors.New("not found dependence " + depLink)
				break
			}

			bundle, err := NewSQLBundle(Argument{
				WorkDir: depPath,
			})
			if err != nil {
				break
			}
			err = bundle.readConfig()
			if err != nil {
				break
			}
			err = collectMigrations(bundle, script)
			if err != nil {
				break
			}
		}
		if err != nil {
			return
		}
	}
	err = collectSql(script, sb.SourceDir, sb.Config.GroupId, sb.Config.ArtifactId, sb.ReadVersion())
	return nil
}
