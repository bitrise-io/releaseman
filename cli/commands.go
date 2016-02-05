package cli

import "github.com/codegangsta/cli"

const (
	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// CIKey ...
	CIKey = "ci"
	// CIModeEnvKey ...
	CIModeEnvKey = "CI"

	// DevelopmentBranchKey ...
	DevelopmentBranchKey = "development-branch"

	// ReleaseBranchKey ...
	ReleaseBranchKey = "release-branch"

	// StartStateKey ...
	StartStateKey = "start-state"

	// EndStateKey ...
	EndStateKey = "end-state"

	// ChangelogPathKey ...
	ChangelogPathKey = "changelog-path"

	// ReleaseVersionKey ...
	ReleaseVersionKey = "release-version"
)

var (
	commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Creates new release",
			Action: create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  DevelopmentBranchKey,
					Usage: "Development branch",
				},
				cli.StringFlag{
					Name:  ReleaseBranchKey,
					Usage: "Release branch",
				},
				cli.StringFlag{
					Name:  StartStateKey,
					Usage: "Start state [initial commit/tag]",
				},
				cli.StringFlag{
					Name:  EndStateKey,
					Usage: "End state [tag/current state]",
				},
				cli.StringFlag{
					Name:  ChangelogPathKey,
					Usage: "Changelog path",
				},
				cli.StringFlag{
					Name:  ReleaseVersionKey,
					Usage: "Release version",
				},
			},
		},
		{
			Name:   "init",
			Usage:  "Initialize release_config.yml",
			Action: initRelease,
			Flags:  []cli.Flag{},
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
		cli.BoolFlag{
			Name:   CIKey,
			Usage:  "If true it indicates that we're used by another tool so don't require any user input!",
			EnvVar: CIModeEnvKey,
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}
