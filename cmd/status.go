// File: cmd/status.go

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neh-cli/neh/cmd/shared"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// The statusCmd is a command to remove all query history.
// By using this command, you can delete all previous query history
// stored on the server.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status",
	Long:  `Check the status of the connection endpoint and the subscribed plan.`,
	Run:   runStatusCmd,
}

// The init function is used to add the statusCmd to the rootCmd.
// By using this function, the status command will be recognized
// as part of the CLI.
func init() {
	rootCmd.AddCommand(statusCmd)
}

// The runStatusCmd function is called when the status command is executed.
// This function uses the shared.ExecuteWebSocketCommand function to send
// the "status" command to the server.
// It is set to send the command without waiting for a response from the server.
func runStatusCmd(cmd *cobra.Command, args []string) {
	queryMessage := ""
	clipboardMessage := ""

	// Get and display the endpoint URL
	endpoint := shared.GetNehServerEndpoint("status")
	fmt.Printf("Checking current status...\n")
	fmt.Printf("Connecting to: %s\n", endpoint)

	// Display current language setting
	lang := getCurrentLanguage()
	fmt.Printf("Current language: %s\n", lang)

	err := shared.ExecuteWebSocketCommand("status", queryMessage, clipboardMessage)
	if err != nil {
		fmt.Println(err)
	}
}

// getCurrentLanguage gets the language from config file
func getCurrentLanguage() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "en (default)"
	}

	configPath := filepath.Join(homeDir, ".config", "neh", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "en (default)"
	}

	var cfg struct {
		Lang string `yaml:"lang"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "en (default)"
	}

	if cfg.Lang == "" {
		return "en (default)"
	}

	return cfg.Lang
}
