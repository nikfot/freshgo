package files_test

import (
	"freshgo/internal/files"
	"strings"
	"testing"
)

func TestGetRCFiles(t *testing.T) {
	rcs, err := files.GetRCFiles()
	if len(rcs) == 0 {
		t.Errorf("TestGetRCFiles - FAIL - Neither bashrc nor zshrc found.")
	}
	if err != nil && !strings.Contains(err.Error(), "[DEBUG]") {
		t.Errorf("TestGetRCFiles - FAIL - error: %v", err)

	}
}
