package lunch

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Menu struct {
	Timestamp  time.Time `json:"timestamp"`
	LunchItems []string  `json:"lunch_items"`
}

type Getter interface {
	GetLunches(url string) ([]Menu, error)
}

func (l *Menu) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Lunch Menu %02d.%02d.%02d\n", l.Timestamp.Day(), l.Timestamp.Month(), l.Timestamp.Year()))
	for _, item := range l.LunchItems {
		buffer.WriteString("- " + item + "\n")
	}

	return strings.TrimSuffix(buffer.String(), "\n")
}
