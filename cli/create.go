package cli

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/configs"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/codegangsta/cli"
)

//=======================================
// Utility
//=======================================

func printRollBackMessage() {
	fmt.Println()
	log.Infoln("How to roll-back?")
	log.Infoln("* if you want to undo the last commit you can call:")
	log.Infoln("    $ git reset --hard HEAD~1")
	log.Infoln("* to roll back to the remote state:")
	log.Infoln("    $ git reset --hard origin/[branch-name]")
	fmt.Println()
}

func printDoneMessage(config configs.Config) {
	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s released 🚀", config.Release.Version))
	log.Infoln("Take a look at your git, and if you are happy with the release, push the changes.")
}

func askForDevelopmentBranch() (string, error) {
	branches, err := git.LocalBranches()
	if err != nil {
		return "", err
	}

	fmt.Println()
	developmentBranch, err := goinp.SelectFromStrings("Select your development branch!", branches)
	if err != nil {
		return "", err
	}

	// 'git branch --list' marks the current branch with (* )
	if strings.HasPrefix(developmentBranch, "* ") {
		developmentBranch = strings.TrimPrefix(developmentBranch, "* ")
	}
	return developmentBranch, nil
}

func askForReleaseBranch() (string, error) {
	branches, err := git.LocalBranches()
	if err != nil {
		return "", err
	}

	fmt.Println()
	releaseBranch, err := goinp.SelectFromStrings("Select your release branch!", branches)
	if err != nil {
		return "", err
	}

	// 'git branch --list' marks the current branch with (* )
	if strings.HasPrefix(releaseBranch, "* ") {
		releaseBranch = strings.TrimPrefix(releaseBranch, "* ")
	}

	return releaseBranch, nil
}

func askForReleaseVersion() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in release version!")
}

func askForChangelogPath() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in changelog path!")
}

func firstCommitAfterTag(taggedCommit git.CommitModel, commits []git.CommitModel) (git.CommitModel, bool) {
	for _, commit := range commits {
		if taggedCommit.Date.Sub(commit.Date) < 0 {
			return commit, true
		}
	}
	return git.CommitModel{}, false
}

func sectionHeader(startTag, endTag string) string {
	return fmt.Sprintf("%s - %s\n", startTag, endTag)
}

func sectionBody(commits []git.CommitModel) string {
	body := ""
	for _, commit := range commits {
		body += fmt.Sprintf(" * %s\n", commit.Message)
	}
	body += "\n"
	return body
}

func reverse(commits []git.CommitModel) []git.CommitModel {
	reversed := []git.CommitModel{}
	for i := len(commits) - 1; i >= 0; i-- {
		reversed = append(reversed, commits[i])
	}
	return reversed
}

func commitsBetween(startDate *time.Time, endDate *time.Time, commits []git.CommitModel) []git.CommitModel {
	relevantCommits := []git.CommitModel{}
	isRelevantCommit := false

	for _, commit := range commits {
		if !isRelevantCommit && (startDate == nil || (*startDate).Sub(commit.Date) <= 0) {
			isRelevantCommit = true
		}

		if isRelevantCommit && (endDate == nil || (*endDate).Sub(commit.Date) <= 0) {
			return relevantCommits
		}

		if isRelevantCommit {
			relevantCommits = append(relevantCommits, commit)
		}
	}

	return reverse(relevantCommits)
}

func writeChnagelog(changelogPath string, commits, taggedCommits []git.CommitModel, nextVersion string) error {
	changelog := "\n"

	if len(taggedCommits) > 0 {
		// Commits between initial commit and first tag
		relevantCommits := commitsBetween(nil, &(taggedCommits[0].Date), commits)
		header := sectionHeader("", taggedCommits[0].Tag)
		body := sectionBody(relevantCommits)
		if body != "\n" {
			changelog = header + body + changelog
		}

		if len(taggedCommits) > 1 {
			// Commits between tags
			for i := 0; i < len(taggedCommits)-1; i++ {
				startTaggedCommit := taggedCommits[i]
				endTaggedCommit := taggedCommits[i+1]

				relevantCommits = commitsBetween(&(startTaggedCommit.Date), &(endTaggedCommit.Date), commits)
				header := sectionHeader(startTaggedCommit.Tag, endTaggedCommit.Tag)
				body := sectionBody(relevantCommits)
				if body != "\n" {
					changelog = header + body + changelog
				}
			}
		}

		// Commits between last tag and current state
		relevantCommits = commitsBetween(&(taggedCommits[len(taggedCommits)-1].Date), nil, commits)
		header = sectionHeader(taggedCommits[len(taggedCommits)-1].Tag, "")
		body = sectionBody(relevantCommits)
		if body != "\n" {
			changelog = header + body + changelog
		}
	} else {
		relevantCommits := commitsBetween(nil, nil, commits)
		header := sectionHeader("", "")
		body := sectionBody(relevantCommits)
		if body != "\n" {
			changelog = header + body + changelog
		}
	}

	return fileutil.WriteStringToFile(changelogPath, changelog)
}

//=======================================
// Main
//=======================================

func create(c *cli.Context) {
	//
	// Fail if git is not clean
	if areChanges, err := git.AreUncommitedChanges(); err != nil {
		log.Fatalf("Failed to get uncommited changes, error: %#v", err)
	} else if areChanges {
		log.Fatalf("There are uncommited changes in your git, please commit your changes before continue release!")
	}

	printRollBackMessage()

	//
	// Build config from file
	config := configs.Config{}
	configPath := ""
	if c.IsSet("config") {
		configPath = c.String("config")
	} else {
		configPath = configs.DefaultConfigPth
	}

	if exist, err := pathutil.IsPathExists(configPath); err != nil {
		log.Warnf("Failed to check if path exist, error: %#v", err)
	} else if exist {
		config, err = configs.NewConfigFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to parse release config at (%s), error: %#v", configPath, err)
		}
	}

	var err error
	if c.IsSet(DevelopmentBranchKey) {
		config.Release.DevelopmentBranch = c.String(DevelopmentBranchKey)
	}
	if config.Release.DevelopmentBranch == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: development branch")
		} else {
			config.Release.DevelopmentBranch, err = askForDevelopmentBranch()
			if err != nil {
				log.Fatalf("Failed to ask for development branch, error: %s", err)
			}
		}
	}

	//
	// Checkout the release start branch
	currentBranch, err := git.CurrentBranchName()
	if err != nil {
		log.Fatalf("Failed to get current branch name, error: %#v", err)
	}

	if config.Release.DevelopmentBranch != currentBranch {
		log.Warnf("Your current branch (%s), should be the development branch (%s)!", currentBranch, config.Release.DevelopmentBranch)

		fmt.Println()
		checkout, err := goinp.AskForBool(fmt.Sprintf("Would you like to checkout development branch (%s)?", config.Release.DevelopmentBranch))
		if err != nil {
			log.Fatalf("Failed to ask for checkout, error: %#v", err)
		}

		if !checkout {
			log.Fatalf("Current branch should be the development branch (%s)!", config.Release.DevelopmentBranch)
		}

		if err := git.CheckoutBranch(config.Release.DevelopmentBranch); err != nil {
			log.Fatalf("Failed to checkout branch (%s), error: %#v", config.Release.DevelopmentBranch, err)
		}
	}
	//
	//

	if c.IsSet(ReleaseBranchKey) {
		config.Release.ReleaseBranch = c.String(ReleaseBranchKey)
	}
	if config.Release.ReleaseBranch == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: release branch")
		} else {
			config.Release.ReleaseBranch, err = askForReleaseBranch()
			if err != nil {
				log.Fatalf("Failed to ask for release branch, error: %s", err)
			}
		}
	}

	if c.IsSet(ReleaseVersionKey) {
		config.Release.Version = c.String(ReleaseVersionKey)
	}
	if config.Release.Version == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: release version")
		} else {
			config.Release.Version, err = askForReleaseVersion()
			if err != nil {
				log.Fatalf("Failed to ask for release version, error: %s", err)
			}
		}
	}

	if c.IsSet(ChangelogPathKey) {
		config.Changelog.Path = c.String(ChangelogPathKey)
	}
	if config.Changelog.Path == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: changelog path")
		} else {
			config.Changelog.Path, err = askForChangelogPath()
			if err != nil {
				log.Fatalf("Failed to ask for changelog path, error: %s", err)
			}
		}
	}

	//
	// Print config
	fmt.Println()
	log.Infof("Your config:")
	log.Infof(" * Development branch: %s", config.Release.DevelopmentBranch)
	log.Infof(" * Release branch: %s", config.Release.ReleaseBranch)
	log.Infof(" * Release version: %s", config.Release.Version)
	log.Infof(" * Changelog path: %s", config.Changelog.Path)
	fmt.Println()

	if !configs.IsCIMode {
		ok, err := goinp.AskForBool("Are you ready for release?")
		if err != nil {
			log.Fatalf("Failed to ask for input, error: %s", err)
		}
		if !ok {
			log.Fatal("Aborted release")
		}
	}

	//
	// Generate changelog and release
	startCommit, err := git.FirstCommit()
	if err != nil {
		log.Fatalf("Failed to get first commit, error: %#v", err)
	}

	endCommit, err := git.LatestCommit()
	if err != nil {
		log.Fatalf("Failed to get latest commit, error: %#v", err)
	}

	taggedCommits, err := git.ListTaggedCommits()
	if err != nil {
		log.Fatalf("Failed to get tagged commits, error: %#v", err)
	}

	startDate := startCommit.Date
	endDate := endCommit.Date
	relevantTags := taggedCommits

	if config.Changelog.Path != "" {
		if exist, err := pathutil.IsPathExists(config.Changelog.Path); err != nil {
			log.Fatalf("Failed to check if path exist, error: %#v", err)
		} else if exist {
			if len(taggedCommits) > 0 {
				lastTaggedCommit := taggedCommits[len(taggedCommits)-1]
				startDate = lastTaggedCommit.Date
				relevantTags = []git.CommitModel{lastTaggedCommit}
			}
		}
	}

	fmt.Println()
	log.Infof("Collect commits between (%s - %s)", startDate, endDate)

	fmt.Println()
	log.Infof("=> Generating changelog...")
	commits, err := git.GetCommitsBetween(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get commits, error: %#v", err)
	}
	if err := writeChnagelog(config.Changelog.Path, commits, relevantTags, config.Release.Version); err != nil {
		log.Fatalf("Failed to write changelog, error: %#v", err)
	}

	fmt.Println()
	log.Infof("=> Adding changes to git...")
	if err := git.Add([]string{config.Changelog.Path}); err != nil {
		log.Fatalf("Failed to git add, error: %s", err)
	}

	if err := git.Commit(fmt.Sprintf("v%s", config.Release.Version)); err != nil {
		log.Fatalf("Failed to git commit, error: %s", err)
	}

	fmt.Println()
	log.Infof("=> Merging changes into release branch...")
	if err := git.CheckoutBranch(config.Release.ReleaseBranch); err != nil {
		log.Fatalf("Failed to git checkout, error: %s", err)
	}

	mergeCommitMessage := fmt.Sprintf("Merge %s into %s, release: v%s", config.Release.DevelopmentBranch, config.Release.ReleaseBranch, config.Release.Version)
	if err := git.Merge(config.Release.DevelopmentBranch, mergeCommitMessage); err != nil {
		log.Fatalf("Failed to git merge, error: %s", err)
	}

	fmt.Println()
	log.Infof("=> Tagging release branch...")
	if err := git.Tag(config.Release.Version); err != nil {
		log.Fatalf("Failed to git tag, error: %s", err)
	}

	if err := git.CheckoutBranch(config.Release.DevelopmentBranch); err != nil {
		log.Fatalf("Failed to git checkout, error: %s", err)
	}

	printDoneMessage(config)
}
