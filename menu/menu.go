package menu

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"
)

// ErrEmptyMenu is returned when Menu doesn't contain any MenuItems
var ErrEmptyMenu = errors.New("empty menu")

// Menu defines a menu for a given time from a given source.
type Menu struct {
	Timestamp time.Time
	MenuItems []string
	Source    string
}

// Runner object holds a getter to get a menu from some source and an array of senders which will send the menu to their
// respective endpoints
type Runner struct {
	getter  Getter
	senders []Sender
}

// Getter defines the behaviour of getting a menu from a source
type Getter interface {
	GetMenu() (Menu, error)
}

// Sender defines the behaviour of sending a menu to an endpoint
type Sender interface {
	SendMenu(menus Menu) error
}

// IsEmpty returns true if the Menu object contains no MenuItems
func (m *Menu) IsEmpty() bool {
	return len(m.MenuItems) < 1
}

// ToString returns a string representation of a Menu object
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

// NewRunner is a constructor for the Runner object
func NewRunner(getter Getter) *Runner {
	return &Runner{getter: getter}
}

// AddSender adds an object implementing the Sender interface to the list of senders
func (r *Runner) AddSender(sender Sender) *Runner {
	r.senders = append(r.senders, sender)
	return r
}

// Run executes the Getter and feeds the result to all the Senders
func (r *Runner) Run() error {
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
