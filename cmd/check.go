package cmd

import (
	"fmt"
	"freshgo/internal/checks"

	"github.com/spf13/cobra"
)

// check runs a check on the go installation.
var check = &cobra.Command{
	Use:   "check",
	Short: "Checks your go installation.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := checks.InstallationStatus()
		if err != nil {
			fmt.Printf("found not recommended go installation - %s", err)
			return
		}
		fmt.Printf("Go Installation is %s: \n - root: %s\n - executable: %s\n - version: %s\n - os: %s\n", status.Summary, status.Root, status.Executable, status.Version, status.Runtime+"_"+status.Architecture)
	},
}

func init() {
	rootCmd.AddCommand(check)
}
