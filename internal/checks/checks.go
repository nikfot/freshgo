package checks

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"freshgo/internal/client"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	archivesFolder = "archives"
	freshgoFolder  = ".freshgo"
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
	ArchivesDir  string
	Summary      string
}

func InstallationStatus() (status Status, err error) {
	status.FreshgoDir, err = FreshgoDir()
	if err != nil {
		return Status{}, fmt.Errorf("error: %s", status.Root)
	}
	status.ArchivesDir, err = ArchivesDir()
	if err != nil {
		return Status{}, fmt.Errorf("error: %s", status.Root)
	}
	status.Version, err = CurrentVersionCMD()
	if err != nil {
		status.Version = "0.0.0"
		status.Summary = "NOT FOUND"
		return status, nil
	}
	status.Root = runtime.GOROOT()
	if !status.rootIsOK() {
		status.Summary = "NOT OK"
		return Status{}, fmt.Errorf("error: go root set to non-recommended directory: %s", status.Root)
	}
	status.Runtime = runtime.GOOS
	status.Architecture = runtime.GOARCH
	status.Executable, err = status.GetGoExecutable()
	if err != nil {
		status.Summary = "NOT OK"
		return Status{}, err
	}
	if !status.executablePathIsOK() {
		status.Summary = "NOT OK"
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

func (s *Status) GetGoRoot() (string, error) {
	defaultGoSrcPath := "/usr/local/go"
	if strings.TrimSpace(s.Root) == "" {
		_, err := os.Stat(defaultGoSrcPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(defaultGoSrcPath, 644)
				if err != nil {
					return "/unknown", err
				}
			} else {
				return "/unknown", err
			}
		}
		return defaultGoSrcPath, nil
	}
	return s.Root, nil
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

func FreshgoDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/unknown", err
	}
	freshgoDir := filepath.Join(home, freshgoFolder)
	_, err = os.Stat(freshgoDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(freshgoDir, 644)
			if err != nil {
				return "/unknown", err
			}
		} else {
			return "/unknown", err
		}
	}
	return freshgoDir, nil
}
func ArchivesDir() (string, error) {
	freshgoDir, err := FreshgoDir()
	if err != nil {
		return "/unknown", err
	}
	archivesDir := filepath.Join(freshgoDir, archivesFolder)
	_, err = os.Stat(archivesDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(archivesDir, 644)
			if err != nil {
				return "/unknown", err
			}
		} else {
			return "/unknown", err
		}
	}
	return archivesDir, nil
}

func ArchiveExistsAndLegit(version, versionPath string) (bool, error) {
	status, err := InstallationStatus()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	versionPack, err := os.Open(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer versionPack.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, versionPack); err != nil {
		log.Fatal(err)
	}

	versions, err := client.GetVersions()
	if err != nil {
		return false, fmt.Errorf("%s - cannot check if archive exists and legit", err)
	}
	for i := range versions {
		if versions[i].Name == "go"+version {
			for j := range versions[i].Files {
				if versions[i].Files[j].OS == status.Runtime && versions[i].Files[j].Arch == status.Architecture {
					return versions[i].Files[j].SHA256 == fmt.Sprintf("%x", hash.Sum(nil)), nil
				}
			}
		}
	}
	return false, fmt.Errorf("error: version SHA not found")
}
