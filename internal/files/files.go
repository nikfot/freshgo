package files

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	if _, err := os.Stat(dir); !os.IsExist(err) {
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
func UnTarGz(tarGzName, xpath string) (err error) {
	gzreader, err := os.Open(tarGzName)
	defer gzreader.Close()
	if err != nil {
		return err
	}
	ungzStream, err := gzip.NewReader(gzreader)
	if err != nil {
		return err
	}
	tarStream := tar.NewReader(ungzStream)
	absPath, err := filepath.Abs(xpath)
	// untar each segment
	for {
		hdr, err := tarStream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// determine proper file path info
		finfo := hdr.FileInfo()
		fileName := hdr.Name
		absFileName := filepath.Join(absPath, fileName)
		// if a dir, create it, then go to next segment
		if finfo.Mode().IsDir() {
			if err := os.MkdirAll(absFileName, 0755); err != nil {
				return err
			}
			continue
		}
		// create new file with original file mode
		file, err := os.OpenFile(
			absFileName,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			finfo.Mode().Perm(),
		)
		if err != nil {
			return err
		}
		//fmt.Printf("x %s\n", absFileName)
		n, cpErr := io.Copy(file, tarStream)
		if closeErr := file.Close(); closeErr != nil {
			return err
		}
		if cpErr != nil {
			return cpErr
		}
		if n != finfo.Size() {
			return fmt.Errorf("wrote %d, want %d", n, finfo.Size())
		}
	}
	return nil
}

func SudoCopyDir(src, dst string) error {
	var out bytes.Buffer
	var cmd *exec.Cmd
	if sudoExists() {
		cmd = exec.Command("sudo", "cp", "-r", src, dst)
	} else {
		cmd = exec.Command("cp", "-r", src, dst)
	}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func sudoExists() bool {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", `"which sudo"`)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err == nil
}

func ExportToPath(dir string) error {
	path := os.Getenv("PATH")
	err := os.Setenv("PATH", path+":"+dir)
	if err != nil {
		return err
	}
	return nil
}
