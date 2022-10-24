package checks

import (
	"bytes"
	"os/exec"
	"strings"
)

func GetGoSrcPath(OS string) (dir string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	path := ""
	switch strings.ToLower(OS) {
	case "windows":
	default:
		cmd := exec.Command("which", "go")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return "", err
		}
		path = strings.TrimSpace(strings.ReplaceAll(out.String(), "/bin/go", ""))
		if path == "" {
			return "", nil
		}
	}
	return path, nil
}
