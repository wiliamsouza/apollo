package api

import (
	"net/http"
)

func packageList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	w.Write([]byte(results))
}

func packageUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func packageDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
