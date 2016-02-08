package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/depman/pathutil"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

//=======================================
// Utility
//=======================================

func collectReleaseConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	//
	// Fill development branch
	if config, err = fillDevelopmetnBranch(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Ensure current branch
	if err := ensureCurrentBranch(config); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill release branch
	if config, err = fillReleaseBranch(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill release version
	if config, err = fillVersion(config, c); err != nil {
		return releaseman.Config{}, err
	}

	return config, nil
}

//=======================================
// Main
//=======================================

func createRelease(c *cli.Context) {
	//
	// Fail if git is not clean
	if areChanges, err := git.AreUncommitedChanges(); err != nil {
		log.Fatalf("Failed to get uncommited changes, error: %#v", err)
	} else if areChanges {
		log.Fatalf("There are uncommited changes in your git, please commit your changes before continue release!")
	}

	printRollBackMessage()

	//
	// Build config
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

	config, err := collectReleaseConfigParams(config, c)
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

	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s released 🚀", config.Release.Version))
	log.Infoln("Take a look at your git, and if you are happy with the release, push the changes.")
}
