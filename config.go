package sqlbundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const PACKAGE_JSON = "package.json"

type PackageJSON struct {
	GroupId      string   `json:"group"`
	ArtifactId   string   `json:"artifact"`
	Version      string   `json:"version"`
	Dependencies []string `json:"dependencies"`
}

func ReadPackageJSON(file string) (config PackageJSON, err error) {
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.New(file + " file not found")
		return
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		err = errors.New(fmt.Sprintf("read application config file get error %v", err))
		return
	}
	err = json.Unmarshal(bytes, &config)
	return
}
