package releaseman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReleaseConfigFromBytes(t *testing.T) {
	configStr := `
release:
  development_branch: develop
  release_branch: master
  version: 1.1.0
changelog:
  path: "./_changelog/changelog.md"
`
	config, err := NewConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.Equal(t, "develop", config.Release.DevelopmentBranch)
	require.Equal(t, "master", config.Release.ReleaseBranch)
	require.Equal(t, "1.1.0", config.Release.Version)

	require.Equal(t, "./_changelog/changelog.md", config.Changelog.Path)

	configStr = `
release:
  development_branch: develop
  release_branch: master
  version: 1.1.0
`
	config, err = NewConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.Equal(t, "develop", config.Release.DevelopmentBranch)
	require.Equal(t, "master", config.Release.ReleaseBranch)
	require.Equal(t, "1.1.0", config.Release.Version)

	configStr = `
changelog:
  path: "./_changelog/changelog.md"
`

	config, err = NewConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.Equal(t, "./_changelog/changelog.md", config.Changelog.Path)
}
