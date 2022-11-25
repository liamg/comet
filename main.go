package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := stdin(); err != nil {
		fail("Error: %s", err)
	}

	if err := checkGitInPath(); err != nil {
		fail("Error: %s", err)
	}

	gitRoot, err := findGitDir()
	if err != nil {
		fail("Error: %s", err)
	}

	if err := os.Chdir(gitRoot); err != nil {
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

func stdin() error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("cannot check STDIN: %w", err)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		out, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("cannot read from STDIN: %w", err)
		}

		if string(out) != "" {
			if !verifyCommitMessage(string(out)) {
				return fmt.Errorf("invalid commit message")
			}

			os.Exit(0)
		}
	}

	return nil
}

func fail(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
