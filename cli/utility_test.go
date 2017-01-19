package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBumpVersion(t *testing.T) {
	ver, err := bumpedVersion("", -1)
	require.EqualError(t, err, "Invalid (negative) segment index: -1")
	require.Equal(t, "", ver)

	ver, err = bumpedVersion("", 0)
	require.EqualError(t, err, "Malformed version: ")
	require.Equal(t, "", ver)

	ver, err = bumpedVersion("-1", 0)
	require.EqualError(t, err, "Malformed version: -1")
	require.Equal(t, "", ver)

	ver, err = bumpedVersion("1.0.0", 12)
	require.EqualError(t, err, "Version does not have enough segments (segments count: 3) to increment segment at idx (12)")
	require.Equal(t, "", ver)

	ver, err = bumpedVersion("1", 0)
	require.NoError(t, err)
	require.Equal(t, "2.0.0", ver)

	ver, err = bumpedVersion("1.1", 0)
	require.NoError(t, err)
	require.Equal(t, "2.1.0", ver)

	ver, err = bumpedVersion("1.1", 1)
	require.NoError(t, err)
	require.Equal(t, "1.2.0", ver)

	ver, err = bumpedVersion("1.1.1", 0)
	require.NoError(t, err)
	require.Equal(t, "2.1.1", ver)

	ver, err = bumpedVersion("1.1.1", 1)
	require.NoError(t, err)
	require.Equal(t, "1.2.1", ver)

	ver, err = bumpedVersion("1.1.1", 2)
	require.NoError(t, err)
	require.Equal(t, "1.1.2", ver)
}
