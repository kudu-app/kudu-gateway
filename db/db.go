package db

import (
	"github.com/knq/firebase"
)

var Item *firebase.DatabaseRef

func Setup(credPath string) error {
	var err error
	Item, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsFile(credPath),
	)
	return err
}
