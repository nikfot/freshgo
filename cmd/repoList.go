package cmd

import (
	"freshgo/internal/repos"

	"github.com/spf13/cobra"
)

// repoList lists all go repos under a dir.
var repoList = &cobra.Command{
	Use:   "list",
	Short: "Lists go repos under a dir",
	Long:  `List directory and go version of go repos`,
	Run: func(cmd *cobra.Command, args []string) {
		basedir, err := cmd.Flags().GetString("basedir")
		if err != nil {
			return
		}
		repos.List(basedir)
	},
}

func init() {
	repoList.Flags().StringP("basedir", "d", "./", "Choose which directory to look into for go repos.")
	repo.AddCommand(repoList)
}
