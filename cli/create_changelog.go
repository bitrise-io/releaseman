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

func askForChangelogPath() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in changelog path!")
}

func askForChangelogTemplatePath() (string, error) {
	fmt.Println()
	return goinp.AskForString("Type in changelog template path, or press enter to use default one!")
}

func collectChangelogConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
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

func createChangelog(c *cli.Context) {
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

	config, err := collectChangelogConfigParams(config, c)
	if err != nil {
		log.Fatalf("Failed to collect config params, error: %#v", err)
	}

	//
	// Print config
	fmt.Println()
	log.Infof("Your config:")
	log.Infof(" * Development branch: %s", config.Release.DevelopmentBranch)
	log.Infof(" * Release version: %s", config.Release.Version)
	log.Infof(" * Changelog path: %s", config.Changelog.Path)
	if config.Changelog.TemplatePath != "" {
		log.Infof(" * Changelog template path: %s", config.Changelog.TemplatePath)
	}
	fmt.Println()

	if !releaseman.IsCIMode {
		ok, err := goinp.AskForBool("Are you ready for creating Changelog?")
		if err != nil {
			log.Fatalf("Failed to ask for input, error: %s", err)
		}
		if !ok {
			log.Fatal("Aborted create Changelog")
		}
	}

	//
	// Generate Changelog
	startCommit, err := git.FirstCommit()
	if err != nil {
		log.Fatalf("Failed to get first commit, error: %#v", err)
	}

	endCommit, err := git.LatestCommit()
	if err != nil {
		log.Fatalf("Failed to get latest commit, error: %#v", err)
	}

	taggedCommits, err := git.TaggedCommits()
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
	log.Infof("=> Generating Changelog...")
	commits, err := git.GetCommitsBetween(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get commits, error: %#v", err)
	}
	if err := releaseman.WriteChnagelog(commits, relevantTags, config); err != nil {
		log.Fatalf("Failed to write Changelog, error: %#v", err)
	}

	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s Changelog created (%s) 🚀", config.Release.Version, config.Changelog.Path))
}
