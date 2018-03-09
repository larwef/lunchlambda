package testutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func Setup() (mux *http.ServeMux, url string, teardown func()) {
	mux = http.NewServeMux()
	handler := http.NewServeMux()
	handler.Handle("/", mux)
	server := httptest.NewServer(handler)

	os.Setenv("HOOK_URL", server.URL+"/hook")
	os.Setenv("LUNCH_URL", server.URL+"/lunch")

	return mux, server.URL, server.Close
}

func GetTestFileAsString(t *testing.T, filepath string) string {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func AssertNotError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got unexpected error: %s", err)
	}
}

func AssertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v %v to be equal to %v %v", reflect.TypeOf(actual).Name(), actual, reflect.TypeOf(expected).Name(), expected)
	}
}
