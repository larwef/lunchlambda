package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/larwef/lunchlambda/menu"
	"log"
	"net/http"
)

const (
	// BulletPoint is printet in front of all menu items
	BulletPoint = ":knife_fork_plate:"
)

type (
	// Slack is used to post a menu to Slack. Implements the Sender interface
	Slack struct {
		sinkURL string
	}

	message struct {
		Attachments []attachment `json:"attachments"`
	}

	attachment struct {
		Title  string `json:"title"`
		Text   string `json:"text"`
		Footer string `json:"footer"`
	}
)

// NewSlack is a constructor for the Slack object
func NewSlack(url string) *Slack {
	return &Slack{sinkURL: url}
}

// SendMenu posts a menu to Slack
func (s *Slack) SendMenu(m menu.Menu) error {
	if len(m.MenuItems) < 1 {
		return menu.ErrEmptyMenu
	}

	log.Printf("Sending menu to: %s\n", s.sinkURL)

	var buffer bytes.Buffer
	for _, item := range m.MenuItems {
		buffer.WriteString(BulletPoint + " " + item + "\n")
	}

	// Make prettier
	var audioString string
	if m.AudioURL != "" {
		audioString = "\nAs mp3: " + m.AudioURL
	}

	a := attachment{
		Title:  fmt.Sprintf("Menu %s %02d.%02d.%02d", m.Timestamp.Weekday(), m.Timestamp.Day(), m.Timestamp.Month(), m.Timestamp.Year()),
		Text:   string(buffer.Bytes()[:buffer.Len()-1]),
		Footer: "Source: " + m.Source + audioString + "\nNB: Menu may vary from what's presented",
	}
	mes := message{
		Attachments: []attachment{a},
	}

	payload, err := json.Marshal(mes)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.sinkURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("received response: \"%s\" on POST", resp.Status)
	}

	return nil
}
