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

func collectConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
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

func create(c *cli.Context) {
	//
	// Fail if git is not clean
	if err := ensureCleanGit(); err != nil {
		log.Fatalf("Ensure clean git failed, error: %#v", err)
	}

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

	config, err := collectConfigParams(config, c)
	if err != nil {
		log.Fatalf("Failed to collect config params, error: %#v", err)
	}

	printRollBackMessage()

	//
	// Validate config
	config.Print(releaseman.FullMode)

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
				startCommit = taggedCommits[len(taggedCommits)-1]
				startDate = startCommit.Date
				relevantTags = []git.CommitModel{startCommit}
				appendChangelog = true
			}
		}
	}

	printCollectingCommits(startCommit, config.Release.Version)

	fmt.Println()
	log.Infof("=> Generating changelog...")
	commits, err := git.GetCommitsBetween(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to get commits, error: %#v", err)
	}
	if err := releaseman.WriteChangelog(commits, relevantTags, config, appendChangelog); err != nil {
		log.Fatalf("Failed to write changelog, error: %#v", err)
	}

	//
	// Create release git changes
	fmt.Println()
	log.Infof("=> Adding changes to git...")
	changes, err := git.GetChangedFiles()
	if err != nil {
		log.Fatalf("Failed to get changes, error: %s", err)
	}
	if err := git.Add(changes); err != nil {
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

	fmt.Println()
	log.Infoln(colorstring.Greenf("v%s released 🚀", config.Release.Version))
	log.Infoln("Take a look at your git, and if you are happy with the release, push the changes.")
}
