package sqlbundle

const SQL_BUNDLE_CONFIG = "package.json"

type SQLBundleConfig struct {
	GroupId      string   `json:"group"`
	ArtifactId   string   `json:"artifact"`
	Version      string   `json:"version"`
	Dependencies []string `json:"dependencies"`
}

func ReadSQLBundleConfig(file string) (SQLBundleConfig, error) {
	return SQLBundleConfig{}, nil
}
