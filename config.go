package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/enescakir/emoji"
)

type prefix struct {
	T string `json:"title"`
	D string `json:"description"`
	E string   `json:"emoji"`
}

type config struct {
	Prefixes       []prefix `json:"prefixes"`
	SignOffCommits bool     `json:"signOffCommits"`
	Emoji 	bool     `json:"emoji"`
}

func (i prefix) Title() string       { return i.T }
func (i prefix) Description() string { return i.D }
func (i prefix) FilterValue() string { return i.T }

var defaultPrefixes = []list.Item{
	prefix{
		T: "feat",
		D: "Unleashed a magical unicorn",
		E: emoji.Parse(":sparkles:"),
	},
	prefix{
		T: "fix",
		D: "Vanquished a sneaky bug",
		E: emoji.Parse(":bug:"),
	},
	prefix{
		T: "docs",
		D: "Penned down a riveting tale about our code",
		E: emoji.Parse(":memo:"),
	},
	prefix{
		T: "test",
		D: "An army of knights to battle-test our code is here",
		E: emoji.Parse(":white_check_mark:"),
	},
	prefix{
		T: "build",
		D: "Constructed a towering fortress of code",
		E: emoji.Parse(":hammer:"),
	},
	prefix{
		T: "ci",
		D: "Enlisted a team of robot to automate our code",
		E: emoji.Parse(":robot:"),
	},
	prefix{
		T: "perf",
		D: "Taught our code some secret ninja moves",
		E: emoji.Parse(":zap:"),
	},
	prefix{
		T: "refactor",
		D: "Gave our code a dazzling makeover",
		E: emoji.Parse(":recycle:"),
	},
	prefix{
		T: "revert",
		D: "Turned back time on our code's misadventures",
		E: emoji.Parse(":rewind:"),
	},
	prefix{
		T: "style",
		D: "Dressed our code in the trendiest fashion",
		E: emoji.Parse(":lipstick:"),
	},
	prefix{
		T: "chore",
		D: "Performed superhero janitorial duties",
		E: emoji.Parse(":broom:"),
		
	},
	prefix{
		T: "add",
		D: "Unleashed file-fairies for a codebase expansion!",
		E: emoji.Parse(":heavy_plus_sign:"),
	},
	prefix{
		T: "init",
		D: "Batman! (this commit has no parents)",
		E: emoji.Parse(":tada:"),
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
