package cmd

import (
	"go-versions/internal/versions"

	"github.com/spf13/cobra"
)

// getLatest gets the latest version
var getLatest = &cobra.Command{
	Use:   "latest",
	Short: "Get the latest go version semantic version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		versions.Latest()
	},
}

func init() {
	rootCmd.AddCommand(getLatest)
}
