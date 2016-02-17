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

func collectChangelogConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
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
	// Fill release version
	if config, err = fillVersion(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill changelog path
	if config, err = fillChangelogPath(config, c); err != nil {
		return releaseman.Config{}, err
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
	// Validate config
	config.Print(releaseman.ChangelogMode)

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
	// Run set version script
	if c.IsSet(SetVersionScriptKey) {
		setVersionScript := c.String(SetVersionScriptKey)
		if err := runSetVersionScript(setVersionScript, config.Release.Version); err != nil {
			log.Fatalf("Failed to run set version script, error: %#v", err)
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

	taggedCommits, err := git.VersionTaggedCommits()
	if err != nil {
		log.Fatalf("Failed to get tagged commits, error: %#v", err)
	}

	startDate := startCommit.Date
	endDate := endCommit.Date
	relevantTags := taggedCommits
	appendChangelog := false

	if config.Changelog.Path != "" {
		if exist, err := pathutil.IsPathExists(config.Changelog.Path); err != nil {
			log.Fatalf("Failed to check if path exist, error: %#v", err)
		} else if exist {
			if len(taggedCommits) > 0 {
				lastTaggedCommit := taggedCommits[len(taggedCommits)-1]
				startCommit = lastTaggedCommit
				startDate = lastTaggedCommit.Date
				relevantTags = []git.CommitModel{lastTaggedCommit}
				appendChangelog = true
			}
		}
	}

	printCollectingCommits(startCommit, config.Release.Version)

	fmt.Println()
	log.Infof("=> Generating Changelog...")
	commits, err := git.GetCommitsBetween(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get commits, error: %#v", err)
	}
	if err := releaseman.WriteChangelog(commits, relevantTags, config, appendChangelog); err != nil {
		log.Fatalf("Failed to write Changelog, error: %#v", err)
	}

	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s Changelog created (%s) 🚀", config.Release.Version, config.Changelog.Path))
}
