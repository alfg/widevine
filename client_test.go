package widevine

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testResponse struct {
	Test string `json:"test"`
}

var server *httptest.Server

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()

	teardown()
	os.Exit(retCode)
}

func teardown() {
	defer server.Close()
}

func setup() {
	// Mock server.
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"test": "testing123"}`)
	}))
}

func TestGet(t *testing.T) {

	// Test client.
	resp := testResponse{}
	client, _ := NewClient()
	err := client.get(server.URL, &resp)
	if err != nil {
		t.Error()
	}

	if resp.Test != "testing123" {
		t.Error()
	}
}

func TestPost(t *testing.T) {

	// Test client.
	resp := testResponse{}
	client, _ := NewClient()
	err := client.post(server.URL, &resp, map[string]interface{}{})
	if err != nil {
		t.Error()
	}

	if resp.Test != "testing123" {
		t.Error()
	}
}
