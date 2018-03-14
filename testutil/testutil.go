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

// HandlerAssert is a helper object used to assert that a handler is called during testing
type HandlerAssert struct {
	t        *testing.T
	isCalled bool
	handler  string
}

// NewHandlerAssert is a constructor for the HandlerAssert
func NewHandlerAssert(t *testing.T, handler string) *HandlerAssert {
	return &HandlerAssert{t: t, handler: handler}
}

// IsCalled checks if the Called function was called on the HandlerAssert object
func (h *HandlerAssert) IsCalled() {
	if !h.isCalled {
		h.t.Fatalf("Handler %s not called", h.handler)
	}
}

// Called is used in the handler to set the called variable checked in the IsCalled function
func (h *HandlerAssert) Called() {
	h.isCalled = true
}

// Setup returns a mux to register handlers on, url to the mock server and a teardown function
func Setup() (mux *http.ServeMux, url string, teardown func()) {
	mux = http.NewServeMux()
	handler := http.NewServeMux()
	handler.Handle("/", mux)
	server := httptest.NewServer(handler)

	os.Setenv("HOOK_URL", server.URL+"/hook")
	os.Setenv("MENU_URL", server.URL+"/menu")

	return mux, server.URL, server.Close
}

// GetTestFileAsString gets the file as a string... Test will fail if the file cannot be read.
func GetTestFileAsString(t *testing.T, filepath string) string {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// AssertNotError asserts if an error equals nil or fails the test
func AssertNotError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got unexpected error: %s", err)
	}
}

// AssertEqual asserts if two object are same type and equal value
func AssertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v %v to be equal to %v %v", reflect.TypeOf(actual).Name(), actual, reflect.TypeOf(expected).Name(), expected)
	}
}

// AssertJSONSEqual asserts two json strings for equality by unmarshalling them to account for formatting.
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
