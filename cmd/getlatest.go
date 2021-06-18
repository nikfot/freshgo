package cmd

import (
	"go-versions/versions"

	"github.com/spf13/cobra"
)

// getLatest gets the latest version
var getLatest = &cobra.Command{
	Use:   "latest",
	Short: "Get the latest go version semantic version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		versions.PrintLatest()
	},
}

func init() {
	rootCmd.AddCommand(getLatest)
}
