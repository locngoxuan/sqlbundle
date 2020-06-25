package sqlbundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var version = "1.0.0"
var sqlTemplate = `--+up BEGIN
-- TODO: write sql statement here
--+up END

--+down BEGIN
-- TODO: write sql statement here
--+down END`

type SQLBundle struct {
	Argument
	Config     *PackageJSON
	WorkingDir string
	SourceDir  string
	BuildDir   string
	DepsDir    string
	ConfigFile string
}

func NewSQLBundle(arg Argument) (bundle SQLBundle, err error) {
	workDir := arg.WorkDir
	if strings.TrimSpace(workDir) == "" {
		workDir, err = filepath.Abs(".")
	}
	bundle.Argument = arg
	bundle.WorkingDir = workDir
	bundle.SourceDir = filepath.Join(workDir, "src")
	bundle.BuildDir = filepath.Join(workDir, "build")
	bundle.DepsDir = filepath.Join(workDir, "deps")
	bundle.ConfigFile = filepath.Join(workDir, PACKAGE_JSON)
	return
}

func printInfo(v ...interface{}) {
	_, _ = fmt.Fprintln(os.Stdout, v...)
}

func Handle(command string, bundle SQLBundle) error {
	switch command {
	case "init":
		return bundle.Init()
	case "clean":
		return bundle.Clean()
	case "create":
		return bundle.Create()
	case "install":
		return bundle.Install()
	case "pack":
		return bundle.Pack()
	case "publish":
		return bundle.Publish()
	case "version":
		printInfo(fmt.Sprintf("Version %s", version))
		return nil
	case "upgrade":
		if isEmpty(bundle.Argument.DBDriver) || isEmpty(bundle.Argument.DBString) {
			return errors.New("missing database driver/connection string configuration")
		}
		return bundle.Upgrade()
	case "downgrade":
		if isEmpty(bundle.Argument.DBDriver) || isEmpty(bundle.Argument.DBString) {
			return errors.New("missing database driver/connection string configuration")
		}
		return bundle.Downgrade()
	case "help":
		f.Usage()
		return nil
	default:
		printInfo(fmt.Sprintf("command `%s` not found. Try `sqlbundle help` for more detail!", command))
	}
	return nil
}

func (sb *SQLBundle) ReadVersion() string {
	if isEmpty(sb.Argument.Version) {
		err := sb.readConfig()
		if err != nil {
			printInfo(err)
			os.Exit(1)
		}
		if isEmpty(sb.Config.Version) {
			printInfo("not found version from package.json")
			os.Exit(1)
		}
		return strings.TrimSpace(sb.Config.Version)
	}
	return strings.TrimSpace(sb.Argument.Version)
}

func (sb *SQLBundle) Init() error {
	err := os.MkdirAll(sb.SourceDir, 0755)
	if err != nil {
		return err
	}

	if strings.TrimSpace(sb.Group) == "" {
		return errors.New("missing group")
	}

	if strings.TrimSpace(sb.Artifact) == "" {
		return errors.New("missing artifact")
	}

	if strings.TrimSpace(sb.ConfigFile) == "" {
		return errors.New("missing version")
	}

	config := PackageJSON{
		GroupId:      sb.Group,
		ArtifactId:   sb.Artifact,
		Version:      sb.Argument.Version,
		Dependencies: make([]string, 0),
	}

	f, err := os.Create(sb.ConfigFile)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("can not create file " + sb.ConfigFile)
	}
	defer func() {
		_ = f.Close()
	}()

	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, string(bytes))
	return nil
}

func (sb *SQLBundle) Create() error {
	fileName := sb.Filename
	if strings.TrimSpace(fileName) == "" {
		return errors.New("missing filename")
	}
	err := os.MkdirAll(sb.SourceDir, 0755)
	if err != nil {
		return err
	}
	if filepath.Ext(fileName) != ".sql" {
		fileName = fmt.Sprintf("%s.sql", fileName)
	}

	fileName = fmt.Sprintf("%s_%s", makeTimeSequence(), fileName)
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

func (sb *SQLBundle) readConfig() error {
	if sb.Config == nil {
		config, err := ReadPackageJSON(sb.ConfigFile)
		if err != nil {
			return err
		}
		sb.Config = &config
	}
	return nil
}

func (sb *SQLBundle) Clean() error {
	err := os.RemoveAll(sb.BuildDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(sb.DepsDir)
	if err != nil {
		return err
	}
	return nil
}

func (sb *SQLBundle) Install() error {
	err := sb.readConfig()
	if err != nil {
		return err
	}

	err = os.RemoveAll(sb.DepsDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(sb.DepsDir, 0755)
	if err != nil {
		return err
	}

	if sb.Config.Dependencies == nil || len(sb.Config.Dependencies) == 0 {
		_ = os.RemoveAll(sb.DepsDir)
		return nil
	}

	for _, dep := range sb.Config.Dependencies {
		tarPath, err := downloadDependency(sb.DepsDir, dep)
		if err != nil {
			return err
		}
		_, tarFile := filepath.Split(tarPath)
		tarFile = strings.TrimSuffix(tarFile, filepath.Ext(tarFile))
		dest := filepath.Join(sb.DepsDir, tarFile)
		err = os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
		err = untarFile(tarPath, dest)
		if err != nil {
			return err
		}
		err = os.RemoveAll(tarPath)
		if err != nil {
			return err
		}
	}

	if ok, _ := isDirEmpty(sb.DepsDir); ok {
		_ = os.RemoveAll(sb.DepsDir)
	}
	return nil
}

func (sb *SQLBundle) Pack() error {
	err := sb.Install()
	if err != nil {
		return err
	}

	//create build dir
	err = os.RemoveAll(sb.BuildDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(sb.BuildDir, 0755)
	if err != nil {
		return err
	}

	packDirPath := filepath.Join(sb.BuildDir, "package")
	err = os.MkdirAll(packDirPath, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(packDirPath, "src"), 0755)
	if err != nil {
		return err
	}

	//write package json
	f, err := os.Create(filepath.Join(packDirPath, PACKAGE_JSON))
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("can not create file " + sb.ConfigFile)
	}
	defer func() {
		_ = f.Close()
	}()

	newConfig := *sb.Config
	newConfig.Version = sb.ReadVersion()
	bytes, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, string(bytes))
	//end write package json

	err = copyDirectory(sb.SourceDir, filepath.Join(packDirPath, "src"))
	if err != nil {
		return err
	}

	if exists(sb.DepsDir) {
		err = os.MkdirAll(filepath.Join(packDirPath, "deps"), 0755)
		if err != nil {
			return err
		}

		err = copyDirectory(sb.DepsDir, filepath.Join(packDirPath, "deps"))
		if err != nil {
			return err
		}

		if ok, _ := isDirEmpty(filepath.Join(packDirPath, "deps")); ok {
			_ = os.RemoveAll(filepath.Join(packDirPath, "deps"))
		}
	}

	if err != nil {
		return err
	}

	dest := filepath.Join(sb.BuildDir, fmt.Sprintf("%s-%s.tar", sb.Config.ArtifactId, sb.ReadVersion()))
	return tarFile(dest, []string{packDirPath})
}

func (sb *SQLBundle) Publish() error {
	if isEmpty(sb.Repository) || isEmpty(sb.Username) || isEmpty(sb.Password) {
		return errors.New("missing repository configuration")
	}

	err := sb.readConfig()
	if err != nil {
		return err
	}

	tarName := fmt.Sprintf("%s-%s.tar", sb.Config.ArtifactId, sb.ReadVersion())
	tarFile := filepath.Join(sb.BuildDir, tarName)
	_, err = os.Stat(tarFile)
	if err != nil {
		return err
	}
	args := strings.Split(sb.Config.GroupId, ".")
	args = append(args, sb.Config.ArtifactId)
	args = append(args, sb.ReadVersion())
	modulePath := strings.Join(args, "/")

	link := fmt.Sprintf("%s/%s/%s", sb.Repository, modulePath, tarName)
	return uploadFile(link, tarFile, sb.Username, sb.Password)
}
