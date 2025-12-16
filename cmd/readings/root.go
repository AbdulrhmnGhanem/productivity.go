package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"productivity.go/internal/config"
	"productivity.go/internal/notion"
	"productivity.go/internal/readings"
	"productivity.go/internal/storage"
	"productivity.go/internal/sync"
	"productivity.go/internal/tui"
)

var (
	tagFlag string
)

var rootCmd = &cobra.Command{
	Use:   "readings",
	Short: "A CLI for managing your weekly readings from Notion",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\nRun 'readings setup' to configure.\n", err)
			os.Exit(1)
		}

		if err := cfg.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration invalid: %v\nRun 'readings setup' to configure.\n", err)
			os.Exit(1)
		}

		store, err := storage.NewSQLite()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
			os.Exit(1)
		}
		defer store.Close()

		notionClient := notion.NewClient(cfg.NotionAPIKey, cfg.NotionDatabaseID)
		svc := readings.NewService(store, notionClient)

		// Launch TUI
		if err := tui.Start(svc); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}

		// Trigger background sync
		if err := sync.TriggerBackgroundSync(); err != nil {
			// Just log to stderr, don't fail the command
			fmt.Fprintf(os.Stderr, "Failed to trigger background sync: %v\n", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&tagFlag, "tag", "t", "", "Filter by tag")
	rootCmd.AddCommand(setupCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
