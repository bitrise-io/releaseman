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
	RawCommand string
	Name       string
	Args       []string
}

// NewPrintableCommand ...
func NewPrintableCommand(commandParts ...string) PrintableCommand {
	name := commandParts[0]
	args := []string{}
	if len(commandParts) > 1 {
		args = commandParts[1:len(commandParts)]
	}

	return PrintableCommand{
		RawCommand: strings.Join(commandParts, " "),
		Name:       name,
		Args:       args,
	}
}

// Run ...
func (command PrintableCommand) Run() (string, error) {
	log.Debugf("=> (%#v)", command)

	out, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr(command.Name, command.Args...)
	if err != nil {
		log.Fatalf("Failed to execute:\ncommand:(%s),\noutput:(%s),\nerror:(%#v)", command.RawCommand, out, err)
	}
	log.Debugf("output:\n(%s)", out)

	return out, err
}

//=======================================
// Util
//=======================================

func splitByNewLineAndStrip(str string) []string {
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
