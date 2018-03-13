package menu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var ErrEmptyMenu = errors.New("empty menu")

type (
	Slack struct {
		sinkURL string
	}

	data struct {
		Text string `json:"text"`
	}
)

func NewSlack(url string) *Slack {
	return &Slack{sinkURL: url}
}

func (s *Slack) SendMenu(menu Menu) error {
	if len(menu.MenuItems) < 1 {
		return ErrEmptyMenu
	}

	log.Printf("Sending menu to: %s\n", s.sinkURL)

	d := data{Text: menu.ToString()}

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
