package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/bitrise-tools/releaseman/releaseman"
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

func printDoneMessage(config releaseman.Config) {
	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s released 🚀", config.Release.Version))
	log.Infoln("Take a look at your git, and if you are happy with the release, push the changes.")
}

func firstCommitAfterTag(taggedCommit git.CommitModel, commits []git.CommitModel) (git.CommitModel, bool) {
	for _, commit := range commits {
		if taggedCommit.Date.Sub(commit.Date) < 0 {
			return commit, true
		}
	}
	return git.CommitModel{}, false
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
	config := releaseman.Config{}
	configPath := ""
	if c.IsSet("config") {
		configPath = c.String("config")
	} else {
		configPath = releaseman.DefaultConfigPth
	}

	if exist, err := pathutil.IsPathExists(configPath); err != nil {
		log.Warnf("Failed to check if path exist, error: %#v", err)
	} else if exist {
		config, err = releaseman.NewConfigFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to parse release config at (%s), error: %#v", configPath, err)
		}
	}

	config, err := collectConfigParams(config, c)
	if err != nil {
		log.Fatalf("Failed to collect config params, error: %#v", err)
	}

	//
	// Checkout the development branch
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
	// Print config
	fmt.Println()
	log.Infof("Your config:")
	log.Infof(" * Development branch: %s", config.Release.DevelopmentBranch)
	log.Infof(" * Release branch: %s", config.Release.ReleaseBranch)
	log.Infof(" * Release version: %s", config.Release.Version)
	log.Infof(" * Changelog path: %s", config.Changelog.Path)
	fmt.Println()

	if !releaseman.IsCIMode {
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
	if err := releaseman.WriteChnagelog(commits, relevantTags, config); err != nil {
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
