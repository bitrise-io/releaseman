package git

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

//=======================================
// Model
//=======================================

// PrintableCommand ...
type PrintableCommand struct {
	rawCommand string
	name       string
	args       []string
}

// NewPrintableCommand ...
func NewPrintableCommand(commandParts ...string) PrintableCommand {
	name := commandParts[0]
	args := []string{}
	if len(commandParts) > 1 {
		args = commandParts[1:len(commandParts)]
	}

	return PrintableCommand{
		rawCommand: strings.Join(commandParts, " "),
		name:       name,
		args:       args,
	}
}

// Run ...
func (command PrintableCommand) Run() (string, error) {
	log.Debugf("=> (%#v)", command)

	out, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr(command.name, command.args...)
	if err != nil {
		log.Fatalf("Failed to execute:\ncommand:(%s),\noutput:(%s),\nerror:(%#v)", command.rawCommand, out, err)
	}
	log.Debugf("output:\n(%s)", out)

	return out, err
}

//=======================================
// Util
//=======================================

func splitByNewLine(str string) []string {
	split := strings.Split(str, "\n")
	out := []string{}
	for _, part := range split {
		if part == "" {
			continue
		}

		out = append(out, strip(part))
	}
	return out
}

func strip(str string) string {
	dirty := true
	strippedStr := str
	for dirty {
		hasWhiteSpacePrefix := false
		if strings.HasPrefix(strippedStr, " ") {
			hasWhiteSpacePrefix = true
			strippedStr = strings.TrimPrefix(strippedStr, " ")
		}

		hasWhiteSpaceSuffix := false
		if strings.HasSuffix(strippedStr, " ") {
			hasWhiteSpaceSuffix = true
			strippedStr = strings.TrimSuffix(strippedStr, " ")
		}

		if !hasWhiteSpacePrefix && !hasWhiteSpaceSuffix {
			dirty = false
		}
	}
	return strippedStr
}
