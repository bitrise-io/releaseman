package releaseman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReleaseConfigFromBytes(t *testing.T) {
	configStr := `
config:
  release:
    development_branch: master
    release_branch: master
    version: "1.1.0"
  changelog:
    changelog_path: "./changelog"
`

	config, err := NewConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.NotEqual(t, nil, config.Release)
	require.Equal(t, "develop", config.Release.DevelopmentBranch)
	require.Equal(t, "master", config.Release.ReleaseBranch)
	require.Equal(t, "1.1.0", config.Release.Version)

	require.NotEqual(t, nil, config.Changelog)
	require.Equal(t, "./changelog", config.Changelog.Path)

	configStr = `
config:
  release:
    development_branch: master
    release_branch: master
    version: "1.1.0"
`

	config, err = NewConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.NotEqual(t, nil, config.Release)
	require.Equal(t, "develop", config.Release.DevelopmentBranch)
	require.Equal(t, "master", config.Release.ReleaseBranch)
	require.Equal(t, "1.1.0", config.Release.Version)

	require.Equal(t, nil, config.Changelog)

	configStr = `
config:
changelog:
  changelog_path: "./changelog"
`

	config, err = NewConfigFromBytes([]byte(configStr))
	require.NotEqual(t, nil, err)
}
