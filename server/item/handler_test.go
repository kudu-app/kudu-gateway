package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Index)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("Expected response to have status %d : got %d",
			http.StatusOK, status)
	}

	if contentType := http.DetectContentType(res.Body.Bytes()); strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to have content-type %s : got %s",
			"application/json", contentType)
	}
}

func TestPostHandler(t *testing.T) {
	var err error

	testData := Item{
		Goal:  "Foo",
		Tag:   "Bar",
		Notes: "Baz",
	}
	jsonData, err := json.Marshal(testData)

	req, err := http.NewRequest("POST", "/items", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Expected response to have status %d : got %d",
			http.StatusCreated, status)
	}
	if contentType := http.DetectContentType(res.Body.Bytes()); strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to have content-type %s : got %s",
			"application/json", contentType)
	}
}

func TestPostHandlerWithInvalidData(t *testing.T) {
	var err error

	testData := Item{
		Goal:  "",
		Tag:   "Bar",
		Notes: "Baz",
	}
	jsonData, err := json.Marshal(testData)

	req, err := http.NewRequest("POST", "/items", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusBadRequest {
		t.Errorf("Expected response to have status %d : got %d",
			http.StatusBadRequest, status)
	}
	if contentType := http.DetectContentType(res.Body.Bytes()); strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to have content-type %s : got %s",
			"application/json", contentType)
	}
}

func TestGetHandler(t *testing.T) {
	urlstr := fmt.Sprintf("/items/%s", "FOO")
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Get)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusNotFound {
		t.Errorf("Expected response to have status %d : got %d",
			http.StatusNotFound, status)
	}
	if contentType := http.DetectContentType(res.Body.Bytes()); strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to have content-type %s : got %s",
			"application/json", contentType)
	}
}
