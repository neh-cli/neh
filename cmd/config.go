// File: cmd/config.go

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Lang string `yaml:"lang"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var langCmd = &cobra.Command{
	Use:   "lang [ja|en|...]",
	Short: "Set language (e.g. en, ja, zh, es, fr, de)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := args[0]
		configPath := getConfigPath()

		cfg := Config{Lang: lang}
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			fmt.Println("Failed to create config directory:", err)
			return
		}

		f, err := os.Create(configPath)
		if err != nil {
			fmt.Println("Failed to create config file:", err)
			return
		}
		defer f.Close()

		encoder := yaml.NewEncoder(f)
		defer encoder.Close()
		if err := encoder.Encode(&cfg); err != nil {
			fmt.Println("Failed to write config:", err)
		} else {
			fmt.Printf("Language set to '%s' in %s\n", lang, configPath)
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := getConfigPath()

		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			return
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Printf("Failed to parse config file: %v\n", err)
			return
		}

		fmt.Println("Current configuration:")
		fmt.Printf("lang: %s\n", cfg.Lang)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete all configuration (with confirmation)",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := filepath.Dir(getConfigPath())

		fmt.Printf("Are you sure you want to delete %s? (y/N): ", configDir)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Aborted. Configuration not deleted.")
			return
		}

		if err := os.RemoveAll(configDir); err != nil {
			fmt.Printf("Failed to delete config directory: %v\n", err)
			return
		}
		fmt.Println("Configuration reset: directory deleted.")
	},
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "neh", "config.yml")
}

func init() {
	configCmd.AddCommand(langCmd)
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(configCmd)
}
