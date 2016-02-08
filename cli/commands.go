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

	// BumpVersionKey ...
	BumpVersionKey = "bump-version"
	// PatchKey ...
	PatchKey = "patch"
	// MinorKey ...
	MinorKey = "minor"
	// MajorKey ...
	MajorKey = "major"
)

var (
	commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create changelog and release new version",
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
					Name:  VersionKey,
					Usage: "Release version",
				},
				cli.StringFlag{
					Name:  BumpVersionKey,
					Value: "patch",
					Usage: "Bump version (options: patch, minor, major).",
				},
				cli.StringFlag{
					Name:  ChangelogPathKey,
					Usage: "Change log path",
				},
			},
		},
		{
			Name:   "create-changelog",
			Usage:  "Create changelog",
			Action: createChangelog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  DevelopmentBranchKey,
					Usage: "Development branch",
				},
				cli.StringFlag{
					Name:  VersionKey,
					Usage: "Release version",
				},
				cli.StringFlag{
					Name:  BumpVersionKey,
					Value: "patch",
					Usage: "Bump version (options: patch, minor, major).",
				},
				cli.StringFlag{
					Name:  ChangelogPathKey,
					Usage: "changelog path",
				},
			},
		},
		{
			Name:   "create-release",
			Usage:  "Release new version",
			Action: createRelease,
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
					Name:  VersionKey,
					Usage: "Release version",
				},
				cli.StringFlag{
					Name:  BumpVersionKey,
					Value: "patch",
					Usage: "Bump version (options: patch, minor, major).",
				},
			},
		},
		{
			Name:   "init",
			Usage:  "Initialize release configuration",
			Action: initRelease,
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
