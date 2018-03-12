package menu

import (
	"bytes"
	"fmt"
	"time"
)

type Menu struct {
	Timestamp time.Time
	MenuItems []string
	Source    string
}

type Runner struct {
	sinks []Getter
}

type Getter interface {
	GetMenu() (Menu, error)
}

type Sender interface {
	SendMenu(menus Menu) error
}

func (m *Menu) IsEmpty() bool {
	return len(m.MenuItems) < 1
}

func (m *Menu) ToString() string {
	if len(m.MenuItems) < 1 {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Menu %02d.%02d.%02d\n", m.Timestamp.Day(), m.Timestamp.Month(), m.Timestamp.Year()))
	for _, item := range m.MenuItems {
		buffer.WriteString("- " + item + "\n")
	}
	buffer.WriteString(fmt.Sprintf("Source: %s\n", m.Source))
	buffer.WriteString("NB: Menu may vary from what's presented")

	return buffer.String()
}
