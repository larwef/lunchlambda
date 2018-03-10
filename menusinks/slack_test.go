package menusinks

import (
	"fmt"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/testutil"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestSlack_SendMenu(t *testing.T) {
	mux, url, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		testutil.AssertJSONSEqual(t, string(bytes), testutil.GetTestFileAsString(t, "../testdata/slackRequest.json"))
		fmt.Fprint(w, "ok")
	})

	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	m := menu.Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
	}

	err := NewSlack(url).SendMenu(m)
	testutil.AssertNotError(t, err)
}

func TestSlack_SendMenu_404(t *testing.T) {
	mux, url, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	m := menu.Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
	}

	err := NewSlack(url).SendMenu(m)
	testutil.AssertEqual(t, err.Error(), "received response: \"404 Not Found\" on POST")
}

func TestSlack_SendMenu_EmptyMenu(t *testing.T) {
	m := menu.Menu{}
	err := NewSlack("").SendMenu(m)
	testutil.AssertEqual(t, err.Error(), EmptyMenuError.Error())
}
