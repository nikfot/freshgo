package checks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	OS           string
	Architecture string
)

type Status struct {
	Root         string
	Executable   string
	Runtime      string
	Version      string
	Architecture string
	FreshgoDir   string
}

func InstallationStatus() (status Status, err error) {
	status.Root = runtime.GOROOT()
	if !status.rootIsOK() {
		return Status{}, fmt.Errorf("error: go root set to non-recommended directory: %s", status.Root)
	}
	status.Runtime = runtime.GOOS
	status.Architecture = runtime.GOARCH
	status.Version = runtime.Version()
	status.Executable, err = status.GetGoExecutable()
	if err != nil {
		return Status{}, err
	}
	if !status.executablePathIsOK() {
		return Status{}, fmt.Errorf("error: go executable setup incorrectly in dir: %s", status.Executable)
	}
	return status, nil
}

func (s *Status) GetGoExecutable() (dir string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	path := ""
	switch strings.ToLower(s.Runtime) {
	case "windows":
	default:
		cmd := exec.Command("which", "go")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return "", err
		}
		path = strings.TrimSpace(out.String())
	}
	return path, nil
}

func (s *Status) rootIsOK() bool {
	unwantedRootDirectories := []string{"/", "/usr/local/bin", "/usr/local/sbin"}
	for i := range unwantedRootDirectories {
		if s.Root == unwantedRootDirectories[i] || s.Root == filepath.Dir(unwantedRootDirectories[i]) {
			return false
		}
	}
	return true
}

func (s *Status) executablePathIsOK() bool {
	if s.Root == filepath.Dir(s.Executable) || filepath.Dir(s.Executable) != s.Root+"/bin" {
		return false
	}
	return true
}

func CurrentVersionCMD() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(strings.Split(out.String(), " ")[2], "go"), nil
}

func FreshgoDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	freshgoDir := home + "/.freshgo"
	_, err = os.Stat(freshgoDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(freshgoDir, 644)
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}
