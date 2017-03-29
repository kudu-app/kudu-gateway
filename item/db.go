package item

import (
	"github.com/knq/firebase"
	"github.com/rnd/kudu/db"
)

type Item struct {
	Goal    string                   `json:"goal"`
	Tag     string                   `json:"tag"`
	Notes   string                   `json:"notes"`
	Created firebase.ServerTimestamp `json:"created"`
}

func (i *Item) Keys(keys *map[string]interface{}) error {
	return db.Item.Ref("/item").Get(keys, firebase.Shallow)
}

func (i *Item) Index(keys *map[string]interface{}) error {
	return db.Item.Ref("/item").Get(keys)
}

func (i *Item) Add() (string, error) {
	return db.Item.Ref("/item").Push(i)
}

func (i *Item) Get(id string, res *Item) error {
	return db.Item.Ref("/item/" + id).Get(res)
}

func (i *Item) Set(id string, res *Item) error {
	return db.Item.Ref("/item/" + id).Set(res)
}

func (i *Item) Update(id string, res *Item) error {
	return db.Item.Ref("/item/" + id).Update(res)
}

func (i *Item) Delete(id string) error {
	return db.Item.Ref("/item/" + id).Remove()
}
