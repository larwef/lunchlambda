package lunch

import (
	"github.com/larwef/lunchlambda/testutil"
	"testing"
	"time"
)

func TestLunch_ToString(t *testing.T) {
	lunchItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	lunches := Menu{
		Timestamp:  timeStamp,
		LunchItems: lunchItems,
	}

	expected := "Lunch Menu 07.03.2018\n- Some vegetarian alternative\n- Some main dish\n- Some soup"
	testutil.AssertEqual(t, lunches.ToString(), expected)
}
