package router

import "testing"

func TestInit(t *testing.T) {
	r := Init()
	if r == nil {
		t.Error("Expected route to not empty")
	}
}
