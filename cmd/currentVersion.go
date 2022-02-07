package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version outputs the version of freshgo.
var version = &cobra.Command{
	Use:   "version",
	Short: "Show current version of freshgo.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(version)
}
