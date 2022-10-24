package checks_test

import (
	"freshgo/internal/checks"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveExustsAndLegit(t *testing.T) {
	version := "1.19.2"
	status, err := checks.InstallationStatus()
	if err != nil {
		t.Errorf("TestArchiveExustsAndLegit - FAIL - status error %s", err)
		return
	}
	vPath := filepath.Join(status.ArchivesDir, "go1.19.2.tar.gz")
	_, err = os.Stat(vPath)
	exists := false
	if err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("TestArchiveExustsAndLegit - FAIL - %s", err)
			return
		}
	} else {
		exists = true
	}
	suffix := ""
	switch status.Runtime {
	case "windows":
		suffix = ".zip"
	default:
		suffix = ".tar.gz"
	}
	versionPath := filepath.Join(status.ArchivesDir, "go"+version+suffix)
	found, err := checks.ArchiveExistsAndLegit(version, versionPath)
	if err != nil {
		t.Errorf("TestArchiveExustsAndLegit - FAIL - archive error %s", err)
		return
	}
	if found != exists {
		t.Errorf("TestArchiveExustsAndLegit - FAIL - expected found: %v, got found: %v", exists, found)
	}
}
