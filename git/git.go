package git

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

//=======================================
// Utility
//=======================================

func parseCommit(commitLineStr string) (string, string, bool) {
	re := regexp.MustCompile(`(?P<commit>[0-9a-z]*)\s(?P<message>.*)`)
	results := re.FindAllStringSubmatch(commitLineStr, -1)

	for _, v := range results {
		commitHash := v[1]
		commitMessage := v[2]

		if commitHash != "" && commitMessage != "" {
			return commitHash, commitMessage, true
		}
	}
	return "", "", false
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
	return splitByNewLine(out), nil
}

// ListTags ...
func ListTags() ([]string, error) {
	out, err := NewPrintableCommand("git", "tag", "--list").Run()
	if err != nil {
		return []string{}, err
	}
	return splitByNewLine(out), nil
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
func FirstCommit() (string, error) {
	out, err := NewPrintableCommand("git", "rev-list", "--max-parents=0", "HEAD").Run()
	if err != nil {
		return "", err
	}
	return strip(out), nil
}

// LatestCommit ...
func LatestCommit() (string, error) {
	out, err := NewPrintableCommand("git", "rev-parse", "HEAD").Run()
	if err != nil {
		return "", err
	}
	return strip(out), nil
}

// CommitHashOfTag ...
func CommitHashOfTag(tag string) (string, error) {
	out, err := NewPrintableCommand("git", "show-ref", tag).Run()
	if err != nil {
		return "", err
	}

	hash, _, success := parseCommit(out)
	if !success {
		return "", fmt.Errorf("Failed to parse commit: %s", out)
	}
	return hash, nil
}

// CommitMessages ...
func CommitMessages(startCommit, endCommit string) ([]map[string]string, error) {
	out, err := NewPrintableCommand("git", "log", "--pretty=oneline", "--reverse").Run()
	if err != nil {
		return []map[string]string{}, err
	}

	commitList := splitByNewLine(out)

	commitMapList := []map[string]string{}
	isRelevantCommit := (startCommit == "") // collecting from repo init if no start commit

	for _, commit := range commitList {
		commitHash, commitMessage, success := parseCommit(commit)
		if !success {
			log.Warningf("Failed to parse commit line (%s)", commit)
			continue
		}

		if !isRelevantCommit && startCommit == commitHash {
			isRelevantCommit = true
		}

		if isRelevantCommit {
			commitMapList = append(commitMapList, map[string]string{commitHash: commitMessage})
		}

		if isRelevantCommit && endCommit != "" && endCommit == commitHash {
			break
		}
	}
	return commitMapList, nil
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
