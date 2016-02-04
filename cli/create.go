package cli

import (
	"fmt"
	"strings"

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

func printDoneMessage(config configs.ReleaseConfig) {
	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s released 🚀", config.ReleaseVersion))
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

func askForStartState() (string, error) {
	tags, err := git.ListTags()
	if err != nil {
		return "", err
	}

	states := append([]string{configs.InitialCommitStr}, tags...)

	fmt.Println()
	startState, err := goinp.SelectFromStrings("Select release start state!", states)
	if err != nil {
		return "", err
	}

	if startState == configs.InitialCommitStr {
		startState = ""
	}

	return startState, nil
}

func askForEndState(startState string) (string, error) {
	tags, err := git.ListTags()
	if err != nil {
		return "", err
	}

	// disclude all statets before the from state
	states := append(tags, configs.CurrentStateStr)
	fromStateIdx := -1
	for idx, state := range states {
		if state == startState {
			fromStateIdx = idx
		}
	}
	if fromStateIdx > -1 {
		states = states[fromStateIdx+1 : len(states)]
	}

	fmt.Println()
	endState, err := goinp.SelectFromStrings("Select release end state!", states)
	if err != nil {
		log.Fatalf("Failed to select state, error: %#v", err)
	}

	if endState == configs.CurrentStateStr {
		endState = ""
	}

	return endState, nil
}

func askForReleaseVersion() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in release version!")
}

func askForChangelogPath() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in changelog path!")
}

func getCommitHashes(startState, endState string) (string, string, error) {
	firstCommitHash := ""
	lastCommitHash := ""

	if startState == "" && endState == "" {
		firstCommitHash, _ = git.FirstCommit()
		lastCommitHash, _ = git.LatestCommit()
	} else if startState == "" {
		firstCommitHash, _ = git.FirstCommit()
		lastCommitHash, _ = git.CommitHashOfTag(endState)
	} else if endState == "" {
		firstCommitHash, _ = git.CommitHashOfTag(startState)
		lastCommitHash, _ = git.LatestCommit()
	} else {
		firstCommitHash, _ = git.CommitHashOfTag(startState)
		lastCommitHash, _ = git.CommitHashOfTag(endState)
	}

	return firstCommitHash, lastCommitHash, nil
}

func writeChnagelog(changelogPath, startState, endState string, commits []map[string]string) error {
	if startState == "" {
		startState = configs.InitialCommitStr
	}

	if endState == "" {
		endState = configs.CurrentStateStr
	}

	changelog := "\n"
	changelog += fmt.Sprintf("%s - %s\n", startState, endState)
	for _, commit := range commits {
		for _, message := range commit {
			changelog += fmt.Sprintf(" * %s\n", message)
		}
	}

	if exist, err := pathutil.IsPathExists(changelogPath); err != nil {
		log.Fatalf("Failed to check if path exist (%s), error: %s", changelogPath, err)
	} else if exist {
		log.Infoln("   Previous changelog exist, appending current.")

		previousChangeLog, err := fileutil.ReadStringFromFile(changelogPath)
		if err != nil {
			log.Fatalf("Failed to append new changelog, error: %s", err)
		}

		changelog += previousChangeLog
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
	config := configs.ReleaseConfig{}
	configPath := ""
	if c.IsSet("config") {
		configPath = c.String("config")
	} else {
		configPath = configs.DefaultConfigPth
	}

	if exist, err := pathutil.IsPathExists(configPath); err != nil {
		log.Warnf("Failed to check if path exist, error: %#v", err)
	} else if exist {
		config, err = configs.NewReleaseConfigFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to parse release config at (%s), error: %#v", configPath, err)
		}
	}

	var err error
	if c.IsSet(DevelopmentBranchKey) {
		config.DevelopmentBranch = c.String(DevelopmentBranchKey)
	}
	if config.DevelopmentBranch == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: development branch")
		} else {
			config.DevelopmentBranch, err = askForDevelopmentBranch()
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

	if config.DevelopmentBranch != currentBranch {
		log.Warnf("Your current branch (%s), should be the development branch (%s)!", currentBranch, config.DevelopmentBranch)

		fmt.Println()
		checkout, err := goinp.AskForBool(fmt.Sprintf("Would you like to checkout development branch (%s)?", config.DevelopmentBranch))
		if err != nil {
			log.Fatalf("Failed to ask for checkout, error: %#v", err)
		}

		if !checkout {
			log.Fatalf("Current branch should be the development branch (%s)!", config.DevelopmentBranch)
		}

		if err := git.CheckoutBranch(config.DevelopmentBranch); err != nil {
			log.Fatalf("Failed to checkout branch (%s), error: %#v", config.DevelopmentBranch, err)
		}
	}
	//
	//

	if c.IsSet(ReleaseBranchKey) {
		config.ReleaseBranch = c.String(ReleaseBranchKey)
	}
	if config.ReleaseBranch == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: release branch")
		} else {
			config.ReleaseBranch, err = askForReleaseBranch()
			if err != nil {
				log.Fatalf("Failed to ask for release branch, error: %s", err)
			}
		}
	}

	if c.IsSet(StartStateKey) {
		config.StartState = c.String(StartStateKey)
	}
	if config.StartState == "" {
		if configs.IsCIMode {
			// In CI it means start from repo create
		} else {
			config.StartState, err = askForStartState()
			if err != nil {
				log.Fatalf("Failed to ask for start state, error: %s", err)
			}
		}
	}

	if c.IsSet(EndStateKey) {
		config.EndState = c.String(EndStateKey)
	}
	if config.EndState == "" {
		if configs.IsCIMode {
			// In CI it means create changelog until current state
		} else {
			config.EndState, err = askForEndState(config.StartState)
			if err != nil {
				log.Fatalf("Failed to ask for end state, error: %s", err)
			}
		}
	}

	if c.IsSet(ReleaseVersionKey) {
		config.ReleaseVersion = c.String(ReleaseVersionKey)
	}
	if config.ReleaseVersion == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: release version")
		} else {
			config.ReleaseVersion, err = askForReleaseVersion()
			if err != nil {
				log.Fatalf("Failed to ask for release version, error: %s", err)
			}
		}
	}

	if c.IsSet(ChangelogPathKey) {
		config.ChangelogPath = c.String(ChangelogPathKey)
	}
	if config.ChangelogPath == "" {
		if configs.IsCIMode {
			log.Fatalln("Missing required input: changelog path")
		} else {
			config.ChangelogPath, err = askForChangelogPath()
			if err != nil {
				log.Fatalf("Failed to ask for changelog path, error: %s", err)
			}
		}
	}

	//
	// Print config
	fmt.Println()
	log.Infof("Your config:")
	log.Infof(" * Development branch: %s", config.DevelopmentBranch)
	log.Infof(" * Release branch: %s", config.ReleaseBranch)
	if config.StartState == "" && config.EndState == "" {
		log.Infof(" * Create release from initial commit until current state")
	} else if config.StartState == "" {
		log.Infof(" * Create release from initial commit until %s", config.EndState)
	} else if config.EndState == "" {
		log.Infof(" * Create release from %s until current state", config.StartState)
	} else {
		log.Infof(" * Create release from %s until %s", config.StartState, config.EndState)
	}
	log.Infof(" * Changelog path: %s", config.ChangelogPath)
	log.Infof(" * Release version: %s", config.ReleaseVersion)
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
	startCommitHash, endCommitHash, err := getCommitHashes(config.StartState, config.EndState)
	if err != nil {
		log.Fatalf("Failed to get commit hashes, error: %s", err)
	}

	fmt.Println()
	log.Infof("start commit hash: %s", startCommitHash)
	log.Infof("end commit hash: %s", endCommitHash)

	fmt.Println()
	log.Infof("=> Generating changelog...")
	commits, err := git.CommitMessages(startCommitHash, endCommitHash)
	if err := writeChnagelog(config.ChangelogPath, config.StartState, config.EndState, commits); err != nil {
		log.Fatalf("Failed to write changelog, error: %#v", err)
	}

	fmt.Println()
	log.Infof("=> Adding changes to git...")
	if err := git.Add([]string{config.ChangelogPath}); err != nil {
		log.Fatalf("Failed to git add, error: %s", err)
	}

	if err := git.Commit(fmt.Sprintf("v%s", config.ReleaseVersion)); err != nil {
		log.Fatalf("Failed to git commit, error: %s", err)
	}

	fmt.Println()
	log.Infof("=> Merging changes into release branch...")
	if err := git.CheckoutBranch(config.ReleaseBranch); err != nil {
		log.Fatalf("Failed to git checkout, error: %s", err)
	}

	mergeCommitMessage := fmt.Sprintf("Merge %s into %s, release: v%s", config.DevelopmentBranch, config.ReleaseBranch, config.ReleaseVersion)
	if err := git.Merge(config.DevelopmentBranch, mergeCommitMessage); err != nil {
		log.Fatalf("Failed to git merge, error: %s", err)
	}

	fmt.Println()
	log.Infof("=> Tagging release branch...")
	if err := git.Tag(config.ReleaseVersion); err != nil {
		log.Fatalf("Failed to git tag, error: %s", err)
	}

	if err := git.CheckoutBranch(config.DevelopmentBranch); err != nil {
		log.Fatalf("Failed to git checkout, error: %s", err)
	}

	printDoneMessage(config)
}
