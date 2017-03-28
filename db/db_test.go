package db

import (
	"os"
	"testing"
)

var TestDataDir = os.Getenv("GOPATH") + "/src/github.com/rnd/kudu/testdata/"

func TestSetup(t *testing.T) {
	var err error

	err = Setup(TestDataDir + "creds/test-datacreds.json")
	if err != nil {
		t.Error(err)
	}
}

func TestSetupWithInvalidPath(t *testing.T) {
	var err error

	err = Setup(TestDataDir + "foo/bar.json")
	if err == nil {
		t.Errorf(
			"Expected setup to fail with invalid credential path: %s",
			TestDataDir+"foo/bar.json")
	}
}
