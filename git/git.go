package git

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	version "github.com/hashicorp/go-version"
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
	if i < 0 {
		return time.Time{}, fmt.Errorf("Invalid time stamp (%s)", unixTimeStampStr)
	}
	tm := time.Unix(i, 0)

	return tm, nil
}

func parseCommit(commitLineStr string) (CommitModel, error) {
	// commit b738dee2d32def019a4d553249004364046dc1bd
	// commit: b738dee2d32def019a4d553249004364046dc1bd
	// date: 1455631980
	// author: Viktor Benei
	// message: Merge branch 'master' of github.com:bitrise-tools/releaseman
	hashPrefix := "commit: "
	datePrefix := "date: "
	authorPrefix := "author: "
	messagePrefix := "message: "

	hash := ""
	dateStr := ""
	author := ""
	message := ""

	commitSplits := splitByNewLineAndStrip(commitLineStr)
	if len(commitSplits) < 5 {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: (%s)", commitLineStr)
	}

	messageStart := false
	for _, line := range commitSplits {
		if strings.HasPrefix(line, hashPrefix) {
			hash = strings.TrimPrefix(line, hashPrefix)
		} else if strings.HasPrefix(line, datePrefix) {
			dateStr = strings.TrimPrefix(line, datePrefix)
		} else if strings.HasPrefix(line, authorPrefix) {
			author = strings.TrimPrefix(line, authorPrefix)
		} else if strings.HasPrefix(line, messagePrefix) {
			messageStart = true
		}

		if messageStart {
			if strings.HasPrefix(line, messagePrefix) {
				message += strings.TrimPrefix(line, messagePrefix)
			} else {
				message += fmt.Sprintf("\n%s", line)
			}
		}
	}

	if hash == "" || dateStr == "" || author == "" {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: (%s)", commitLineStr)
	}

	date, err := parseDate(dateStr)
	if err != nil {
		return CommitModel{}, err
	}

	return CommitModel{
		Hash:    hash,
		Message: message,
		Date:    date,
		Author:  author,
	}, nil
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

// VersionTaggedCommits ...
func VersionTaggedCommits() ([]CommitModel, error) {
	out, err := NewPrintableCommand("git", "tag", "--list").Run()
	if err != nil {
		return []CommitModel{}, err
	}
	taggedCommits := []CommitModel{}
	tags := splitByNewLineAndStrip(out)
	for _, tag := range tags {
		// is tag sem-ver tag?
		_, err := version.NewVersion(tag)
		if err != nil {
			continue
		}

		out, err = NewPrintableCommand("git", "rev-list", "-n", "1", `--pretty=format:commit: %H%ndate: %ct%nauthor: %an%nmessage: %s`, tag).Run()
		if err != nil {
			return []CommitModel{}, err
		}

		commit, err := parseCommit(Strip(out))
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
	return Strip(out), nil

}

// AreUncommitedChanges ...
func AreUncommitedChanges() (bool, error) {
	out, err := NewPrintableCommand("git", "status", "--porcelain").Run()
	if err != nil {
		return false, err
	}
	return (out != ""), nil
}

// GetChangedFiles ...
func GetChangedFiles() ([]string, error) {
	out, err := NewPrintableCommand("git", "status", "--porcelain").Run()
	if err != nil {
		return []string{}, err
	}

	changes := []string{}
	changeList := splitByNewLineAndStrip(out)
	for _, change := range changeList {
		changeSplits := strings.Split(change, " ")

		normalizedChangeSplits := changeSplits[1:len(changeSplits)]
		normalizedChangeStr := strings.Join(normalizedChangeSplits, " ")

		changes = append(changes, normalizedChangeStr)
	}

	return changes, nil
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
	out, err := NewPrintableCommand("git", "rev-list", "--max-parents=0", `--pretty=format:commit: %H%ndate: %ct%nauthor: %an%nmessage: %s`, "HEAD").Run()
	if err != nil {
		return CommitModel{}, err
	}
	commit, err := parseCommit(Strip(out))
	if err != nil {
		return CommitModel{}, fmt.Errorf("Failed to parse commit: %#v", err)
	}
	return commit, nil
}

// LatestCommit ...
func LatestCommit() (CommitModel, error) {
	out, err := NewPrintableCommand("git", "log", "-1", `--pretty=format:commit: %H%ndate: %ct%nauthor: %an%nmessage: %s`).Run()
	if err != nil {
		return CommitModel{}, err
	}
	commit, err := parseCommit(Strip(out))
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
	out, err := NewPrintableCommand("git", "log", `--pretty=format:commit: %H%ndate: %ct%nauthor: %an%nmessage: %s`, "--reverse").Run()
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
