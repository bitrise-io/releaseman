package releaseman

import (
	"testing"
	"time"

	"github.com/bitrise-tools/releaseman/git"
	"github.com/stretchr/testify/require"
)

func TestCommitsBetween(t *testing.T) {
	unixTimestampStart := time.Unix(1454498673, 0)
	unixTimestampEnd := time.Unix(1454498683, 0)
	unixTimestampBeforeEnd := time.Unix(1454498682, 0)

	allCommits := []git.CommitModel{
		git.CommitModel{
			Date: time.Unix(1454498663, 0),
			Hash: "1",
		},
		git.CommitModel{
			Date: unixTimestampStart,
			Hash: "2",
		},
		git.CommitModel{
			Date: unixTimestampBeforeEnd,
			Hash: "3",
		},
		git.CommitModel{
			Date: time.Unix(1454498693, 0),
			Hash: "4",
		},
	}

	commits := commitsBetween(nil, nil, allCommits)
	require.Equal(t, 4, len(commits))
	require.Equal(t, "1", commits[0].Hash)
	require.Equal(t, "2", commits[1].Hash)
	require.Equal(t, "3", commits[2].Hash)
	require.Equal(t, "4", commits[3].Hash)

	commits = commitsBetween(nil, &unixTimestampEnd, allCommits)
	require.Equal(t, 3, len(commits))
	require.Equal(t, "1", commits[0].Hash)
	require.Equal(t, "2", commits[1].Hash)
	require.Equal(t, "3", commits[2].Hash)

	commits = commitsBetween(&unixTimestampStart, nil, allCommits)
	require.Equal(t, 3, len(commits))
	require.Equal(t, "2", commits[0].Hash)
	require.Equal(t, "3", commits[1].Hash)
	require.Equal(t, "4", commits[2].Hash)

	commits = commitsBetween(&unixTimestampStart, &unixTimestampEnd, allCommits)
	require.Equal(t, 2, len(commits))
	require.Equal(t, "2", commits[0].Hash)
	require.Equal(t, "3", commits[1].Hash)
}

func TestReversedSections(t *testing.T) {
	sections := []ChangelogSectionModel{}

	reversed := reversedSections(sections)
	require.Equal(t, 0, len(reversed))

	sections = []ChangelogSectionModel{
		ChangelogSectionModel{
			StartTaggedCommit: git.CommitModel{
				Tag: "1",
			},
		},
		ChangelogSectionModel{
			StartTaggedCommit: git.CommitModel{
				Tag: "2",
			},
		},
		ChangelogSectionModel{
			StartTaggedCommit: git.CommitModel{
				Tag: "3",
			},
		},
	}

	reversed = reversedSections(sections)
	require.Equal(t, "3", reversed[0].StartTaggedCommit.Tag)
	require.Equal(t, "2", reversed[1].StartTaggedCommit.Tag)
	require.Equal(t, "1", reversed[2].StartTaggedCommit.Tag)
}

func TestBumpVersion(t *testing.T) {
	bumpedVersion, err := BumpedVersion("", 0)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", bumpedVersion)

	bumpedVersion, err = BumpedVersion("-1", 0)
	require.NotEqual(t, nil, err)
	require.Equal(t, "", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.0.0", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1.1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.1.0", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1.1", 1)
	require.Equal(t, nil, err)
	require.Equal(t, "1.2.0", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1.1.1", 0)
	require.Equal(t, nil, err)
	require.Equal(t, "2.1.1", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1.1.1", 1)
	require.Equal(t, nil, err)
	require.Equal(t, "1.2.1", bumpedVersion)

	bumpedVersion, err = BumpedVersion("1.1.1", 2)
	require.Equal(t, nil, err)
	require.Equal(t, "1.1.2", bumpedVersion)
}
