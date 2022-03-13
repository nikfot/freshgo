package cmd

import (
	"fmt"
	"freshgo/internal/versions"

	"github.com/spf13/cobra"
)

// selectVersion selects specific version
var selectVersion = &cobra.Command{
	Use:   "select",
	Short: "Select go version semantic version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		selection, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Println("Error parsing app name flag: ", err)
		}
		versions.Select(selection, false)
	},
}

func init() {
	rootCmd.AddCommand(selectVersion)
	selectVersion.Flags().StringP("version", "v", "latest", "Select the semantic go version to install.")
}
