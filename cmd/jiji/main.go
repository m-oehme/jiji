package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	tea "charm.land/bubbletea/v2"

	"github.com/m-oehme/jiji/internal/app"
	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/version"
)

func main() {
	var cli config.CLI
	kong.Parse(&cli,
		kong.Name("jiji"),
		kong.Description("A terminal UI for Jira"),
		kong.UsageOnError(),
	)

	if cli.Version {
		fmt.Printf("jiji %s (%s) built %s\n", version.Version, version.Commit, version.Date)
		os.Exit(0)
	}

	if err := cli.ValidateConnection(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		os.Exit(1)
	}

	cfg.Jira = config.JiraConnection{
		Host:  cli.Host,
		Email: cli.Email,
		Token: cli.Token,
	}

	client, err := jira.NewCloudAdapter(cfg.Jira.Host, cfg.Jira.Email, cfg.Jira.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Jira: %s\n", err)
		os.Exit(1)
	}

	var logger *slog.Logger
	if cli.Debug {
		logPath := filepath.Join(config.CacheDir(), "debug.log")
		f, err := os.Create(logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating debug log: %s\n", err)
			os.Exit(1)
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "Error closing debug log: %s\n", cerr)
			}
		}()
		logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
		logger.Info("debug logging enabled", "path", logPath)
	} else {
		logger = slog.New(slog.DiscardHandler)
	}

	m := app.New(cfg, client, logger)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
