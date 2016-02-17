package releaseman

import (
	"bytes"
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
const ChangelogTemplate = `{{range .Sections}}### {{.StartTaggedCommit.Tag}} - {{.EndTaggedCommit.Tag}} ({{.EndTaggedCommit.Date.Format "2006 Jan 02"}})

{{range $idx, $commit := .Commits}}* [{{trimm $commit.Hash 7}}] {{$commit.Author}} - {{$commit.Message}} ({{$commit.Date.Format "2006 Jan 02"}})
{{end}}
{{end}}`

var changelogTemplateFuncMap = template.FuncMap{
	"trimm": func(str string, length int) string {
		if len(str) < length {
			return str
		}

		return str[0:length]
	},
}

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
			relevantCommits = append([]git.CommitModel{commit}, relevantCommits...)
		}
	}

	return relevantCommits
}

func reversedSections(sections []ChangelogSectionModel) []ChangelogSectionModel {
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

	changelog.Sections = reversedSections(changelog.Sections)

	return changelog
}

//=======================================
// Main
//=======================================

// WriteChangelog ...
func WriteChangelog(commits, taggedCommits []git.CommitModel, config Config, append bool) error {
	changelog := generateChangelog(commits, taggedCommits, config.Release.Version)
	log.Debugf("Changelog: %#v", changelog)

	changelogItemTemplateStr := ChangelogTemplate
	if config.Changelog.ItemTemplate != "" {
		changelogItemTemplateStr = config.Changelog.ItemTemplate
	}

	tmpl := template.New("changelog").Funcs(changelogTemplateFuncMap)
	tmpl, err := tmpl.Parse(changelogItemTemplateStr)
	if err != nil {
		log.Fatalf("Failed to parse template, error: %#v", err)
	}

	var changelogBytes bytes.Buffer
	err = tmpl.Execute(&changelogBytes, changelog)
	if err != nil {
		log.Fatalf("Failed to execute template, error: %#v", err)
	}
	changelogStr := changelogBytes.String()

	if append {
		log.Debugln("Changelog exist, append new version")
		previosChangelogStr, err := fileutil.ReadStringFromFile(config.Changelog.Path)
		if err != nil {
			return err
		}

		changelogStr = changelogStr + previosChangelogStr
	}

	return fileutil.WriteStringToFile(config.Changelog.Path, changelogStr)
}
