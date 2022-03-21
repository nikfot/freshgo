package repos

import (
	"freshgo/internal/files"

	"github.com/hashicorp/go-version"
)

type RepoList struct {
	Repos []Repo
}
type Repo struct {
	Name    string
	Path    string
	Version version.Version
}

// List lists all go repos under a directory, along with their version.
func List(basedir string) {
	files.SearchFile(basedir, "go.mod")
}
