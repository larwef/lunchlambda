package lunchlambda

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func setup() (mux *http.ServeMux, url string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	server := httptest.NewServer(apiHandler)

	return mux, url, server.Close
}

func TestHandler(t *testing.T) {
	Handler()
}
