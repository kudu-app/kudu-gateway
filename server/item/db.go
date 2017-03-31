package item

import (
	"github.com/knq/firebase"
	"github.com/rnd/kudu/server/db"
)

// Item represents firebase database model for /item ref.
type Item struct {
	Goal    string                   `json:"goal"`
	Tag     string                   `json:"tag"`
	Notes   string                   `json:"notes"`
	Created firebase.ServerTimestamp `json:"created"`
}

// Keys retreive all /item keys.
func (i *Item) Keys(keys *map[string]interface{}) error {
	return db.Item.Ref("/item").Get(keys, firebase.Shallow)
}

// Index retreive all /item values.
func (i *Item) Index(keys *map[string]interface{}) error {
	return db.Item.Ref("/item").Get(keys)
}

// Add push new item to /item ref.
func (i *Item) Add() (string, error) {
	return db.Item.Ref("/item").Push(i)
}

// Get retreive item which have specified id.
func (i *Item) Get(id string, res *Item) error {
	return db.Item.Ref("/item/" + id).Get(res)
}

// Update updates item values which have specified id.
func (i *Item) Update(id string, res *Item) error {
	return db.Item.Ref("/item/" + id).Update(res)
}

// Delete remove item which have specified id.
func (i *Item) Delete(id string) error {
	return db.Item.Ref("/item/" + id).Remove()
}
