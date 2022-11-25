package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const maxGitRecursion = 32

func checkGitInPath() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("cannot find git in PATH: %w", err)
	}
	return nil
}

func findGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf(string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

func commit(msg string, body bool, signOff bool) error {
	args := append([]string{
		"commit", "-m", msg,
	}, os.Args[1:]...)
	if body {
		args = append(args, "-e")
	}
	if signOff {
		args = append(args, "-s")
	}
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func verifyCommitMessage(msg string) bool {
	prefixes, _, err := loadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	title, msg, ok := strings.Cut(msg, ":")
	if !ok || msg == "" {
		return false
	}

	for _, v := range prefixes {
		if v.FilterValue() == title {
			return true
		}
	}

	return false
}