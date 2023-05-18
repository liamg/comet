package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

type prefix struct {
	T string `json:"title"`
	D string `json:"description"`
}

type config struct {
	Prefixes       []prefix `json:"prefixes"`
	SignOffCommits bool     `json:"signOffCommits"`
}

func (i prefix) Title() string       { return i.T }
func (i prefix) Description() string { return i.D }
func (i prefix) FilterValue() string { return i.T }

var defaultPrefixes = []list.Item{
	prefix{
		T: "feat",
		D: "Unleashed a magical unicorn",
	},
	prefix{
		T: "fix",
		D: "Vanquished a sneaky bug",
	},
	prefix{
		T: "docs",
		D: "Penned down a riveting tale about our code",
	},
	prefix{
		T: "test",
		D: "An army of knights to battle-test our code is here",
	},
	prefix{
		T: "build",
		D: "Constructed a towering fortress of code",
	},
	prefix{
		T: "ci",
		D: "Enlisted a team of robot to automate our code",
	},
	prefix{
		T: "perf",
		D: "Taught our code some secret ninja moves",
	},
	prefix{
		T: "refactor",
		D: "Gave our code a dazzling makeover",
	},
	prefix{
		T: "revert",
		D: "Turned back time on our code's misadventures",
	},
	prefix{
		T: "style",
		D: "Dressed our code in the trendiest fashion",
	},
	prefix{
		T: "chore",
		D: "Performed superhero janitorial duties",
		
	},
	prefix{
		T: "add",
		D: "Unleashed file-fairies for a codebase expansion!",
	},
	prefix{
		T: "init",
		D: "Batman! (this commit has no parents)",
	},
}

const configFile = ".comet.json"

func loadConfig() ([]list.Item, bool, error) {

	if _, err := os.Stat(configFile); err == nil {
		return loadConfigFile(configFile)
	}

	if home, err := os.UserHomeDir(); err == nil {
		path := filepath.Join(home, configFile)
		if _, err := os.Stat(path); err == nil {
			return loadConfigFile(path)
		}
	}

	if _, err := os.Stat(configFile); err == nil {
		return loadConfigFile(configFile)
	}

	return defaultPrefixes, false, nil
}

func loadConfigFile(path string) ([]list.Item, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read config file: %w", err)
	}
	var c config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, false, fmt.Errorf("invalid json in config file '%s': %w", path, err)
	}
	return convertPrefixes(c.Prefixes), c.SignOffCommits, nil
}

func convertPrefixes(prefixes []prefix) []list.Item {
	var output []list.Item
	for _, prefix := range prefixes {
		output = append(output, prefix)
	}
	if len(output) == 0 {
		return defaultPrefixes
	}
	return output
}
