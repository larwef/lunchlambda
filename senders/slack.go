package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/larwef/lunchlambda/menu"
	"log"
	"net/http"
)

type (
	// Slack is used to post a menu to Slack. Implements the Sender interface
	Slack struct {
		sinkURL string
	}

	data struct {
		Text string `json:"text"`
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

	d := data{Text: m.ToString()}

	payload, err := json.Marshal(d)
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
