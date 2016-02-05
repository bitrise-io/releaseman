package git

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

//=======================================
// Models
//=======================================

// CommitModel ...
type CommitModel struct {
	Hash    string
	Message string
	Date    time.Time
	Author  string
	Tag     string
}

//=======================================
// Utility
//=======================================

func parseDate(unixTimeStampStr string) (time.Time, error) {
	i, err := strconv.ParseInt(unixTimeStampStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)

	return tm, nil
}

func parseCommit(commitLineStr string) (CommitModel, error) {
	// ba58d366e3565a0f52250dce992fe29c29750f79 1454582002 go tests added
	re := regexp.MustCompile(`(?P<hash>[0-9a-z]+) (?P<date>[0-9]+) (?P<author>.*) (?P<message>.+)`)
	results := re.FindAllStringSubmatch(commitLineStr, -1)

	for _, v := range results {
		hash := v[1]
		date, err := parseDate(v[2])
		if err != nil {
			return CommitModel{}, err
		}
		author := v[3]
		message := v[4]

		return CommitModel{
			Hash:    hash,
			Message: message,
			Date:    date,
			Author:  author,
		}, nil
	}
	return CommitModel{}, fmt.Errorf("Failed to parse commit: %s", commitLineStr)
}

//=======================================
// Git functions
//=======================================

// LocalBranches ...
func LocalBranches() ([]string, error) {
	out, err := NewPrintableCommand("git", "branch", "--list").Run()
	if err != nil {
		return []string{}, err
	}
	return splitByNewLineAndStrip(out), nil
}

// TaggedCommits ...
func TaggedCommits() ([]CommitModel, error) {
	out, err := NewPrintableCommand("git", "tag", "--list").Run()
	if err != nil {
		return []CommitModel{}, err
	}
	taggedCommits := []CommitModel{}
	tags := splitByNewLineAndStrip(out)
	for _, tag := range tags {
		out, err = NewPrintableCommand("git", "rev-list", "-n", "1", `--pretty=format:%H %ct %an %s`, tag).Run()
		if err != nil {
			return []CommitModel{}, err
		}

		commit, err := parseCommit(strip(out))
		if err != nil {
			return []CommitModel{}, fmt.Errorf("Failed to parse commit: %#v", err)
		}
		commit.Tag = tag

		taggedCommits = append(taggedCommits, commit)
	}
	return taggedCommits, nil
}

// CurrentBranchName ...
func CurrentBranchName() (string, error) {
	out, err := NewPrintableCommand("git", "symbolic-ref", "--short", "HEAD").Run()
	if err != nil {
		return "", err
	}
	return strip(out), nil
}

// AreUncommitedChanges ...
func AreUncommitedChanges() (bool, error) {
	out, err := NewPrintableCommand("git", "status", "--porcelain").Run()
	if err != nil {
		return false, err
	}
	return (out != ""), nil
}

// CheckoutBranch ...
func CheckoutBranch(branch string) error {
	if _, err := NewPrintableCommand("git", "checkout", branch).Run(); err != nil {
		return err
	}
	return nil
}

// FirstCommit ...
func FirstCommit() (CommitModel, error) {
	out, err := NewPrintableCommand("git", "rev-list", "--max-parents=0", `--pretty=format:%H %ct %an %s`, "HEAD").Run()
	if err != nil {
		return CommitModel{}, err
	}
	commit, err := parseCommit(strip(out))
	if err != nil {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: %#v", err)
	}
	return commit, nil
}

// LatestCommit ...
func LatestCommit() (CommitModel, error) {
	out, err := NewPrintableCommand("git", "log", "-1", `--pretty=format:%H %ct %an %s`).Run()
	if err != nil {
		return CommitModel{}, err
	}
	commit, err := parseCommit(strip(out))
	if err != nil {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: %#v", err)
	}
	return commit, nil
}

// CommitOfTag ...
func CommitOfTag(tag string) (CommitModel, error) {
	out, err := NewPrintableCommand("git", "show-ref", tag).Run()
	if err != nil {
		return CommitModel{}, err
	}
	commit, err := parseCommit(out)
	if err != nil {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: %#v", err)
	}
	return commit, nil
}

// GetCommitsBetween ...
func GetCommitsBetween(startDate, endDate time.Time) ([]CommitModel, error) {
	out, err := NewPrintableCommand("git", "log", `--pretty=format:%H %ct %an %s`, "--reverse").Run()
	if err != nil {
		return []CommitModel{}, err
	}
	commitList := splitByNewLineAndStrip(out)

	commits := []CommitModel{}
	isRelevantCommit := false

	for _, commitListItem := range commitList {
		commit, err := parseCommit(commitListItem)
		if err != nil {
			return []CommitModel{}, fmt.Errorf("Failed to parse commit, error: %#v", err)
		}

		if !isRelevantCommit && startDate.Sub(commit.Date) <= 0 {
			isRelevantCommit = true
		}

		if isRelevantCommit {
			commits = append(commits, commit)
		}

		if isRelevantCommit && endDate.Sub(commit.Date) <= 0 {
			break
		}
	}
	return commits, nil
}

// Add ...
func Add(files []string) error {
	for _, file := range files {
		if _, err := NewPrintableCommand("git", "add", file).Run(); err != nil {
			return err
		}
	}
	return nil
}

// Commit ...
func Commit(message string) error {
	if _, err := NewPrintableCommand("git", "commit", "-m", message).Run(); err != nil {
		return err
	}
	return nil
}

// Merge ...
func Merge(branch, commitMessage string) error {
	if _, err := NewPrintableCommand("git", "merge", branch, "--no-ff", "-m", commitMessage).Run(); err != nil {
		return err
	}
	return nil
}

// Tag ...
func Tag(version string) error {
	if _, err := NewPrintableCommand("git", "tag", version).Run(); err != nil {
		return err
	}
	return nil
}
