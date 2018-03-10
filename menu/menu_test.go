package menu

import (
	"github.com/larwef/lunchlambda/testutil"
	"testing"
	"time"
)

func TestMenu_ToString(t *testing.T) {
	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	menu := Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
	}

	expected := "Menu 07.03.2018\n- Some vegetarian alternative\n- Some main dish\n- Some soup\nNB: Menu may vary from what's presented"
	testutil.AssertEqual(t, menu.ToString(), expected)
}

func TestMenu_ToString_EmptyMenu(t *testing.T) {
	menu := Menu{}
	testutil.AssertEqual(t, menu.ToString(), "")
}
