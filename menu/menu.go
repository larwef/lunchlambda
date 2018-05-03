package menu

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// ErrEmptyMenu is returned when Menu doesn't contain any MenuItems
var ErrEmptyMenu = errors.New("empty menu")

// Menu defines a menu for a given time from a given source.
type Menu struct {
	Timestamp time.Time
	MenuItems []string
	Source    string
	AudioURL  string
}

// Runner object holds a getter to get a menu from some source and an array of senders which will send the menu to their
// respective endpoints
type Runner struct {
	getter    Getter
	modifiers []Modifier
	senders   []Sender
}

// Getter defines the behaviour of getting a menu from a source
type Getter interface {
	GetMenu() (Menu, error)
}

// Modifier defines an interface for modyfying menu object before sending it
type Modifier interface {
	Modify(menu *Menu) error
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

// AddModifier adds a modifier to the Runner obbject
func (r *Runner) AddModifier(modifier Modifier) *Runner {
	r.modifiers = append(r.modifiers, modifier)
	return r
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

	for _, modifier := range r.modifiers {
		modifier.Modify(&menu)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(r.senders))
	for _, element := range r.senders {
		go func(sender Sender) {
			if err := sender.SendMenu(menu); err != nil {
				log.Printf("Encountered error when sending menu: %v", err)
			}
			waitGroup.Done()
		}(element)
	}

	waitGroup.Wait()
	return nil
}
