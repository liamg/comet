package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := checkGitInPath(); err != nil {
		fail("Error: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	dir, err := findGitDir(cwd, 0)
	if err != nil {
		fail("Error: %s", err)
	}

	if err := os.Chdir(dir); err != nil {
		fail("Error: could not change directory: %s", err)
	}

	prefixes, signOff, err := loadConfig()
	if err != nil {
		fail("Error: %s", err)
	}

	m := newModel(prefixes)
	if err := tea.NewProgram(m).Start(); err != nil {
		fail("Error: %s", err)
	}

	fmt.Println("")

	if !m.Finished() {
		fail("Aborted.")
	}

	msg, withBody := m.CommitMessage()
	if err := commit(msg, withBody, signOff); err != nil {
		fail("Error creating commit: %s", err)
	}
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
