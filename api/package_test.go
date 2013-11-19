package api

import (
	"strings"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestPackageList(t *testing.T) {
	request, _ := http.NewRequest("GET", "test/package", nil)
	response := httptest.NewRecorder()

	packageList(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response code expected %v got %v:\n",
		http.StatusOK, response.Code)
	}

	ct := response.HeaderMap["Content-Type"][0]
	if !strings.EqualFold(ct, "application/json") {
		t.Fatalf("Content-Type does not equal 'application/json'")
	}
}

func TestPackageUpload(t *testing.T) {
	request, _ := http.NewRequest("POST", "test/package", nil)
	response := httptest.NewRecorder()

	packageUpload(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("Response code expected %v got %v:\n",
		http.StatusCreated, response.Code)
	}

	ct := response.HeaderMap["Content-Type"][0]
	if !strings.EqualFold(ct, "application/json") {
		t.Fatalf("Content-Type does not equal 'application/json'")
	}
}

func TestPackageDownload(t *testing.T) {
	request, _ := http.NewRequest("GET", "test/package/pkg.tgz", nil)
	response := httptest.NewRecorder()

	packageDownload(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response code expected %v got %v:\n",
		http.StatusOK, response.Code)
	}

	ct := response.HeaderMap["Content-Type"][0]
	if !strings.EqualFold(ct, "application/json") {
		t.Fatalf("Content-Type does not equal 'application/json'")
	}
}
