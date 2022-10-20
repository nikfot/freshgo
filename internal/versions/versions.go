package versions

import (
	"encoding/json"
	"fmt"
	"freshgo/internal/checks"
	"freshgo/internal/files"
	gvhttp "freshgo/pkg/http"
	"io"
	"os"
	"strings"
	"time"

	vers "github.com/hashicorp/go-version"
)

const (
	versionTmpPathLin = "/tmp/freshgo/"
	gobin             = "/usr/local/go/bin"
	//linuxOS        = "linux-amd64"
	//versionPrefix = "<a class=\"download\" href=\"/dl/"
)

type GoVersion struct {
	Name   string `json:"version"`
	Stable bool   `json:"stable"`
	Files  []File `json:"files"`
}

type File struct {
	Name    string `json:"filename"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
	Size    int    `json:"size"`
	Kind    string `json:"kind"`
}

func Select(selection string, onlyNewer bool) {
	instStatus, err := checks.InstallationStatus()
	if err != nil {
		fmt.Printf("error: install status not ok - %s", err)
	}
	isUpgrade := true
	inPath := true
	if err != nil {
		isUpgrade = false
		envpath := os.Getenv("PATH")
		if !strings.Contains(envpath, gobin) {
			inPath = false
		}
	}
	versions, err := getVersions()
	if err != nil {
		fmt.Println("error: getting versions failed - ", err)
	}
	var versReq GoVersion
	if selection == "latest" {
		versReq = LookUpLatest(versions, true)
	} else {
		versReq, err = lookUpVersion(versions, selection)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	semVer, err := vers.NewVersion(strings.TrimPrefix(versReq.Name, "go"))
	if err != nil {
		fmt.Println("error: could not get latest version.")
		return
	}
	newer, isUpgrade, err := compare(semVer)
	if err != nil {
		fmt.Println("error: could not compare versions - ", err)
	}
	if onlyNewer && !newer {
		return
	}
	if !promptInstall(isUpgrade) {
		return
	}
	if promptBackup(isUpgrade) {
		if err != nil {
			fmt.Println("Error getting go bin dir: ", err)
			return
		}
		err = files.BackUp(instStatus.Root, instStatus.Version)
		if err != nil {
			fmt.Println("Error taking backup: ", err)
			return
		}
	}
	err = InstallVersion(semVer, isUpgrade)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !inPath {
		fmt.Println(" - Go binary not found in path. Adding export to shell config file.")
		err = files.AddToPath(gobin)
		if err != nil {
			fmt.Printf("error: could not update profile - %s \n", err)
		}
	}
	fmt.Println("[FRESHGO, OK TO GO!]")
}

func Latest() {
	Select("latest", true)
}

func InstallVersion(version *vers.Version, isUpgrade bool) error {
	instStatus, err := checks.InstallationStatus()
	if err != nil {
		fmt.Printf("error: install status not ok - %s", err)
	}
	defaultGoSrcPath := "/usr/local/go"
	if _, err := os.Stat(instStatus.Root); err == nil {
		os.MkdirAll(defaultGoSrcPath, os.ModePerm)
	}
	if _, err := os.Stat(versionTmpPathLin); os.IsNotExist(err) {
		err := os.Mkdir(versionTmpPathLin, os.ModePerm)
		if err != nil {
			return err
		}
	}
	downloadPath := versionTmpPathLin + "go" + version.String()
	fmt.Printf(" - Downloading version %s to path %s.\n", "go"+version.String(), downloadPath)
	dlVers := dlGoVersionFormat(version.String())
	err = downloadToPath("https://go.dev"+dlVers, downloadPath)
	if err != nil {
		return err
	}
	if isUpgrade {
		fmt.Println(" - Deleting current version.")
		err = deleteCurrentVersion(instStatus.Root)
		if err != nil {
			return err
		}
	}
	fmt.Printf(" - Untaring downloaded version from %s to %s.\n", downloadPath, versionTmpPathLin)
	err = files.UnTarGz(downloadPath, versionTmpPathLin)
	// err = otiai10.Copy(downloadPath, curGoSrcPath)
	if err != nil {
		return err
	}
	fmt.Printf(" - Copying from %s to %s.\n", versionTmpPathLin+"go", instStatus.Root)
	err = files.SudoCopyDir(versionTmpPathLin+"go", instStatus.Root)
	if err != nil {
		return err
	}
	if !isUpgrade {
		err := files.ExportToPath(instStatus.Root + "/bin")
		if err != nil {
			return err
		}
	}
	newStatus, err := checks.InstallationStatus()
	if err != nil {
		return err
	}
	fmt.Println("Successfully updated go version to: ", newStatus.Version)
	return nil
}

func List() error {
	versions, err := getVersions()
	if err != nil {
		return err
	}
	for i := range versions {
		fmt.Print("• " + strings.TrimPrefix(versions[i].Name, "go") + " ")
	}
	return nil
}
func lookUpVersion(versions []GoVersion, name string) (GoVersion, error) {
	for i := range versions {
		if strings.TrimPrefix(versions[i].Name, "go") == strings.TrimPrefix(name, "go") {
			return versions[i], nil
		}
	}
	return GoVersion{}, fmt.Errorf("error: version '%s' notfound", name)
}
func downloadToPath(url string, path string) error {
	cli := gvhttp.NewHTTPClient("Freshgo", "", 60*time.Second, nil, false)
	resp, err := cli.Request("GET", url, nil, "", "", nil)
	if err != nil {
		return fmt.Errorf("error getting versions: %v", err)
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, strings.NewReader(string(resp)))
	return err
}

func LookUpLatest(versions []GoVersion, wantStable bool) (version GoVersion) {
	for i := range versions {
		if !wantStable {
			return versions[i]
		} else if wantStable == versions[i].Stable {
			return versions[i]
		}
	}
	return GoVersion{}
}

func promptInstall(upgrade bool) bool {
	if upgrade {
		fmt.Print("Would you like to upgrade?[Y/n]")
	} else {
		fmt.Print("Would you like to install?[Y/n]")
	}
	var prompt string
	fmt.Scanln(&prompt)
	return prompt == "Y"
}

func promptBackup(upgrade bool) bool {
	if upgrade {
		fmt.Print("Would you like to backup current go-version[Y/n]")
		var prompt string
		fmt.Scanln(&prompt)
		return prompt == "Y"
	}
	return false
}
func getVersions() ([]GoVersion, error) {
	cli := gvhttp.NewHTTPClient("Freshgo", "", 10*time.Second, nil, false)
	resp, err := cli.Request("GET", "https://go.dev/dl/?mode=json&include=all", nil, "", "", nil)
	if err != nil {
		return nil, err
	}
	var versions []GoVersion
	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func compare(upstream *vers.Version) (newer, isUpgrade bool, err error) {
	comp := 1
	c, err := checks.CurrentVersionCMD()
	if err != nil {
		fmt.Println("[INFO]: no installed go version.")
		return true, false, nil
	} else {
		isUpgrade = true
		current, err := vers.NewVersion(c)
		if err != nil {
			fmt.Println(err)
		}
		comp = upstream.Compare(current)
		is := ""
		switch comp {
		case -1:
			is = "older than"
			newer = false
		case 1:
			is = "newer than"
			newer = true
		default:
			is = "equal to"
			newer = false
		}
		fmt.Printf("The latest go version is %v, which is %s the current %v \n", upstream.String(), is, current.String())
	}
	return newer, isUpgrade, err
}

func deleteCurrentVersion(root string) error {
	err := files.Remove(strings.TrimSpace(root))
	if err != nil {
		return fmt.Errorf("%s", "error: no go src path exists")
	}
	return nil
}

func dlGoVersionFormat(version string) string {
	// in case this is the first minor version ie 1.18
	version = strings.TrimSuffix(version, ".0")
	switch strings.ToLower(checks.OS) {
	case "windows":
		return "/dl/go" + version + "." + strings.ToLower(checks.OS) + "-" + strings.ToLower(checks.Architecture) + ".zip"
	default:
		return "/dl/go" + version + "." + strings.ToLower(checks.OS) + "-" + strings.ToLower(checks.Architecture) + ".tar.gz"
	}
}

func init() {
	if checks.OS == "" {
		checks.OS = "linux"
		checks.Architecture = "amd64"
	}
}
