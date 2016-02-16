package cli

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

//=======================================
// Ask for user input
//=======================================

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

func askForChangelogTemplatePath() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in changelog template path, or press enter to use default one!")
}

//=======================================
// Fill config
//=======================================

func fillDevelopmetnBranch(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	if c.IsSet(DevelopmentBranchKey) {
		config.Release.DevelopmentBranch = c.String(DevelopmentBranchKey)
	}
	if config.Release.DevelopmentBranch == "" {
		if releaseman.IsCIMode {
			return releaseman.Config{}, errors.New("Missing required input: development branch")
		}
		config.Release.DevelopmentBranch, err = askForDevelopmentBranch()
		if err != nil {
			return releaseman.Config{}, err
		}
	}

	if config.Release.DevelopmentBranch == "" {
		return releaseman.Config{}, errors.New("Missing required input: development branch")
	}

	return config, nil
}

func fillReleaseBranch(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	if c.IsSet(ReleaseBranchKey) {
		config.Release.ReleaseBranch = c.String(ReleaseBranchKey)
	}
	if config.Release.ReleaseBranch == "" {
		if releaseman.IsCIMode {
			return releaseman.Config{}, errors.New("Missing required input: release branch")
		}

		config.Release.ReleaseBranch, err = askForReleaseBranch()
		if err != nil {
			return releaseman.Config{}, err
		}
	}

	if config.Release.ReleaseBranch == "" {
		return releaseman.Config{}, errors.New("Missing required input: release branch")
	}

	return config, nil
}

func versionSegmentIdx(segmentStr string) (int, error) {
	segmentIdx := 0
	switch segmentStr {
	case PatchKey:
		segmentIdx = 2
	case MinorKey:
		segmentIdx = 1
	case MajorKey:
		segmentIdx = 0
	default:
		return -1, fmt.Errorf("Invalid segment name (%s)", segmentStr)
	}
	return segmentIdx, nil
}

func fillVersion(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	tags, err := git.TaggedCommits()
	if err != nil {
		return releaseman.Config{}, err
	}

	if c.IsSet(BumpVersionKey) {
		if len(tags) == 0 {
			return releaseman.Config{}, errors.New("There are no tags, nothing to bump")
		}

		segmentIdx, err := versionSegmentIdx(c.String(BumpVersionKey))
		if err != nil {
			return releaseman.Config{}, err
		}
		lastVersion := tags[len(tags)-1].Tag

		config.Release.Version, err = releaseman.BumpedVersion(lastVersion, segmentIdx)
		if err != nil {
			return releaseman.Config{}, err
		}
	} else if c.IsSet(VersionKey) {
		config.Release.Version = c.String(VersionKey)
	}

	if config.Release.Version == "" {
		if releaseman.IsCIMode {
			return releaseman.Config{}, errors.New("Missing required input: release version")
		}

		if len(tags) > 0 {
			fmt.Println()
			log.Infof("Your previous tag: %s", tags[len(tags)-1].Tag)
		}

		version, err := askForReleaseVersion()
		if err != nil {
			return releaseman.Config{}, err
		}

		for _, taggedCommit := range tags {
			if taggedCommit.Tag == version {
				return releaseman.Config{}, fmt.Errorf("Tag (%s) already exist", version)
			}
		}

		config.Release.Version = version
	}

	if config.Release.Version == "" {
		return releaseman.Config{}, errors.New("Missing required input: release version")
	}

	return config, nil
}

func fillChangelogPath(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	if c.IsSet(ChangelogPathKey) {
		config.Changelog.Path = c.String(ChangelogPathKey)
	}
	if config.Changelog.Path == "" {
		if releaseman.IsCIMode {
			return releaseman.Config{}, errors.New("Missing required input: changelog path")
		}

		config.Changelog.Path, err = askForChangelogPath()
		if err != nil {
			return releaseman.Config{}, err
		}
	}

	if config.Changelog.Path == "" {
		return releaseman.Config{}, errors.New("Missing required input: changelog path")
	}

	return config, nil
}

//=======================================
// Ensure
//=======================================

func ensureCleanGit() error {
	if areChanges, err := git.AreUncommitedChanges(); err != nil {
		return err
	} else if areChanges {
		return errors.New("There are uncommited changes in your git, please commit your changes before continue release!")
	}
	return nil
}

func ensureCurrentBranch(config releaseman.Config) error {
	currentBranch, err := git.CurrentBranchName()
	if err != nil {
		return err
	}

	if config.Release.DevelopmentBranch != currentBranch {
		if releaseman.IsCIMode {
			return fmt.Errorf("Your current branch (%s), should be the development branch (%s)!", currentBranch, config.Release.DevelopmentBranch)
		}

		log.Warnf("Your current branch (%s), should be the development branch (%s)!", currentBranch, config.Release.DevelopmentBranch)

		fmt.Println()
		checkout, err := goinp.AskForBool(fmt.Sprintf("Would you like to checkout development branch (%s)?", config.Release.DevelopmentBranch))
		if err != nil {
			return err
		}

		if !checkout {
			return fmt.Errorf("Current branch should be the development branch (%s)!", config.Release.DevelopmentBranch)
		}

		if err := git.CheckoutBranch(config.Release.DevelopmentBranch); err != nil {
			return err
		}
	}

	return nil
}

//=======================================
// Print common messages
//=======================================

func printRollBackMessage() {
	fmt.Println()
	log.Infoln("How to roll-back?")
	log.Infoln("* if you want to undo the last commit you can call:")
	log.Infoln("    $ git reset --hard HEAD~1")
	log.Infoln("* to delete tag:")
	log.Infoln("    $ git tag -d [TAG]")
	log.Infoln("    $ git push origin :refs/tags/[TAG]")
	log.Infoln("* to roll back to the remote state:")
	log.Infoln("    $ git reset --hard origin/[branch-name]")
	fmt.Println()
}

func printCollectingCommits(startCommit git.CommitModel, nextVersion string) {
	fmt.Println()
	if startCommit.Tag != "" {
		log.Infof("Collecting commits between (%s - %s)", startCommit.Tag, nextVersion)
	} else {
		log.Infof("Collecting commits between (initial commit - %s)", nextVersion)
	}
}
