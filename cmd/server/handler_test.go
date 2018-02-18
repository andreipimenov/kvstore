package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := NewRouter(&Config{Port: 8080}, nil)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	expected := `{"message":"pong"}`
	if rr.Body.String() != expected {
		t.Errorf("Wrong body: got %v, expected %v", rr.Body.String(), expected)
	}
}
