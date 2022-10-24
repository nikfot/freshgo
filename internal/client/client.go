package client

import (
	"encoding/json"
	"fmt"
	gvhttp "freshgo/pkg/http"
	"io"
	"os"
	"strings"
	"time"
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

const (
	goVersionsURL = "https://go.dev/dl/?mode=json&include=all"
)

func DownloadToPath(url string, path string) error {
	cli := gvhttp.NewHTTPClient("Freshgo", "", 60*time.Second, nil, false)
	resp, err := cli.Request("GET", url, nil, "", "", nil)
	if err != nil {
		return fmt.Errorf("error: could not get versions: %v", err)
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

func GetVersions() ([]GoVersion, error) {
	var versions []GoVersion
	cli := gvhttp.NewHTTPClient("Freshgo", "", 10*time.Second, nil, false)
	resp, err := cli.Request("GET", goVersionsURL, nil, "", "", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}
