package cli

import (
	"bufio"
	"fmt"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

//=======================================
// Utility
//=======================================

func collectConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error
	if c.IsSet(DevelopmentBranchKey) {
		config.Release.DevelopmentBranch = c.String(DevelopmentBranchKey)
	}
	if config.Release.DevelopmentBranch == "" {
		if releaseman.IsCIMode {
			log.Fatalln("Missing required input: development branch")
		} else {
			config.Release.DevelopmentBranch, err = askForDevelopmentBranch()
			if err != nil {
				log.Fatalf("Failed to ask for development branch, error: %s", err)
			}
		}
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

	if c.IsSet(ReleaseBranchKey) {
		config.Release.ReleaseBranch = c.String(ReleaseBranchKey)
	}
	if config.Release.ReleaseBranch == "" {
		if releaseman.IsCIMode {
			log.Fatalln("Missing required input: release branch")
		} else {
			config.Release.ReleaseBranch, err = askForReleaseBranch()
			if err != nil {
				log.Fatalf("Failed to ask for release branch, error: %s", err)
			}
		}
	}

	if c.IsSet(VersionKey) {
		config.Release.Version = c.String(VersionKey)
	}
	if config.Release.Version == "" {
		if releaseman.IsCIMode {
			log.Fatalln("Missing required input: release version")
		} else {
			tags, err := git.TaggedCommits()
			if err != nil {
				log.Fatalf("Failed to list tagged commits, error: %#v", err)
			}

			if len(tags) > 0 {
				fmt.Println()
				log.Infof("Your previous tags:")
				for _, taggedCommit := range tags {
					fmt.Printf("* %s\n", taggedCommit.Tag)
				}
			}

			version, err := askForReleaseVersion()
			if err != nil {
				log.Fatalf("Failed to ask for release version, error: %s", err)
			}

			for _, taggedCommit := range tags {
				if taggedCommit.Tag == version {
					log.Fatalf("Tag (%s) already exist", version)
				}
			}

			config.Release.Version = version
		}
	}

	if c.IsSet(ChangelogPathKey) {
		config.Changelog.Path = c.String(ChangelogPathKey)
	}
	if config.Changelog.Path == "" {
		if releaseman.IsCIMode {
			log.Fatalln("Missing required input: changelog path")
		} else {
			config.Changelog.Path, err = askForChangelogPath()
			if err != nil {
				log.Fatalf("Failed to ask for changelog path, error: %s", err)
			}
		}
	}

	if c.IsSet(ChangelogTemplatePathKey) {
		config.Changelog.TemplatePath = c.String(ChangelogTemplatePathKey)
	}
	if config.Changelog.TemplatePath == "" {
		if releaseman.IsCIMode {
			// Use default changelog template
		} else {
			config.Changelog.TemplatePath, err = askForChangelogTemplatePath()
			if err != nil {
				log.Fatalf("Failed to ask for changelog path, error: %s", err)
			}
		}
	}

	return config, nil
}

//=======================================
// Main
//=======================================

func initRelease(c *cli.Context) {
	//
	// Fail if git is not clean
	if areChanges, err := git.AreUncommitedChanges(); err != nil {
		log.Fatalf("Failed to get uncommited changes, error: %#v", err)
	} else if areChanges {
		log.Fatalf("There are uncommited changes in your git, please commit your changes before continue release!")
	}

	if exist, err := pathutil.IsPathExists(releaseman.DefaultConfigPth); err != nil {
		log.Fatalf("Failed to check path (%s), error: %#v", releaseman.DefaultConfigPth, err)
	} else if exist {
		if releaseman.IsCIMode {
			log.Fatalf("Release config already exist at (%s)", releaseman.DefaultConfigPth)
		} else {
			ok, err := goinp.AskForBool(fmt.Sprintf("Release config already exist at (%s), would you like to overwrite it?", releaseman.DefaultConfigPth))
			if err != nil {
				log.Fatalf("Failed to ask for bool, error: %#v", err)
			} else if !ok {
				log.Fatalln("Create release config aborted")
			}
		}
	}

	releaseConfig, err := collectConfigParams(releaseman.Config{}, c)
	if err != nil {
		log.Fatalf("Failed to collect config params, error: %#v", err)
	}

	log.Infof("releaseConfig: %#v", releaseConfig)

	tmpl, err := template.New("config").Parse(releaseman.ReleaseConfigTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template, error: %#v", err)
	}

	file, err := os.Create(releaseman.DefaultConfigPth)
	if err != nil {
		log.Fatalf("Failed to create realse config at (%s), error: %#v", releaseman.DefaultConfigPth, err)
	}
	fileWriter := bufio.NewWriter(file)

	err = tmpl.Execute(fileWriter, releaseConfig)
	if err != nil {
		log.Fatalf("Failed to execute template, error: %#v", err)
	}

	if err = fileWriter.Flush(); err != nil {
		log.Fatalf("Failed to flush release config file, error: %#v", err)
	}
}
