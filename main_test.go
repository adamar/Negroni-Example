package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimplePage(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	SimplePage(resp, req, "mainpage")

	if resp.Code != http.StatusOK {
		t.Error("Response Code Incorrect ", resp.Code)
	}

}
