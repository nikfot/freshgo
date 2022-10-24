package versions_test

import (
	"freshgo/internal/versions"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckVersionArchives(t *testing.T) {
	existed := false
	home, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("TestCheckVersionArchives - FAIL - %s ", err)
	}
	freshgo := filepath.Join(home, ".freshgo")
	err = os.Mkdir(freshgo, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			t.Errorf("TestCheckVersionArchives - FAIL - could not initiate test - %s ", err)
		}
		existed = true
	}
	archives := filepath.Join(freshgo, "archives")
	err = os.Mkdir(archives, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			t.Errorf("TestCheckVersionArchives - FAIL - could not initiate test - %s ", err)
		}
		existed = true
	}
	found, err := versions.CheckVersionArchives("go1.19.2")
	if err != nil {
		t.Errorf("TestCheckVersionArchives - FAIL - %s ", err)
	}
	if !found {
		t.Errorf("TestCheckVersionArchives - FAIL - expected dir to exist  but found: %v", found)
	}
	if !existed {
		err = os.RemoveAll(archives)
		if err != nil {
			t.Errorf("TestCheckVersionArchives - FAIL - %s ", err)
		}
	}
}
