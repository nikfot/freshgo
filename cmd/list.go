package cmd

import (
	"freshgo/internal/versions"

	"github.com/spf13/cobra"
)

// list lists all versions
var list = &cobra.Command{
	Use:   "list",
	Short: "List go versions.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		versions.List()
	},
}

func init() {
	rootCmd.AddCommand(list)
}
