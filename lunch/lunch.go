package lunch

import (
	"bytes"
	"fmt"
	"time"
)

type Menu struct {
	Timestamp  time.Time
	LunchItems []string
}

type Getter interface {
	GetMenu() (Menu, error)
}

type Poster interface {
	SendMenu(menus Menu) error
}

func (l *Menu) ToString() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Lunch Menu %02d.%02d.%02d\n", l.Timestamp.Day(), l.Timestamp.Month(), l.Timestamp.Year()))
	for _, item := range l.LunchItems {
		buffer.WriteString("- " + item + "\n")
	}
	buffer.WriteString("NB: Menu may vary from what's presented")

	return buffer.String()
}
