package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommit(t *testing.T) {
	commitLine := "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c first change"
	hash, message, ok := parseCommit(commitLine)
	require.Equal(t, true, ok)
	require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", hash)
	require.Equal(t, "first change", message)

	commitLine = "first change"
	hash, message, ok = parseCommit(commitLine)
	require.Equal(t, true, ok)
	require.Equal(t, "first", hash)
	require.Equal(t, "change", message)

	commitLine = "85D8658733F73AE6d5407e8e4c2b81a5f2ed016c first change"
	hash, message, ok = parseCommit(commitLine)
	require.Equal(t, false, ok)

	commitLine = "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c"
	hash, message, ok = parseCommit(commitLine)
	require.Equal(t, false, ok)

	commitLine = ""
	hash, message, ok = parseCommit(commitLine)
	require.Equal(t, false, ok)
}
