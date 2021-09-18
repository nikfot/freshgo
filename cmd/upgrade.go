package cmd

import (
	"fmt"
	"freshgo/internal/versions"

	"github.com/spf13/cobra"
)

// upgradeVersion represents the sum command
var upgradeVersion = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade current go version",
	Long:  `You can upgrade to a selected version or the latest version vailable.`,
	Run: func(cmd *cobra.Command, args []string) {
		latest, err := cmd.Flags().GetBool("latest")
		if err != nil {
			fmt.Println("Error parsing app name flag: ", err)
		}
		selection, err := cmd.Flags().GetString("select")
		if err != nil {
			fmt.Println("Error parsing app name flag: ", err)
		}
		// versions.UpgradeVersion(latest)
		fmt.Println(latest, selection)
		err = versions.Upgrade(latest, selection)
		if err != nil {
			fmt.Println("Error upgrading: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeVersion)
	// Here you will define your flags and configuration settings.
	upgradeVersion.Flags().BoolP("latest", "l", false, "Select upgrade to latest version.")
	upgradeVersion.Flags().StringP("select", "s", "", "Select the semantic go version to upgrade to.")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sumCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sumCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
