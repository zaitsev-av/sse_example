package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetConnect(t *testing.T) {
	req, err := http.NewRequest("GET", "/connect", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getConnect)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Server is up and running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleObjects(t *testing.T) {
	jsonStr := []byte(`[{"id":1, "status":"pending"}, {"id":2, "status":"in_progress"}]`)
	req, err := http.NewRequest("POST", "/objects", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleObjects)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleSSE(t *testing.T) {
	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSSE)
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Header().Get("Content-Type"), "text/event-stream") {
		t.Errorf("handler returned wrong content type: got %v", rr.Header().Get("Content-Type"))
	}
}
