package menu

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"
)

var ErrEmptyMenu = errors.New("empty menu")

type Menu struct {
	Timestamp time.Time
	MenuItems []string
	Source    string
}

type runner struct {
	getter  Getter
	senders []Sender
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
func NewRunner(getter Getter) *runner {
	return &runner{getter: getter}
}

func (r *runner) AddSender(sender Sender) *runner {
	r.senders = append(r.senders, sender)
	return r
}

func (r *runner) Run() error {
	if len(r.senders) < 1 {
		log.Println("No senders registerd for runner")
	}

	menu, err := r.getter.GetMenu()

	if err != nil {
		log.Printf("received error from menusource: %s", err)
		return err
	}

	if menu.IsEmpty() {
		return ErrEmptyMenu
	}

	for _, element := range r.senders {
		if err := element.SendMenu(menu); err != nil {
			log.Printf("Encountered error when sending menu: %v", err)
		}
	}

	return nil
}
