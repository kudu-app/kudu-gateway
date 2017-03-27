package db

import (
	"github.com/knq/firebase"
)

var Item *firebase.DatabaseRef

func Setup(path string) error {
	var err error
	Item, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsFile(path),
	)
	return err
}
