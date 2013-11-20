package api

import (
	"strings"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestPackageList(t *testing.T) {
	results := `[{"_id" : ObjectId("5287a2b1e5adc4f8e90b50a2"),
	           "filename" : "pkg-1.0.tar.gz",
		   "chunkSize" : 262144,
		   "uploadDate" : ISODate("2013-11-16T16:52:01.089Z"),
		   "md5" : "ac98f052e55d0c556e6536796d6b53eb",
		   "length" : 2133},
	           {"_id" : ObjectId("5287a2b1e5adc4f8e90b50a3"),
	           "filename" : "pkg-2.0.tar.gz",
		   "chunkSize" : 262145,
		   "uploadDate" : ISODate("2013-12-16T16:52:01.089Z"),
		   "md5" : "ac98f052e55d0c556e6536796d6b54eb",
		   "length" : 2134}]`

	request, _ := http.NewRequest("GET", "tests/packages", nil)
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

	if !strings.EqualFold(response.Body.String(), results) {
		t.Fatalf("Body does not equal 'results'")
	}
}

func TestPackageUpload(t *testing.T) {
	request, _ := http.NewRequest("POST", "tests/packages", nil)
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
	request, _ := http.NewRequest("GET", "tests/packages/pkg.tgz", nil)
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
