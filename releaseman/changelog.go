package releaseman

import (
	"bufio"
	"os"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/releaseman/git"
)

//=======================================
// Consts
//=======================================

// ChangelogTemplate ...
const ChangelogTemplate = `# Current version: {{.Version}}
===
# Change log
===
{{range .Sections}}### {{.StartTaggedCommit.Tag}} - {{.EndTaggedCommit.Tag}} ({{.EndTaggedCommit.Date}})
---
{{range $idx, $commit := .Commits}} * {{$commit.Message}}{{ "\n" }}{{end}}
{{end}}`

//=======================================
// Models
//=======================================

// ChangelogSectionModel ...
type ChangelogSectionModel struct {
	StartTaggedCommit git.CommitModel
	EndTaggedCommit   git.CommitModel
	Commits           []git.CommitModel
}

// ChangelogModel ...
type ChangelogModel struct {
	Version  string
	Sections []ChangelogSectionModel
}

//=======================================
// Utility
//=======================================

func reverseCommits(commits []git.CommitModel) []git.CommitModel {
	reversed := []git.CommitModel{}
	for i := len(commits) - 1; i >= 0; i-- {
		reversed = append(reversed, commits[i])
	}
	return reversed
}

func commitsBetween(startDate *time.Time, endDate *time.Time, commits []git.CommitModel) []git.CommitModel {
	relevantCommits := []git.CommitModel{}
	isRelevantCommit := false

	for _, commit := range commits {
		if !isRelevantCommit && (startDate == nil || (*startDate).Sub(commit.Date) <= 0) {
			isRelevantCommit = true
		}

		if isRelevantCommit && endDate != nil && (*endDate).Sub(commit.Date) <= 0 {
			return relevantCommits
		}

		if isRelevantCommit {
			relevantCommits = append(relevantCommits, commit)
		}
	}

	return relevantCommits
}

func changeList(commits []git.CommitModel) []string {
	changes := []string{}
	for _, commit := range commits {
		changes = append(changes, commit.Message)
	}
	return changes
}

func reverseSections(sections []ChangelogSectionModel) []ChangelogSectionModel {
	reversed := []ChangelogSectionModel{}
	for i := len(sections) - 1; i >= 0; i-- {
		reversed = append(reversed, sections[i])
	}
	return reversed
}

func generateChangelog(commits, taggedCommits []git.CommitModel, version string) ChangelogModel {
	changelog := ChangelogModel{
		Version:  version,
		Sections: []ChangelogSectionModel{},
	}

	if len(taggedCommits) > 0 {
		if len(taggedCommits) > 1 {
			// Commits between tags
			for i := 0; i < len(taggedCommits)-1; i++ {
				startTaggedCommit := taggedCommits[i]
				endTaggedCommit := taggedCommits[i+1]

				relevantCommits := commitsBetween(&(startTaggedCommit.Date), &(endTaggedCommit.Date), commits)

				section := ChangelogSectionModel{
					StartTaggedCommit: startTaggedCommit,
					EndTaggedCommit:   endTaggedCommit,
					Commits:           relevantCommits,
				}
				changelog.Sections = append(changelog.Sections, section)
			}
		}

		// Commits between last tag and current state
		relevantCommits := commitsBetween(&(taggedCommits[len(taggedCommits)-1].Date), nil, commits)

		section := ChangelogSectionModel{
			StartTaggedCommit: taggedCommits[len(taggedCommits)-1],
			EndTaggedCommit: git.CommitModel{
				Tag:  version,
				Date: time.Now(),
			},
			Commits: relevantCommits,
		}
		changelog.Sections = append(changelog.Sections, section)
	} else {
		relevantCommits := commitsBetween(nil, nil, commits)

		section := ChangelogSectionModel{
			StartTaggedCommit: git.CommitModel{},
			EndTaggedCommit: git.CommitModel{
				Tag:  version,
				Date: time.Now(),
			},
			Commits: relevantCommits,
		}
		changelog.Sections = append(changelog.Sections, section)
	}

	changelog.Sections = reverseSections(changelog.Sections)

	return changelog
}

// WriteChnagelog ...
func WriteChnagelog(commits, taggedCommits []git.CommitModel, config Config) error {
	changelog := generateChangelog(commits, taggedCommits, config.Release.Version)
	log.Debugf("Changelog: %#v", changelog)

	changelogTemplate := ChangelogTemplate
	if config.Changelog.TemplatePath != "" {
		var err error
		changelogTemplate, err = fileutil.ReadStringFromFile(config.Changelog.TemplatePath)
		if err != nil {
			log.Fatalf("Failed to read changelog template, error: %#v", err)
		}
	}

	tmpl, err := template.New("changelog").Parse(changelogTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template, error: %#v", err)
	}

	file, err := os.Create(config.Changelog.Path)
	if err != nil {
		log.Fatalf("Failed to create changelog at (%s), error: %#v", config.Changelog.Path, err)
	}
	fileWriter := bufio.NewWriter(file)

	err = tmpl.Execute(fileWriter, changelog)
	if err != nil {
		log.Fatalf("Failed to execute template, error: %#v", err)
	}

	if err = fileWriter.Flush(); err != nil {
		log.Fatalf("Failed to flush changelog file, error: %#v", err)
	}

	return nil
}
