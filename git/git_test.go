package git

import (
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestVersionTaggedCommits(t *testing.T) {
	t.Log("Test valid sem-ver string")
	{
		_, err := version.NewVersion("1.0.0")
		require.Equal(t, nil, err)
	}

	t.Log("Test shorter sem-ver string")
	{
		_, err := version.NewVersion("1.0")
		require.Equal(t, nil, err)
	}

	t.Log("Test short sem-ver string")
	{
		_, err := version.NewVersion("1")
		require.Equal(t, nil, err)
	}
}

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
	/*
	   // commit b738dee2d32def019a4d553249004364046dc1bd
	   // commit: b738dee2d32def019a4d553249004364046dc1bd
	   // date: 1455631980
	   // author: Viktor Benei
	   // message: Merge branch 'master' of github.com:bitrise-tools/releaseman
	*/
	t.Log("Test valid commit")
	{
		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit: 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
date: 1455631980
author: Krisztián Gödrei
message: first change`

		commit, err := parseCommit(commitLine)
		require.Equal(t, nil, err)
		require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", commit.Hash)
		require.Equal(t, "Krisztián Gödrei", commit.Author)
		require.Equal(t, "first change", commit.Message)
	}

	t.Log("Test commit without hash")
	{
		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit:
date: 1455631980
author: Krisztián Gödrei
message: first change`

		_, err := parseCommit(commitLine)
		require.NotEqual(t, nil, err)
	}

	t.Log("Test commit without date")
	{
		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit: 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
date:
author: Krisztián Gödrei
message: first change`

		_, err := parseCommit(commitLine)
		require.NotEqual(t, nil, err)
	}

	t.Log("Test commit without author")
	{
		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit: 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
date: 1455631980
author:
message: first change`

		commit, err := parseCommit(commitLine)
		t.Logf("commit: %#v", commit)
		require.NotEqual(t, nil, err)
	}

	t.Log("Test commit without message")
	{
		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit: 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
date: 1455631980
author: Krisztián Gödrei
message: `

		commit, err := parseCommit(commitLine)
		require.Equal(t, nil, err)
		require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", commit.Hash)
		require.Equal(t, "Krisztián Gödrei", commit.Author)
	}

	t.Log("Test commit with multiline message")
	{
		multilineMessage := `multiline
test
commit`

		commitLine := `commit 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
commit: 85d8658733f73ae6d5407e8e4c2b81a5f2ed016c
date: 1455631980
author: Krisztián Gödrei
message: ` + multilineMessage

		commit, err := parseCommit(commitLine)
		require.Equal(t, nil, err)
		require.Equal(t, "85d8658733f73ae6d5407e8e4c2b81a5f2ed016c", commit.Hash)
		require.Equal(t, "Krisztián Gödrei", commit.Author)
		require.Equal(t, multilineMessage, commit.Message)
	}
}
