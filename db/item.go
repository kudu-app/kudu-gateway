package db

import (
	"time"

	"github.com/knq/firebase"
)

var ItemRef *firebase.DatabaseRef

type Item struct {
	Goal    string                   `json:"goal"`
	Tag     string                   `json:"tag"`
	Created firebase.ServerTimestamp `json:"created"`
	DueDate time.Time                `json:"due_date"`
}

func (i *Item) Index(keys *map[string]interface{}) error {
	return ItemRef.Ref("/item").Get(keys)
}

func (i *Item) Add() (string, error) {
	return ItemRef.Ref("/item").Push(i)
}

func (i *Item) Get(id string, res *Item) error {
	return ItemRef.Ref("/item/" + id).Get(res)
}

func (i *Item) Set(id string, res *Item) error {
	return ItemRef.Ref("/item/" + id).Set(res)
}

func (i *Item) Update(id string, res *Item) error {
	return ItemRef.Ref("/item/" + id).Update(res)
}

func (i *Item) Delete(id string) error {
	return ItemRef.Ref("/item/" + id).Remove()
}
