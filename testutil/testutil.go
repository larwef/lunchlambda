package testutil

import (
	"encoding/json"
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
	os.Setenv("MENU_URL", server.URL+"/menu")

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

func AssertJSONSEqual(t *testing.T, json1 string, json2 string) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(json1), &o1)
	if err != nil {
		t.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(json2), &o2)
	if err != nil {
		t.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Json strings are not equal:")
		t.Errorf("Actual: %s", json1)
		t.Errorf("Expected: %s", json2)
	}
}
