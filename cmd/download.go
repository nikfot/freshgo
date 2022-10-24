package cmd

import (
	"fmt"
	"freshgo/internal/versions"
	"os"

	"github.com/spf13/cobra"
)

// download downloads a version for later use.
var download = &cobra.Command{
	Use:   "download",
	Short: "Downloads a version for later use",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Printf("error: no valid version number supplied - %s", err)
			os.Exit(1)
		}
		err = versions.DownloadVersion(version)
		if err != nil {
			fmt.Printf("error: could not download version %s - %s", version, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully downloaded version %s in archives path.", version)
	},
}

func init() {
	rootCmd.AddCommand(download)
	download.Flags().StringP("version", "v", "latest", "Select the semantic go version to install.")
}
