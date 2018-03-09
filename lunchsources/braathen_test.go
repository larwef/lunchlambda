package lunchsources

import (
	"fmt"
	"github.com/larwef/lunchlambda/lunch"
	"github.com/larwef/lunchlambda/testutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestBraathen_GetMenu(t *testing.T) {
	mux, url, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "../testdata/pageSource.html"))
	})

	loc, err := time.LoadLocation("Europe/Oslo")
	testutil.AssertNotError(t, err)
	menu, err := NewBraathen(url, time.Date(2018, time.March, 8, 0, 0, 0, 0, loc)).GetMenu()
	testutil.AssertNotError(t, err)

	expected := lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 8, 0, 0, 0, 0, loc),
		LunchItems: []string{"Fiskegrateng med pepperrotsmør", "Asiatisk marinert sopp", "Potetgrateng med spicy salat", "Gulrotsuppe med ingefær"},
	}

	testutil.AssertEqual(t, reflect.DeepEqual(menu, expected), true)
}

func TestBraathen_GetLunches(t *testing.T) {
	mux, url, teardown := testutil.Setup()

	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testutil.GetTestFileAsString(t, "../testdata/pageSource.html"))
	})

	menus, err := NewBraathen(url, time.Time{}).GetMenus()
	testutil.AssertNotError(t, err)
	testutil.AssertEqual(t, len(menus), 5)

	loc, err := time.LoadLocation("Europe/Oslo")
	testutil.AssertNotError(t, err)

	expected := make(map[string]lunch.Menu)
	expected["20180305"] = lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 5, 0, 0, 0, 0, loc),
		LunchItems: []string{"Fersk pasta med mornaysaus", "Potetsalat", "Fersk pasta med vegetar mornaysaus", "Grønnsakssuppe"},
	}
	expected["20180306"] = lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 6, 0, 0, 0, 0, loc),
		LunchItems: []string{"Fiskekaker med mandelpotet og skalldyrsaus", "Råkostsalat med urtevinaigrette", "Bakt brokkoli med bulgur", "Kyllingsuppe"},
	}
	expected["20180307"] = lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 7, 0, 0, 0, 0, loc),
		LunchItems: []string{"Lasagne al forno", "Tomat- og rødløksalat med balsamico", "Falafel med stekte grønnsaker og tahinidressing", "Kremet fiskesuppe"},
	}
	expected["20180308"] = lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 8, 0, 0, 0, 0, loc),
		LunchItems: []string{"Fiskegrateng med pepperrotsmør", "Asiatisk marinert sopp", "Potetgrateng med spicy salat", "Gulrotsuppe med ingefær"},
	}
	expected["20180309"] = lunch.Menu{
		Timestamp:  time.Date(2018, time.March, 9, 0, 0, 0, 0, loc),
		LunchItems: []string{"Røkt svinenakke med rødvinssaus og baconfrest sopp", "Nicoisesalat", "Vegetar Jambalaya", "Fisk Bisque"},
	}

	for key, value := range menus {
		testutil.AssertEqual(t, reflect.DeepEqual(value, expected[key]), true)
	}
}

func Test_getTimestampFromString(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Oslo")
	testutil.AssertNotError(t, err)

	input1 := "Tirsdag 6. mars 2018"
	output1 := time.Date(2018, time.March, 6, 0, 0, 0, 0, loc)

	input2 := "Gibberish 13. februar 2018"
	output2 := time.Date(2018, time.February, 13, 0, 0, 0, 0, loc)

	input3 := "Gibberish 24 dEsEmbEr 2018"
	output3 := time.Date(2018, time.December, 24, 0, 0, 0, 0, loc)

	inputs := []string{input1, input2, input3}
	outputs := []time.Time{output1, output2, output3}

	for i, element := range inputs {
		timestamp, err := getTimestampFromString(element)
		testutil.AssertNotError(t, err)
		testutil.AssertEqual(t, timestamp.Equal(outputs[i]), true)
	}
}

func Test_getMonthNumber(t *testing.T) {
	inputs := []string{"Januar", "februar", "MARS", "aPrIL", "mai", "juni", "Juli", "AuGUST", "SEPTEMBER", "Oktober", "novembeR", "Desember"}
	outputs := []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

	for i, input := range inputs {
		month, err := getMonthNumber(input)
		testutil.AssertNotError(t, err)
		testutil.AssertEqual(t, month, outputs[i])
	}
}

func Test_splitSlice(t *testing.T) {
	slice1 := []string{"1", "2", "3", "4", "5", "6", "Split", "7", "8", "9", "10", "11", "12"}
	slice2 := []string{"Split", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	slice3 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "Split"}
	slice4 := []string{"Split", "1", "2", "Split", "3", "4", "Split", "5", "6", "Split", "7", "8", "Split", "9", "10", "Split", "11", "12"}
	slice5 := []string{"1", "2", "3", "Split", "4", "5", "6", "Split", "7", "8", "9", "Split", "10", "11", "12", "Split"}
	slice6 := []string{"Split", "1", "2", "3", "Split", "4", "5", "Split", "Split", "6", "7", "8", "9", "10", "11", "Split", "12", "Split", "Split"}

	testSlices := [][]string{slice1, slice2, slice3, slice4, slice5, slice6}

	expectedLengths := []int{2, 1, 1, 6, 4, 4}

	for i, element := range testSlices {
		splitSlice := splitSlice(element, "Split")
		testutil.AssertEqual(t, len(splitSlice), expectedLengths[i])
	}
}
