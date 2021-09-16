package files

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetGoSrcPath() (dir string, err error) {
	var out bytes.Buffer
	cmd := exec.Command("which", "go")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(out.String(), "/bin/go", ""), nil
}

func BackUp(dir string, curVersion string) error {
	if _, err := os.Stat(strings.TrimSpace(dir) + strings.TrimSpace(curVersion) + "_backup"); os.IsNotExist(err) {
		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", "sudo mkdir "+strings.TrimSpace(dir)+strings.TrimSpace(curVersion)+"_backup"+" &&  sudo cp -r "+strings.TrimSpace(dir)+" "+strings.TrimSpace(dir)+strings.TrimSpace(curVersion)+"_backup")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error: file %s already exists", strings.TrimSpace(dir)+strings.TrimSpace(curVersion)+"_backup")
}

func GetFileNameFromPath(dir string) string {
	subs := strings.SplitAfter(dir, "/")
	return subs[len(subs)-1]
}

func Remove(dir string) error {
	if _, err := os.Stat(dir); os.IsExist(err) {
		var out bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", "sudo rm -rf "+dir)
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error: file %s does not exist", dir)
}
