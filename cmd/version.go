/*
* filename: cmd/version.go
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Define the version information
const Version = "0.0.35"

// versionCmd represents the `neh version` command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version information",
	Long: `The version command displays the current version of the neh command.
It provides detailed information about the version number, which can be useful for troubleshooting and ensuring compatibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version of neh command: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
