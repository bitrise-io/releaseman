package configs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReleaseConfigFromBytes(t *testing.T) {
	configStr := `
development_branch: develop
release_branch: master
start_state: "initial commit"
end_state: "current state"
changelog_path: "./changelog"
release_version: "1.1.0"
`

	config, err := NewReleaseConfigFromBytes([]byte(configStr))
	require.Equal(t, nil, err)

	require.Equal(t, "develop", config.DevelopmentBranch)
	require.Equal(t, "master", config.ReleaseBranch)
	require.Equal(t, "initial commit", config.StartState)
	require.Equal(t, "current state", config.EndState)
	require.Equal(t, "./changelog", config.ChangelogPath)
	require.Equal(t, "1.1.0", config.ReleaseVersion)
}
