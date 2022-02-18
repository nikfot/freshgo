package versions

import (
	"bytes"
	"fmt"
	"freshgo/internal/files"
	gvhttp "freshgo/pkg/http"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	vers "github.com/hashicorp/go-version"
	"golang.org/x/net/html"
)

var (
	OS           string
	Architecture string
)

const (
	versionTmpPathLin = "/tmp/freshgo/"

	//linuxOS        = "linux-amd64"
	//versionPrefix = "<a class=\"download\" href=\"/dl/"
)

func Select(selection string) {
	_, err := CurrentVersion()
	isUpgrade := true
	if err != nil {
		isUpgrade = false
	}
	if selection == "latest" {
		cli := gvhttp.NewHTTPClient("GoVersionsURL", "", 10*time.Second, nil, false)
		resp, err := cli.Request("GET", "https://go.dev/dl/", nil, "", "", nil)
		if err != nil {
			fmt.Println("Error getting versions: ", err)
		}
		latest, err := vers.NewVersion(LookUpLatest(string(resp)))
		if err != nil {
			fmt.Println(err)
			return
		}
		err = InstallVersion(latest, isUpgrade)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	selectVers, err := vers.NewVersion(selection)
	if err != nil {
		fmt.Println(err)
	}
	err = InstallVersion(selectVers, isUpgrade)
	if err != nil {
		fmt.Println(err)
	}
}
func Latest() {
	cli := gvhttp.NewHTTPClient("GoVersionsURL", "", 10*time.Second, nil, false)
	resp, err := cli.Request("GET", "https://go.dev/dl/", nil, "", "", nil)
	if err != nil {
		fmt.Println("Error getting versions: ", err)
	}
	latest, err := vers.NewVersion(LookUpLatest(string(resp)))
	if err != nil {
		fmt.Println("error: could not get latest version.")
		return
	}
	comp := 1
	current := &vers.Version{}
	c, err := CurrentVersion()
	isUpgrade := false
	if err != nil {
		fmt.Println("[INFO]: no installed go version.")
	} else {
		isUpgrade = true
		current, err = vers.NewVersion(c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(current)
		comp = latest.Compare(current)
		is := ""
		switch comp {
		case -1:
			is = "older than"
		case 1:
			is = "newer than"
		default:
			is = "equal to"
		}
		fmt.Printf("The latest go version is %v, which is %s the current %v \n", latest.String(), is, current.String())
	}
	if comp == 1 {
		if !promptInstall(isUpgrade) {
			return
		}
		if promptBackup(isUpgrade) {
			curDir, err := files.GetGoSrcPath(OS)
			if err != nil {
				fmt.Println("Error getting go bin dir: ", err)
				return
			}
			err = files.BackUp(curDir, current.String())
			if err != nil {
				fmt.Println("Error taking backup: ", err)
				return
			}

		}
		err := InstallVersion(latest, isUpgrade)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InstallVersion(version *vers.Version, isUpgrade bool) error {
	curGoSrcPath := "/usr/local/go"
	if _, err := os.Stat(curGoSrcPath); err == nil {
		os.MkdirAll(curGoSrcPath, os.ModePerm)
	}
	if _, err := os.Stat(versionTmpPathLin); os.IsNotExist(err) {
		err := os.Mkdir(versionTmpPathLin, os.ModePerm)
		if err != nil {
			return err
		}
	}
	downloadPath := versionTmpPathLin + files.GetFileNameFromPath(version.String())
	fmt.Printf(" - Downloading version %s to path %s.\n", version, downloadPath)
	dlVers := dlGoVersionFormat(version.String())
	err := downloadToPath("https://go.dev"+dlVers, downloadPath)
	if err != nil {
		return err
	}
	if isUpgrade {
		curGoSrcPath, err = files.GetGoSrcPath(OS)
		if err != nil {
			return err
		}
		fmt.Println(" - Deleting current version.")
		err = deleteCurrentVersion()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	fmt.Printf(" - Untaring downloaded version from %s to %s.\n", downloadPath, versionTmpPathLin)
	err = files.UnTarGz(downloadPath, versionTmpPathLin)
	// err = otiai10.Copy(downloadPath, curGoSrcPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf(" - Copying from %s to %s.\n", versionTmpPathLin, curGoSrcPath)
	err = files.SudoCopyDir(versionTmpPathLin+"/go", curGoSrcPath)
	if err != nil {
		fmt.Println(err)
	}
	if !isUpgrade {
		err := files.ExportToPath(curGoSrcPath + "/bin")
		if err != nil {
			return err
		}
	}
	u, err := CurrentVersion()
	if err != nil {
		fmt.Println(err)
	}
	updated, err := vers.NewVersion(u)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully updated go version to: ", updated)
	return nil
}

func List() error {
	cli := gvhttp.NewHTTPClient("GoVersionsURL", "", 10*time.Second, nil, false)
	resp, err := cli.Request("GET", "https://go.dev/dl/", nil, "", "", nil)
	if err != nil {
		fmt.Println("Error getting versions: ", err)
	}
	z := html.NewTokenizer(strings.NewReader(string(resp)))
	for z.Next() != html.ErrorToken {
		tt := z.Next()
		found := false
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return fmt.Errorf("could not find latest %s-%s version metadata", OS, Architecture)
		case tt == html.StartTagToken:
			t := z.Token()
			if len(t.Attr) > 1 {
				for i := range t.Attr {
					if t.Attr[i].Key == "id" && strings.Contains(t.Attr[i].Val, "go") {
						fmt.Print("• " + strings.TrimPrefix(t.Attr[i].Val, "go") + " ")
						found = true
					}
				}
			}
			if t.Data == "a" && versionTag(t.Attr, strings.ToLower(OS+"-"+Architecture)) && !found {
				fmt.Print("• " + parseVersion(t.Attr[1].Val) + " ")
			}
		}
	}
	return nil
}

func downloadToPath(url string, path string) error {
	cli := gvhttp.NewHTTPClient("GoDownload", "", 60*time.Second, nil, false)
	resp, err := cli.Request("GET", url, nil, "", "", nil)
	if err != nil {
		fmt.Println("Error getting versions: ", err)
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

func LookUpLatest(body string) (version string) {
	z := html.NewTokenizer(strings.NewReader(body))
	found := false
	for !found {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return ""
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" && latestVersionTag(t.Attr) {
				version = parseVersionFile(t.Attr[1].Val)
				found = true
			}
		}
	}
	return version
}

func latestVersionTag(attr []html.Attribute) bool {
	return (attr[0].Key == "class" && attr[0].Val == "download" && attr[1].Key == "href" && strings.HasPrefix(attr[1].Val, "/dl/go1"))
}
func versionTag(attr []html.Attribute, system string) bool {
	return (attr[0].Key == "class" && attr[0].Val == "download" && attr[1].Key == "href" && strings.Contains(attr[1].Val, system))
}
func parseVersionFile(value string) (version string) {
	switch strings.ToLower(OS) {
	case "windows":
		version = strings.TrimSuffix(strings.TrimPrefix(value, "/dl/go"), ".src.zip")
	default:
		version = strings.TrimSuffix(strings.TrimPrefix(value, "/dl/go"), ".src.tar.gz")
	}
	return version
}
func parseVersion(value string) (version string) {
	switch strings.ToLower(OS) {
	case "windows":
		version = strings.TrimSuffix(strings.TrimPrefix(value, "/dl/go"), "."+strings.ToLower(OS+"-"+Architecture)+".zip")
	default:
		version = strings.TrimSuffix(strings.TrimPrefix(value, "/dl/go"), "."+strings.ToLower(OS+"-"+Architecture)+".tar.gz")
	}
	return version
}
func CurrentVersion() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(strings.Split(out.String(), " ")[2], "go"), nil
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

func deleteCurrentVersion() error {
	curDir, err := files.GetGoSrcPath(OS)
	if err != nil {
		return err
	}
	if curDir != "" {
		err = files.Remove(strings.TrimSpace(curDir))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s", "error: no go src path exists")
	}
	return nil
}

func dlGoVersionFormat(version string) string {
	switch strings.ToLower(OS) {
	case "windows":
		return "/dl/go" + version + "." + strings.ToLower(OS) + "-" + strings.ToLower(Architecture) + ".zip"
	default:
		return "/dl/go" + version + "." + strings.ToLower(OS) + "-" + strings.ToLower(Architecture) + ".tar.gz"
	}
}

func init() {
	if OS == "" {
		OS = "linux"
		Architecture = "amd64"
	}
}
