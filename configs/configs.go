package configs

import (
	"fmt"

	"github.com/bitrise-io/go-utils/fileutil"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultConfigPth ...
	DefaultConfigPth = "./release_config.yml"
	//InitialCommitStr ...
	InitialCommitStr = "initial commit"
	// CurrentStateStr ...
	CurrentStateStr = "current state"
)

var (
	// IsCIMode ...
	IsCIMode = false
)

//=======================================
// Models
//=======================================

// Release ...
type Release struct {
	DevelopmentBranch string `yaml:"development_branch"`
	ReleaseBranch     string `yaml:"release_branch"`
	Version           string `yaml:"version"`
}

// Changelog ...
type Changelog struct {
	Path string `yaml:"path"`
}

// Config ...
type Config struct {
	Release   Release   `yaml:"release,omitempty"`
	Changelog Changelog `yaml:"changelog,omitempty"`
}

// NewConfigFromFile ...
func NewConfigFromFile(pth string) (Config, error) {
	bytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return Config{}, err
	}
	return NewConfigFromBytes(bytes)
}

// NewConfigFromBytes ...
func NewConfigFromBytes(bytes []byte) (Config, error) {
	type FileConfig struct {
		Release   *Release   `yaml:"release,omitempty"`
		Changelog *Changelog `yaml:"changelog,omitempty"`
	}

	fileConfig := FileConfig{}
	if err := yaml.Unmarshal(bytes, &fileConfig); err != nil {
		return Config{}, err
	}

	fmt.Printf("%#v\n", fileConfig)

	if fileConfig.Release == nil {
		return Config{}, fmt.Errorf("Invalid configuration: no release configuration defined")
	}

	config := Config{}
	config.Release = *fileConfig.Release
	if fileConfig.Changelog != nil {
		config.Changelog = *fileConfig.Changelog
	}

	return config, nil
}
