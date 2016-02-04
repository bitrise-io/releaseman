package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPrintableCommand(t *testing.T) {
	command := NewPrintableCommand("git", "branch", "--list")
	require.Equal(t, "git branch --list", command.RawCommand)
	require.Equal(t, "git", command.Name)
	require.Equal(t, []string{"branch", "--list"}, command.Args)

	command = NewPrintableCommand("git")
	require.Equal(t, "git", command.RawCommand)
	require.Equal(t, "git", command.Name)
	require.Equal(t, []string{}, command.Args)

	command = NewPrintableCommand("")
	require.Equal(t, "", command.RawCommand)
	require.Equal(t, "", command.Name)
	require.Equal(t, []string{}, command.Args)
}

func TestRun(t *testing.T) {
	command := NewPrintableCommand("echo", "Hello World!")
	out, err := command.Run()
	require.Equal(t, nil, err)
	require.Equal(t, "Hello World!", out)
}

func TestSplitByNewLineAndStrip(t *testing.T) {
	str := `1. line
2. line
3. line`
	require.Equal(t, []string{"1. line", "2. line", "3. line"}, splitByNewLineAndStrip(str))

	str = `
1. line
2. line
3. line
`
	require.Equal(t, []string{"1. line", "2. line", "3. line"}, splitByNewLineAndStrip(str))

	str = `
    1. line
  2. line
      3. line
`
	require.Equal(t, []string{"1. line", "2. line", "3. line"}, splitByNewLineAndStrip(str))
}

func TestStrip(t *testing.T) {
	str := "test case"
	require.Equal(t, "test case", strip(str))

	str = " test case"
	require.Equal(t, "test case", strip(str))

	str = "test case "
	require.Equal(t, "test case", strip(str))

	str = "   test case   "
	require.Equal(t, "test case", strip(str))

	str = ""
	require.Equal(t, "", strip(str))
}
