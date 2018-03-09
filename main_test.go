package lunchlambda

import (
	"fmt"
	"github.com/larwef/lunchlambda/testutil"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	mux, _, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/lunch", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "testdata/pageSource.html"))
	})

	mux.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertEqual(t, r.Method, http.MethodPost)
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "testdata/pageSource.html"))
	})

	err := Handler()
	testutil.AssertNotError(t, err)

}

func TestHandler_Lunch404(t *testing.T) {
	mux, _, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/lunch", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "testdata/404.html"))
	})

	err := Handler()
	testutil.AssertEqual(t, err.Error(), "received response: \"404 Not Found\" on GET")
}
