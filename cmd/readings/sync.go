package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"productivity.go/internal/config"
	"productivity.go/internal/notion"
	"productivity.go/internal/readings"
	"productivity.go/internal/storage"
)

var syncCmd = &cobra.Command{
	Use:    "sync",
	Hidden: true,
	Short:  "Synchronize articles from Notion to local cache",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			os.Exit(1)
		}

		store, err := storage.NewSQLite()
		if err != nil {
			os.Exit(1)
		}
		defer store.Close()

		notionClient := notion.NewClient(cfg.NotionAPIKey, cfg.NotionDatabaseID, cfg.NotionWeeksDBID)
		svc := readings.NewService(store, notionClient)

		if err := svc.Sync(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Sync failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
