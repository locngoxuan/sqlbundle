package sqlbundle

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var sqlTemplate = `-- up statement
-- TODO: write sql statement here
-- end up

-- down statement
-- TODO: write sql statement here
-- end down`

type SQLBundle struct {
	WorkingDir     string
	SourceDir      string
	BuildDir       string
	BuildDependDir string
	BuildFinaldDir string
}

func NewSQLBundle(workDir string) (SQLBundle SQLBundle, err error) {
	if strings.TrimSpace(workDir) == "" {
		workDir, err = filepath.Abs(".")
	}
	SQLBundle.WorkingDir = workDir
	SQLBundle.SourceDir = filepath.Join(workDir, "src")
	SQLBundle.BuildDir = filepath.Join(workDir, "build")
	SQLBundle.BuildDependDir = filepath.Join(workDir, "build/depend")
	SQLBundle.BuildFinaldDir = filepath.Join(workDir, "build/final")
	return
}

func (sb *SQLBundle) Create(fileName string) error {
	err := os.MkdirAll(sb.SourceDir, 0755)
	if err != nil {
		return err
	}
	if filepath.Ext(fileName) != ".sql" {
		fileName = fmt.Sprintf("%s.sql", fileName)
	}
	fp := filepath.Join(sb.SourceDir, fileName)
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("can not create file " + fileName)
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.WriteString(f, sqlTemplate)
	return err
}

func (sb *SQLBundle) Clean() error {
	_ = os.RemoveAll(sb.BuildDir)
	return nil
}

func (sb *SQLBundle) DownloadDependencies() error {
	return nil
}

func (sb *SQLBundle) Build() error {
	err := os.MkdirAll(sb.BuildDependDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(sb.BuildFinaldDir, 0755)
	if err != nil {
		return err
	}
	err = sb.DownloadDependencies()
	if err != nil {
		return err
	}
	// generate sql
	return nil
}

func (sb *SQLBundle) Pack() error {
	return nil
}
