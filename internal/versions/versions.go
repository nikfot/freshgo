package versions

import (
	"fmt"
	"freshgo/internal/checks"
	"freshgo/internal/client"
	"freshgo/internal/files"

	"os"
	"path/filepath"
	"strings"

	vers "github.com/hashicorp/go-version"
)

const (
	goVersionsURL     = "https://go.dev/dl/?mode=json&include=all"
	versionTmpPathLin = "/tmp/freshgo/"
	gobin             = "/usr/local/go/bin"
	//linuxOS        = "linux-amd64"
	//versionPrefix = "<a class=\"download\" href=\"/dl/"
)

func Select(selection string, onlyNewer bool) {
	instStatus, err := checks.InstallationStatus()
	if err != nil {
		fmt.Printf("error: install status not ok - %s", err)
		return
	}
	goRoot, err := instStatus.GetGoRoot()
	if err != nil {
		fmt.Printf("error: could not get go root- %s", err)
		return
	}
	isUpgrade := true
	inPath := true
	if instStatus.Summary == "NOT FOUND" {
		isUpgrade = false
		envpath := os.Getenv("PATH")
		if !strings.Contains(envpath, gobin) {
			inPath = false
		}
	}
	versions, err := client.GetVersions()
	if err != nil {
		fmt.Printf("error: could not get versions - %s", err)
		os.Exit(1)
	}
	var versReq client.GoVersion
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
		err = files.BackUp(goRoot, instStatus.Version)
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
		return fmt.Errorf("error: install status not ok - %s", err)
	}
	goRoot, err := instStatus.GetGoRoot()
	if err != nil {
		return fmt.Errorf("error: could not get go root- %s", err)
	}
	if _, err := os.Stat(versionTmpPathLin); os.IsNotExist(err) {
		err := os.Mkdir(versionTmpPathLin, os.ModePerm)
		if err != nil {
			return err
		}
	}
	suffix, err := getCompressionSuffix()
	if err != nil {
		return fmt.Errorf("error: could not decide compression suffix - %s", err)
	}
	downloadPath := filepath.Join(instStatus.ArchivesDir, fmt.Sprintf("go%s%s", strings.TrimSuffix(version.String(), ".0"), suffix))
	fmt.Printf(" - Downloading version %s to path %s.\n", "go"+version.String(), downloadPath)
	err = Download(strings.TrimSuffix(version.String(), ".0"), downloadPath)
	if err != nil {
		return fmt.Errorf("error: could not download version - %s", err)
	}
	if isUpgrade {
		fmt.Println(" - Deleting current version.")
		err = deleteCurrentVersion(goRoot)
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
	fmt.Printf(" - Copying from %s to %s.\n", versionTmpPathLin+"go", goRoot)
	err = files.SudoCopyDir(versionTmpPathLin+"go", filepath.Dir(goRoot))
	if err != nil {
		return err
	}
	if !isUpgrade {
		err := files.ExportToPath(goRoot + "/bin")
		if err != nil {
			return err
		}
	}
	newStatus, err := checks.CurrentVersionCMD()
	if err != nil {
		return err
	}
	fmt.Println("Successfully updated go version to: ", newStatus)
	return nil
}
func DownloadVersion(version string) error {
	version = strings.TrimSuffix(version, ".0")
	instStatus, err := checks.InstallationStatus()
	if err != nil {
		return fmt.Errorf("error: install status not ok - %s", err)
	}
	suffix, err := getCompressionSuffix()
	if err != nil {
		return err
	}
	path := filepath.Join(instStatus.ArchivesDir, "go"+version+suffix)
	err = Download(version, path)
	if err != nil {
		return err
	}
	return nil
}

func Download(version, path string) error {
	exists, err := checks.ArchiveExistsAndLegit(version, path)
	if err != nil {
		return err
	}
	if !exists {
		dlVers := dlGoVersionFormat(version)
		err = client.DownloadToPath("https://go.dev"+dlVers, path)
		if err != nil {
			return err
		}
	}
	return nil
}
func List() error {
	versions, err := client.GetVersions()
	if err != nil {
		return err
	}
	for i := range versions {
		fmt.Print("• " + strings.TrimPrefix(versions[i].Name, "go") + " ")
	}
	return nil
}
func lookUpVersion(versions []client.GoVersion, name string) (client.GoVersion, error) {
	for i := range versions {
		if strings.TrimPrefix(versions[i].Name, "go") == strings.TrimPrefix(name, "go") {
			return versions[i], nil
		}
	}
	return client.GoVersion{}, fmt.Errorf("error: version '%s' notfound", name)
}

func LookUpLatest(versions []client.GoVersion, wantStable bool) (version client.GoVersion) {
	for i := range versions {
		if !wantStable {
			return versions[i]
		} else if wantStable == versions[i].Stable {
			return versions[i]
		}
	}
	return client.GoVersion{}
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

func CheckVersionArchives(version string) (bool, error) {
	instStatus, err := checks.InstallationStatus()
	if err != nil {
		return false, fmt.Errorf("error: install status not ok - %s", err)
	}
	_, err = os.Stat(instStatus.ArchivesDir)
	fmt.Println(instStatus.ArchivesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getCompressionSuffix() (string, error) {
	status, err := checks.InstallationStatus()
	if err != nil {
		return "", err
	}
	switch status.Runtime {
	case "windows":
		return ".zip", nil
	default:
		return ".tar.gz", nil
	}
}

func init() {
	if checks.OS == "" {
		checks.OS = "linux"
		checks.Architecture = "amd64"
	}
}
