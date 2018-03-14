package menu_test

import (
	"fmt"
	"github.com/larwef/lunchlambda/getters"
	. "github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/senders"
	"github.com/larwef/lunchlambda/testutil"
	"net/http"
	"testing"
	"time"
)

func TestMenu_ToString(t *testing.T) {
	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	menu := Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
		Source:    "someSource",
	}

	expected := "Menu 07.03.2018\n- Some vegetarian alternative\n- Some main dish\n- Some soup\nSource: someSource\nNB: Menu may vary from what's presented"
	testutil.AssertEqual(t, menu.ToString(), expected)
}

func TestMenu_ToString_EmptyMenu(t *testing.T) {
	menu := Menu{}
	testutil.AssertEqual(t, menu.ToString(), "")
}

func TestRunner(t *testing.T) {
	mux, url, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "../testdata/pageSource.html"))
	})

	slack1HandlerAssert := testutil.NewHandlerAssert(t, "slack1")
	mux.HandleFunc("/slack1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
		slack1HandlerAssert.Called()
	})

	slack2HandlerAssert := testutil.NewHandlerAssert(t, "slack2")
	mux.HandleFunc("/slack2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
		slack2HandlerAssert.Called()
	})

	slack3HandlerAssert := testutil.NewHandlerAssert(t, "slack3")
	mux.HandleFunc("/slack3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
		slack3HandlerAssert.Called()
	})

	loc, err := time.LoadLocation("Europe/Oslo")
	testutil.AssertNotError(t, err)
	braathens := getters.NewBraathen(url, time.Date(2018, time.March, 8, 0, 0, 0, 0, loc))
	testutil.AssertNotError(t, err)

	slack1 := senders.NewSlack(url + "/slack1")
	slack2 := senders.NewSlack(url + "/slack2")
	slack3 := senders.NewSlack(url + "/slack3")

	runner := NewRunner(braathens).AddSender(slack1).AddSender(slack2).AddSender(slack3)

	runner.Run()

	slack1HandlerAssert.IsCalled()
	slack2HandlerAssert.IsCalled()
	slack2HandlerAssert.IsCalled()
}
