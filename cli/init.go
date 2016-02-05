package cli

import (
	"bufio"
	"fmt"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

//=======================================
// Utility
//=======================================

func collectConfigParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error
	config, err = collectChangeLogConfigParams(config, c)
	if err != nil {
		return releaseman.Config{}, err
	}

	return collectReleaseConfigParams(config, c)
}

//=======================================
// Main
//=======================================

func initRelease(c *cli.Context) {
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
