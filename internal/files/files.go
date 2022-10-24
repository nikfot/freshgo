package files

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

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
	if filepath.Dir(dir) != "/" {
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
	return fmt.Errorf("error: cannot remove folders in parent directory - %s", dir)
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
	cmd := exec.Command("sudo", "-V")
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

func SearchFile(dir, filename string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == filename {
			println(info.Name())
		}
		return nil
	})
	return err
}

func AddToPath(gobin string) error {
	shell, err := GetActiveShell()
	if err != nil {
		return err
	}
	confFile := getConfFile(shell)
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error: could not get homedir - %s \n", err)
	}
	err = ExportGoBin(filepath.Join(dirname, confFile), gobin)
	if err != nil {
		return err
	}
	return nil
}

func ExportGoBin(profile, gobin string) error {
	exportPath := fmt.Sprintf("export PATH=$PATH:%s", strings.TrimSpace(gobin))
	err := Prepend(profile, exportPath)
	if err != nil {
		return err
	}
	err = ReloadProfile(profile)
	if err != nil {
		return fmt.Errorf("error: could not reload file - %s %s", err, profile)
	}
	return nil
}

func Copy(source, destination string) (int64, error) {
	sourceFileStat, err := os.Stat(source)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("irregular file: %s", source)
	}
	fileData, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer fileData.Close()
	output, err := os.Create(destination)
	if err != nil {
		return 0, err
	}
	defer output.Close()
	bytesCopied, err := io.Copy(output, fileData)
	return bytesCopied, err
}

func ReloadProfile(dir string) error {
	activeShell, err := GetActiveShell()
	if err != nil {
		return err
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(strings.TrimSpace(activeShell), "-c", fmt.Sprintf("source %s", dir))
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s shell: %s command: %s \n", err, activeShell, fmt.Sprintf("source %s", dir))
	}
	return nil
}

func Prepend(file, new string) error {
	if _, err := os.Stat(file); err != nil {
		return nil
	}
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	existing := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if tmp := scanner.Text(); len(tmp) != 0 {
			existing = append(existing, tmp)
		}
	}
	if err != nil {
		return err
	}
	f.Close()
	nf, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(nf)
	writer.WriteString(fmt.Sprintf("%s\n", new))
	for _, line := range existing {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	nf.Close()
	return nil
}

func getConfFile(shell string) string {
	switch {
	case strings.Contains(shell, "zsh"):
		return ".zshrc"
	case strings.Contains(shell, "bash"):
		return ".bashrc"
	default:
		return ".profile"
	}
}

func GetActiveShell() (string, error) {
	parentpid := os.Getppid()
	proc, err := ps.FindProcess(parentpid)
	if err != nil {
		return "", err
	}
	shell := proc.Executable()
	return shell, nil
}
