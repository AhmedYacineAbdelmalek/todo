/*
Copyright © 2025 Smart Todo CLI
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set by the build script
var Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of the Smart Todo CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Smart Todo CLI %s\n", Version)
		fmt.Println("Built with ❤️ for productive developers")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
