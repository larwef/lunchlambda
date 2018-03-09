package lunchgetters

import (
	"fmt"
	"github.com/larwef/lunchlambda/lunch"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	containerClass = "article-related-images slideshow"
)

type BraathenLunchGetter struct{}

func NewBraathenLunchGetter() *BraathenLunchGetter {
	return &BraathenLunchGetter{}
}

func (b *BraathenLunchGetter) GetLunches(url string) ([]lunch.Menu, error) {
	log.Printf("Getting lunch menu from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return []lunch.Menu{}, fmt.Errorf("received response: \"%s\" on GET", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error parsing")
	}

	node := getNodeByClass(containerClass, doc)
	if node == nil {
		log.Fatalf("didn't find container node")
	}
	text := getContentTextFromNode(node)
	splitSlice := splitSlice(text, "DAGENS MENY")

	var menus []lunch.Menu
	for _, slice := range splitSlice {
		menu := lunch.Menu{}
		menu.Timestamp, err = getTimestampFromString(slice[0])
		if err != nil {
			return []lunch.Menu{}, err
		}
		for _, line := range slice[1:] {
			for _, s := range strings.Split(line, "|") {
				str := strings.Trim(s, " ")
				if str != "" && str != " " {
					menu.LunchItems = append(menu.LunchItems, str)
				}
			}
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func getNodeByClass(class string, node *html.Node) *html.Node {
	var result *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if isClass(class, n) {
			result = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	return result
}

func isClass(class string, node *html.Node) bool {
	for _, attribute := range node.Attr {
		if attribute.Key == "class" && attribute.Val == class {
			return true
		}
	}
	return false
}

func getContentTextFromNode(node *html.Node) []string {
	var textLines []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			s := n.Data
			s = stringMinifier(s)
			s = strings.Trim(s, " \r\n\t")
			if s != "" && s != " " {
				textLines = append(textLines, s)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	return textLines
}

// Removes unneccesary space in sting
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

func splitSlice(startSlice []string, splitElement string) [][]string {
	var splitSlice [][]string
	lastSplitIndex := 0
	for i, element := range startSlice {
		if element == splitElement {
			if i != 0 {
				if slice := startSlice[lastSplitIndex+1 : i]; len(slice) > 0 {
					splitSlice = append(splitSlice, slice)
				}
			}
			lastSplitIndex = i
		}
	}
	if lastSplitIndex != len(startSlice)-1 {
		splitSlice = append(splitSlice, startSlice[lastSplitIndex+1:])
	}
	return splitSlice
}

func getTimestampFromString(str string) (time.Time, error) {
	var year, day int
	var month time.Month
	var loc *time.Location
	var err error

	s := strings.Split(str, " ")
	if day, err = strconv.Atoi(strings.Trim(s[1], " .")); err != nil {
		return time.Time{}, err
	}
	if month, err = getMonthNumber(s[2]); err != nil {
		return time.Time{}, err
	}
	if year, err = strconv.Atoi(strings.Trim(s[3], " .")); err != nil {
		return time.Time{}, err
	}
	if loc, err = time.LoadLocation("Europe/Oslo"); err != nil {
		return time.Time{}, err
	}
	return time.Date(year, month, day, 0, 0, 0, 0, loc), nil
}

func getMonthNumber(month string) (time.Month, error) {
	var timeMonth time.Month
	switch strings.ToUpper(month) {
	case "JANUAR":
		timeMonth = time.January
		break
	case "FEBRUAR":
		timeMonth = time.February
		break
	case "MARS":
		timeMonth = time.March
		break
	case "APRIL":
		timeMonth = time.April
		break
	case "MAI":
		timeMonth = time.May
		break
	case "JUNI":
		timeMonth = time.June
		break
	case "JULI":
		timeMonth = time.July
		break
	case "AUGUST":
		timeMonth = time.August
		break
	case "SEPTEMBER":
		timeMonth = time.September
		break
	case "OKTOBER":
		timeMonth = time.October
		break
	case "NOVEMBER":
		timeMonth = time.November
		break
	case "DESEMBER":
		timeMonth = time.December
		break
	default:
		return 0, fmt.Errorf("couldn't get month from string: %s", month)
	}
	return timeMonth, nil
}
