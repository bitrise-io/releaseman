package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	unixTimestampStr := "1454498673"
	unixTime, err := parseDate(unixTimestampStr)
	require.Equal(t, nil, err)
	require.Equal(t, time.Unix(1454498673, 0), unixTime)

	unixTimestampStr = ""
	unixTime, err = parseDate(unixTimestampStr)
	require.NotEqual(t, nil, err)
	require.Equal(t, time.Time{}, unixTime)

	unixTimestampStr = "-1"
	unixTime, err = parseDate(unixTimestampStr)
	require.NotEqual(t, nil, err)
	require.Equal(t, time.Time{}, unixTime)

	unixTimestampStr = "abc"
	unixTime, err = parseDate(unixTimestampStr)
	require.NotEqual(t, nil, err)
	require.Equal(t, time.Time{}, unixTime)
}

func TestParseCommit(t *testing.T) {
	commitLine := "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c 1454498673 (Krisztián Gödrei) first change"
	commit, err := parseCommit(commitLine)
	require.Equal(t, nil, err)
	require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", commit.Hash)
	require.Equal(t, "Krisztián Gödrei", commit.Author)
	require.Equal(t, "first change", commit.Message)

	commitLine = "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c 1454498673 () first change"
	commit, err = parseCommit(commitLine)
	require.Equal(t, nil, err)
	require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", commit.Hash)
	require.Equal(t, "", commit.Author)
	require.Equal(t, "first change", commit.Message)

	commitLine = "1454498673 first change"
	commit, err = parseCommit(commitLine)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", commit.Hash)
	require.Equal(t, "", commit.Author)
	require.Equal(t, "", commit.Message)

	commitLine = "first change"
	commit, err = parseCommit(commitLine)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", commit.Hash)
	require.Equal(t, "", commit.Author)
	require.Equal(t, "", commit.Message)

	commitLine = ""
	commit, err = parseCommit(commitLine)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", commit.Hash)
	require.Equal(t, "", commit.Author)
	require.Equal(t, "", commit.Message)
}
