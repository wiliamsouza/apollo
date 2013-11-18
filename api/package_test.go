package api_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestListPackage(t *testing.T) {
	request, _ := http.NewRequest("GET", "test/package", nil)
	response := httptest.NewRecorder()

	listPackage(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response code expected 200 got %v:\n", response.Code)
	}
}

func TestUploadPackage(t *testing.T) {
	request, _ := http.NewRequest("POST", "test/package", nil)
	response := httptest.NewRecorder()

	uploadPackage(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response code expected 200 got %v:\n", response.Code)
	}
}

func TestDownloadPackage(t *testing.T) {
	request, _ := http.NewRequest("GET", "test/package/pkg.tgz", nil)
	response := httptest.NewRecorder()

	downloadPackage(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Response code expected 200 got %v:\n", response.Code)
	}
}
