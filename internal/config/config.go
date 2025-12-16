package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/spf13/viper"
)

const (
	NetrcMachineName = "notion.so"
	ConfigFileName   = "productivity.go.toml"
	ConfigDirName    = "productivity.go"
)

type Config struct {
	NotionAPIKey     string
	NotionDatabaseID string
	NotionWeeksDBID  string
}

// Load reads configuration from .netrc and productivity.go.toml
func Load() (*Config, error) {
	cfg := &Config{}

	// 1. Load Notion Database ID from TOML
	if err := loadViperConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// 2. Load Notion API Key from .netrc
	if err := loadNetrcConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to load .netrc: %w", err)
	}

	return cfg, nil
}

func loadViperConfig(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.SetConfigName("productivity.go") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Join(home, ".config", ConfigDirName))
	viper.AddConfigPath(".") // optionally look in current directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if we are just starting up (might be first run)
			// However, for the main app to work, we need the DB ID.
			// We'll return nil here and let the validator check for missing fields later.
			return nil
		}
		return err
	}

	cfg.NotionDatabaseID = CleanDatabaseID(viper.GetString("notion_database_id"))
	cfg.NotionWeeksDBID = CleanDatabaseID(viper.GetString("notion_weeks_db_id"))
	return nil
}

func loadNetrcConfig(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	netrcPath := filepath.Join(home, ".netrc")
	// Check if file exists
	if _, err := os.Stat(netrcPath); os.IsNotExist(err) {
		return nil // No .netrc, that's fine for now
	}

	n, err := netrc.ParseFile(netrcPath)
	if err != nil {
		return err
	}

	machine := n.FindMachine(NetrcMachineName)
	if machine != nil {
		cfg.NotionAPIKey = machine.Password
	}

	return nil
}

// Validate checks if the necessary configuration is present
func (c *Config) Validate() error {
	if c.NotionAPIKey == "" {
		return fmt.Errorf("Notion API Key not found in .netrc (machine: %s)", NetrcMachineName)
	}
	if c.NotionDatabaseID == "" {
		return fmt.Errorf("Notion Database ID not found in %s", ConfigFileName)
	}
	if c.NotionWeeksDBID == "" {
		return fmt.Errorf("Notion Weeks Database ID not found in %s", ConfigFileName)
	}
	return nil
}

// Save writes the configuration to disk
func Save(apiKey, databaseID, weeksDBID string) error {
	if err := saveNetrc(apiKey); err != nil {
		return fmt.Errorf("failed to save .netrc: %w", err)
	}
	if err := saveViper(databaseID, weeksDBID); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}
	return nil
}

func saveNetrc(apiKey string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	netrcPath := filepath.Join(home, ".netrc")

	var n *netrc.Netrc
	if _, err := os.Stat(netrcPath); os.IsNotExist(err) {
		n = &netrc.Netrc{}
	} else {
		var err error
		n, err = netrc.ParseFile(netrcPath)
		if err != nil {
			return err
		}
	}

	machine := n.FindMachine(NetrcMachineName)
	if machine == nil {
		machine = n.NewMachine(NetrcMachineName, "apikey", apiKey, "")
	} else {
		machine.Password = apiKey
		machine.Login = "apikey"
	}

	data, err := n.MarshalText()
	if err != nil {
		return err
	}

	// Set 0600 permissions for security
	return os.WriteFile(netrcPath, data, 0600)
}

func saveViper(databaseID, weeksDBID string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".config", ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	viper.Set("notion_database_id", databaseID)
	viper.Set("notion_weeks_db_id", weeksDBID)
	
	// WriteConfigAs will overwrite or create
	return viper.WriteConfigAs(filepath.Join(configDir, ConfigFileName))
}

// CleanDatabaseID extracts the database ID from a URL if necessary
func CleanDatabaseID(id string) string {
	if strings.Contains(id, "notion.so") {
		if !strings.HasPrefix(id, "http") {
			id = "https://" + id
		}

		u, err := url.Parse(id)
		if err != nil {
			return id
		}

		path := u.Path
		path = strings.TrimSuffix(path, "/")

		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]

			// If it's a UUID (36 chars), return it
			if len(lastPart) == 36 {
				return lastPart
			}

			// If it's Name-ID (ID is 32 chars hex)
			if len(lastPart) >= 32 {
				return lastPart[len(lastPart)-32:]
			}
			return lastPart
		}
	}

	// Handle case where user pasted ID with query params but no domain
	if idx := strings.Index(id, "?"); idx != -1 {
		return id[:idx]
	}

	return id
}
