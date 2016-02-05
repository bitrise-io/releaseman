package releaseman

import (
	"bufio"
	"os"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-tools/releaseman/git"
)

//=======================================
// Consts
//=======================================

// ChangelogTemplate ...
const ChangelogTemplate = `# Current version: {{.Version}}
# Changelog
{{range .Sections}}### {{.HeaderFrom}} - {{.HeaderTo}}
{{range $idx, $change := .Changes}} * {{$change}}{{ "\n" }}{{end}}
{{end}}`

//=======================================
// Models
//=======================================

// ChangelogSectionModel ...
type ChangelogSectionModel struct {
	HeaderFrom string
	HeaderTo   string
	Changes    []string
}

// ChangelogModel ...
type ChangelogModel struct {
	Version  string
	Sections []ChangelogSectionModel
}

//=======================================
// Utility
//=======================================

func changeList(commits []git.CommitModel) []string {
	changes := []string{}
	for _, commit := range commits {
		changes = append(changes, commit.Message)
	}
	return changes
}

func reverse(commits []git.CommitModel) []git.CommitModel {
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

		if isRelevantCommit && (endDate == nil || (*endDate).Sub(commit.Date) <= 0) {
			return relevantCommits
		}

		if isRelevantCommit {
			relevantCommits = append(relevantCommits, commit)
		}
	}

	return reverse(relevantCommits)
}

func generateChangelog(commits, taggedCommits []git.CommitModel, nextVersion string) ChangelogModel {
	changelog := ChangelogModel{Sections: []ChangelogSectionModel{}}

	if len(taggedCommits) > 0 {
		// Commits between initial commit and first tag
		relevantCommits := commitsBetween(nil, &(taggedCommits[0].Date), commits)

		section := ChangelogSectionModel{
			HeaderFrom: "",
			HeaderTo:   taggedCommits[0].Tag,
			Changes:    changeList(relevantCommits),
		}
		changelog.Sections = append(changelog.Sections, section)

		if len(taggedCommits) > 1 {
			// Commits between tags
			for i := 0; i < len(taggedCommits)-1; i++ {
				startTaggedCommit := taggedCommits[i]
				endTaggedCommit := taggedCommits[i+1]

				relevantCommits = commitsBetween(&(startTaggedCommit.Date), &(endTaggedCommit.Date), commits)

				section = ChangelogSectionModel{
					HeaderFrom: startTaggedCommit.Tag,
					HeaderTo:   endTaggedCommit.Tag,
					Changes:    changeList(relevantCommits),
				}
				changelog.Sections = append(changelog.Sections, section)
			}
		}

		// Commits between last tag and current state
		relevantCommits = commitsBetween(&(taggedCommits[len(taggedCommits)-1].Date), nil, commits)

		section = ChangelogSectionModel{
			HeaderFrom: taggedCommits[len(taggedCommits)-1].Tag,
			HeaderTo:   "",
			Changes:    changeList(relevantCommits),
		}
		changelog.Sections = append(changelog.Sections, section)
	} else {
		relevantCommits := commitsBetween(nil, nil, commits)

		section := ChangelogSectionModel{
			HeaderFrom: "",
			HeaderTo:   "",
			Changes:    changeList(relevantCommits),
		}
		changelog.Sections = append(changelog.Sections, section)
	}

	return changelog
}

// WriteChnagelog ...
func WriteChnagelog(changelogPath string, commits, taggedCommits []git.CommitModel, nextVersion string) error {
	changelog := generateChangelog(commits, taggedCommits, nextVersion)
	changelog.Version = nextVersion

	log.Debugf("Changelog: %#v", changelog)

	tmpl, err := template.New("changelog").Parse(ChangelogTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template, error: %#v", err)
	}

	file, err := os.Create(changelogPath)
	if err != nil {
		log.Fatalf("Failed to create changelog at (%s), error: %#v", changelogPath, err)
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
