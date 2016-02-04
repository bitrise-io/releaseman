package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"gopkg.in/yaml.v2"
)

const (
	repoCreateStr   = "repo create"
	currentStateStr = "current state"
	// DefaultConfigPth ...
	DefaultConfigPth = "./config.yml"
)

var (
	// IsCIMode ...
	IsCIMode = false
)

//=======================================
// Models
//=======================================

/*
development_branch: master
release_branch: master
start_state: ""
end_state: ""
changelog_path: "./changelog"
release_version: "1.1.0"
*/

// ReleaseConfig ...
type ReleaseConfig struct {
	DevelopmentBranch string `yaml:"development_branch"`
	ReleaseBranch     string `yaml:"release_branch"`
	StartState        string `yaml:"start_state"`
	EndState          string `yaml:"end_state"`
	ChangelogPath     string `yaml:"changelog_path"`
	ReleaseVersion    string `yaml:"release_version"`
}

// NewReleaseConfigFromFile ...
func NewReleaseConfigFromFile(pth string) (ReleaseConfig, error) {
	bytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return ReleaseConfig{}, err
	}

	config := ReleaseConfig{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return ReleaseConfig{}, err
	}

	log.Debugf("Config: %#v", config)

	return config, nil
}
