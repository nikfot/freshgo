package versions

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	vers "github.com/hashicorp/go-version"
	"golang.org/x/net/html"
)

func PrintLatest() {
	cli := NewHTTPClient("GoVersionsURL", "", 10*time.Second, nil, false)
	resp, err := cli.Request("GET", "https://golang.org/dl/", nil, "", "", nil)
	if err != nil {
		fmt.Println("Error getting versions: ", err)
	}
	latest, err := vers.NewVersion(LookUpLatest(string(resp)))
	if err != nil {
		fmt.Println(err)
	}
	c, err := CurrentVersion()
	if err != nil {
		fmt.Println(err)
	}
	current, err := vers.NewVersion(c)
	if err != nil {
		fmt.Println(err)
	}
	comp := latest.Compare(current)
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
	if comp == 1 {
		promptUpgrade()
	}
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
				version = parseVersion(t.Attr[1].Val)
				found = true
			}
		}
	}
	return version
}

func latestVersionTag(attr []html.Attribute) bool {
	return (attr[0].Key == "class" && attr[0].Val == "download" && attr[1].Key == "href" && strings.HasPrefix(attr[1].Val, "/dl/go1"))
}

func parseVersion(value string) (version string) {
	version = strings.TrimSuffix(strings.TrimPrefix(value, "/dl/go"), ".src.tar.gz")
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

func promptUpgrade() {
	fmt.Print("Would you like to upgrade?[Y/n]")
	var prompt string
	fmt.Scanln(&prompt)
	if prompt == "Y" {
		fmt.Println("Upgrading")
	}

}

func Upgrade(latest bool, selection string) error {
	dir, err := getGoSrcPath()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	//backup()
	return nil
}

func getGoSrcPath() (dir string, err error) {
	var out bytes.Buffer
	cmd := exec.Command("which", "go")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(out.String(), "/bin/go", ""), nil
}
