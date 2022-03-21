package cmd

import (
	"github.com/spf13/cobra"
)

// repo runs commands on local go repos.
var repo = &cobra.Command{
	Use:   "repo",
	Short: "Repo runs commands on local go repos",
	Long:  `Add comands to perform on a repo`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(repo)
}
