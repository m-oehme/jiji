# Jiji

A terminal UI for Jira Cloud, built with Go and [Bubbletea v2](https://github.com/charmbracelet/bubbletea).

## Installation

> Currently Jiji can only be build from source. At the time of writing this I don't think its in a state where releasing binaries makes sense.

### From source

Requires Go 1.26+.

```bash
git clone https://github.com/m-oehme/jiji.git
cd jiji
make build
```

The binary is built to `bin/jiji`.

## Usage

### Connection

Jiji connects to Jira Cloud using an API token. Provide credentials via flags or environment variables:

```bash
jiji --host https://yourorg.atlassian.net --email you@example.com --token YOUR_API_TOKEN
```

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--host`, `-H` | `JIJI_JIRA_HOST` | Jira instance URL |
| `--email`, `-e` | `JIJI_JIRA_EMAIL` | Jira user email |
| `--token`, `-t` | `JIJI_JIRA_TOKEN` | Jira API token |

You can generate an API token at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

### Configuration

UI settings are configured via a TOML file at `$XDG_CONFIG_HOME/jiji/config.toml` (defaults to `~/.config/jiji/config.toml`). See the [default config](https://github.com/m-oehme/jiji/blob/main/internal/config/default.toml) for all available options and their defaults.

#### Tabs

Define tabs with custom JQL queries:

```toml
[[tabs]]
name = "My Issues"
jql = "assignee = currentUser() AND resolution = Unresolved ORDER BY updated DESC"

[[tabs]]
name = "Sprint"
jql = "sprint in openSprints() AND project = MYPROJ ORDER BY rank ASC"
```

#### Keybindings

Press `?` at any time to view all keybindings.

Built-in keybindings follow vim conventions and can be overridden in `config.toml`:

| Action | Default | Description |
|--------|---------|-------------|
| `up` | `k`, `Up` | Move cursor up |
| `down` | `j`, `Down` | Move cursor down |
| `tab_next` | `l`, `Right` | Next tab |
| `tab_prev` | `h`, `Left` | Previous tab |
| `pane_switch` | `Tab` | Switch between list and detail pane |
| `top` | `g` | Jump to top |
| `bottom` | `G` | Jump to bottom |
| `confirm` | `Enter` | Confirm selection |
| `focus_jql` | `/` | Focus JQL input |
| `cancel` | `Esc` | Cancel / go back |
| `quit` | `q` | Quit |
| `help` | `?` | Show keybinding help |
| `transition` | `t` | Transition issue status |
| `comment` | `c` | Add comment |
| `labels` | `L` | Edit labels |
| `summary` | `s` | Edit summary |
| `edit` | `e` | Edit issue description |
| `refresh` | `r` | Refresh current view |

To override a built-in keybinding, set it under `[keybindings.builtin]`:

```toml
[keybindings.builtin]
quit = ["q", "Q"]
up = ["k"]
```

Custom keybindings can execute shell commands with template variables:

```toml
[[keybindings.user.issues]]
key = "H"
command = "claude 'pull jira ticket {{ .Key }}'"
```

### Flags

```
jiji [flags]

  -H, --host     Jira instance URL ($JIJI_JIRA_HOST)
  -e, --email    Jira user email ($JIJI_JIRA_EMAIL)
  -t, --token    Jira API token ($JIJI_JIRA_TOKEN)
  -d, --debug    Write debug logs to ~/.cache/jiji/debug.log
  -v, --version  Print version and exit
```

## Developing

### Prerequisites

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/)

### Make targets

```bash
make build   # Build binary to bin/jiji
make run     # Run with debug logging (reads .env if present)
make test    # Run all tests
make lint    # Run linter
make watch   # Live reload on file changes (requires watchexec)
make clean   # Remove build artifacts
```

## Disclaimer

The 4 commits after the first commit are AI generated. For me this project is for learning Go and how to build TUIs and to improve my workflow with Jira and AI coding agents (that I use at work).
IMO using AI agents is fine if they automate the boring parts of coding. Building the foundation of an app is not super interesting but learning and tinkering with it is.

I hope Jiji is useful to you. If you have feedback about UX, coding practices, ideas or anything else feel free to open a discussion.
