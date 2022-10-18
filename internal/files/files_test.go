package files_test

import (
	"freshgo/internal/files"
	"os"
	"testing"
)

func TestPrepend(t *testing.T) {
	profile := ".profile"
	profile_test := ".profile_test"
	dirname, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("error elevating to root: %s", err)
	}
	profile = dirname + "/" + profile
	profile_test = dirname + "/" + profile_test
	_, err = files.Copy(profile, profile_test)
	if err != nil {
		t.Errorf("could not copy /etc/profile for testing error: %s", err)
	}
	err = files.Prepend(profile_test, "test")
	if err != nil {
		t.Errorf("error editing profile file: %s error: %s", profile_test, err)
	}
	err = os.Remove(profile_test)
	if err != nil {
		t.Errorf("error removing test file: %s error: %s", profile_test, err)
	}
}
