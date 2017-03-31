package item

import (
	"log"
	"os"
	"testing"

	"github.com/rnd/kudu/server/db"
)

var (
	testItem    Item
	testDataDir string
)

func testClear() {
	var err error

	keys := make(map[string]interface{})
	err = testItem.Keys(&keys)
	if err != nil {
		log.Fatal(err)
	}

	for key, _ := range keys {
		err = testItem.Delete(key)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	testDataDir = os.Getenv("GOPATH") + "/src/github.com/rnd/kudu/server/testdata/"
	db.Setup(testDataDir + "creds/test-datacreds.json")
	code := m.Run()
	testClear()
	os.Exit(code)
}

func TestIndex(t *testing.T) {
	var err error

	res := make(map[string]interface{})
	err = testItem.Index(&res)
	if err != nil {
		t.Error(err)
	}
}

func TestAdd(t *testing.T) {
	tests := []Item{
		{
			Goal:  "Foo",
			Tag:   "",
			Notes: "",
		},
		{
			Goal:  "Foo",
			Tag:   "Bar",
			Notes: "Baz",
		},
	}

	for i, test := range tests {
		key, err := test.Add()
		if err != nil {
			t.Error(err)
		}
		if key == "" {
			t.Errorf("test %d expected to have not empty key", i)
		}
	}
}

func TestGet(t *testing.T) {
	var err error
	var res Item

	item = Item{
		Goal:  "Foo",
		Tag:   "Bar",
		Notes: "Baz",
	}

	key, err := item.Add()
	if err != nil {
		t.Error(err)
	}

	err = item.Get(key, &res)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	var err error

	item = Item{
		Goal:  "Foo",
		Tag:   "Bar",
		Notes: "Baz",
	}

	key, err := item.Add()
	if err != nil {
		t.Error(err)
	}

	item.Goal = "FooBar"
	err = item.Update(key, &item)
	if err != nil {
		t.Error(err)
	}
	if item.Goal != "FooBar" {
		t.Errorf("Expected goal value to be 'FooBar' got: %s", item.Goal)
	}
}

func TestDelete(t *testing.T) {
	var err error

	item = Item{
		Goal:  "Foo",
		Tag:   "Bar",
		Notes: "Baz",
	}

	key, err := item.Add()
	if err != nil {
		t.Error(err)
	}

	err = item.Delete(key)
	if err != nil {
		t.Error(err)
	}
}
