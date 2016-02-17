package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBumpVersion(t *testing.T) {
	version, err := bumpedVersion("", 0)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", version)

	version, err = bumpedVersion("-1", 0)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", version)

	version, err = bumpedVersion("1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.0.0", version)

	version, err = bumpedVersion("1.1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.1.0", version)

	version, err = bumpedVersion("1.1", 1)
	require.Equal(t, nil, err)
	require.Equal(t, "1.2.0", version)

	version, err = bumpedVersion("1.1.1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.1.1", version)

	version, err = bumpedVersion("1.1.1", 1)
	require.Equal(t, nil, err)
	require.Equal(t, "1.2.1", version)

	version, err = bumpedVersion("1.1.1", 2)
	require.Equal(t, nil, err)
	require.Equal(t, "1.1.2", version)
}
