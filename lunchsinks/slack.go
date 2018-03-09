package lunchsinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/larwef/lunchlambda/lunch"
	"net/http"
)

type (
	Slack struct {
		sinkUrl string
	}

	data struct {
		Text string `json:"text"`
	}
)

func NewSlack(url string) *Slack {
	return &Slack{sinkUrl: url}
}

func (s *Slack) SendMenu(menu lunch.Menu) error {
	d := data{Text: menu.ToString()}

	payload, err := json.Marshal(d)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.sinkUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("received response: \"%s\" on POST", resp.Status)
	}

	return nil
}
