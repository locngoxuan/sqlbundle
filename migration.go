package sqlbundle

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MigrationScript struct {
	Group      string
	Artifact   string
	Version    string
	FilePath   string
	FileName   string
	NextScript *MigrationScript
	IsApplied  bool
}

type MigrationSQL struct {
	Group    string
	Artifact string
	Version  string
	FilePath string
	FileName string
}

func (ms *MigrationScript) ListAll() {
	s := ms
	for s != nil {
		printInfo(s.Group, s.Artifact, s.Version, s.FilePath)
		s = s.NextScript
	}
}

func (ms *MigrationScript) notAppliedYet() []MigrationSQL {
	s := ms
	paths := make([]MigrationSQL, 0)
	for s != nil {
		if s.IsApplied || isEmpty(s.FilePath) {
			s = s.NextScript
			continue
		}
		paths = append(paths, MigrationSQL{
			Group:    s.Group,
			Artifact: s.Artifact,
			Version:  s.Version,
			FilePath: s.FilePath,
			FileName: s.FileName,
		})
		s = s.NextScript
	}
	return paths
}

func (ms *MigrationScript) markAsApplied(depName, fileName string) {
	s := ms
	for s != nil{
		if fmt.Sprintf("%s.%s", s.Group, s.Artifact) == depName && s.FileName == fileName {
			s.IsApplied = true
			break
		}
		s = s.NextScript
	}
}

func (ms *MigrationScript) append(script *MigrationScript) {
	s := ms
	script.IsApplied = false
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
			Version:   version,
			Group:     group,
			Artifact:  artifact,
			FilePath:  path,
			FileName:  fileName,
			IsApplied: false,
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

func parseStatements(filePath string, up bool) (stmts []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)
	stateMachine := PARSER_START

	for scanner.Scan() {
		if (stateMachine == PARSER_UP_END && up) || (stateMachine == PARSER_DOWN_END && !up) {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(strings.TrimPrefix(line, "--"))
			if strings.HasPrefix(cmd, "+up BEGIN") && up {
				stateMachine = PARSER_UP_BEGIN
			} else if strings.HasPrefix(cmd, "+up END") && stateMachine == PARSER_UP_BEGIN {
				stateMachine = PARSER_UP_END
			} else if strings.HasPrefix(cmd, "+down BEGIN") && !up {
				stateMachine = PARSER_DOWN_BEGIN
			} else if strings.HasPrefix(cmd, "+down END") && stateMachine == PARSER_DOWN_BEGIN {
				stateMachine = PARSER_DOWN_END
			} else {
				// ignore comment
			}
			continue
		}

		line = strings.TrimSpace(line)
		if strings.TrimSpace(line) == "" {
			//ignore empty line
			continue
		}

		if _, err = buf.WriteString(line + " "); err != nil {
			break
			//return nil, false, errors.Wrap(err, "failed to write to buf")
		}

		if strings.HasSuffix(line, ";") {
			statement := buf.String()
			buf.Reset()
			stmts = append(stmts, statement)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// EOF
	return
}
