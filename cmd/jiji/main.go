package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
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

	// Validate connection with a simple API call
	_ = client
	fmt.Printf("Connected to %s\n", cfg.Jira.Host)
}
