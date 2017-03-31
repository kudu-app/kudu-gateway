package db

import (
	"github.com/knq/firebase"
)

// Item is pointer to firebase database ref.
var Item *firebase.DatabaseRef

// Setup sets all the firebase database ref.
func Setup(path string) error {
	var err error
	Item, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsFile(path),
	)
	return err
}
