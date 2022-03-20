package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const maxGitRecursion = 32

func checkGitInPath() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("cannot find git in PATH: %w", err)
	}
	return nil
}

func findGitDir(start string, count int) (string, error) {
	if info, err := os.Stat(filepath.Join(start, ".git")); err == nil {
		if info.IsDir() {
			return start, nil
		}
	}
	above := filepath.Dir(start)
	if above == start || count >= maxGitRecursion {
		return "", fmt.Errorf("not a git repository: .git directory not found")
	}
	return findGitDir(above, count+1)
}

func commit(msg string, body bool) error {
	args := append([]string{
		"commit", "-m", msg,
	}, os.Args[1:]...)
	if body {
		args = append(args, "-e")
	}
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
