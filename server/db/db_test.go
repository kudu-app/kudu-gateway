package db

import (
	"os"
	"testing"
)

var testDataDir = os.Getenv("GOPATH") + "/src/github.com/rnd/kudu/server/testdata/"

func TestSetup(t *testing.T) {
	var err error

	err = Setup(testDataDir + "creds/test-datacreds.json")
	if err != nil {
		t.Error(err)
	}
}

func TestSetupWithInvalidPath(t *testing.T) {
	var err error

	err = Setup(testDataDir + "foo/bar.json")
	if err == nil {
		t.Errorf(
			"Expected setup to fail with invalid credential path: %s",
			testDataDir+"foo/bar.json")
	}
}
