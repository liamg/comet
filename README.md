# Comet

Comet is a simple CLI tool that helps you to use [conventional commits](https://www.conventionalcommits.org/) with git.

You can call `comet` where you'd normally type `git commit`. All flags supported in `git commit` will still work.

![Demo](demo.png)

## Installation

Install with Go (1.17+):

```console
go install github.com/JammUtkarsh/comet@latest
```

Or grab a binary from [the latest release](https://github.com/JammUtkarsh/comet/releases/latest).

## Custom Format

You can customise the options available by creating a `.comet.json` file.

Each repository can have its own `.comet.json` file, or you can create one in your home directory to apply to all repositories.

The content should be in the following format:

```json
{
  "signOffCommits": false,
  "Emoji": false,
  "prefixes": [
    { "title":  "feat", "description":  "a new feature", "emoji": "🚀"},
    { "title":  "fix", "description":  "a bug fix", "emojo": "🐛"},
    { "title":  "bug", "description":  "introducing a bug", "emoji": "🐛"},
  ]
}
```

NOTE: Emoji feature is not yet added to the project.
