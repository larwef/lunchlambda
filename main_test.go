package main

import (
	"fmt"
	"github.com/larwef/lunchlambda/testutil"
	"net/http"
	"testing"
)

func TestHandler_NoMenus(t *testing.T) {
	mux, _, teardown := testutil.Setup()

	defer teardown()

	menuHandlerAssert := testutil.NewHandlerAssert(t, "/menu")
	mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "testdata/pageSource.html"))
		menuHandlerAssert.Called()
	})

	err := Handler()
	menuHandlerAssert.IsCalled()
	testutil.AssertEqual(t, err.Error(), "empty menu")
}

func TestHandler_Menu404(t *testing.T) {
	mux, _, teardown := testutil.Setup()

	defer teardown()

	menuHandlerAssert := testutil.NewHandlerAssert(t, "/menu")
	mux.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "testdata/404.html"))
		menuHandlerAssert.Called()
	})

	err := Handler()
	menuHandlerAssert.IsCalled()
	testutil.AssertEqual(t, err.Error(), "received response: \"404 Not Found\" on GET")
}
