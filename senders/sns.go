package senders

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/larwef/lunchlambda/menu"
)

// SNSSender is used to send a menu to AWS SNS
type SNSSender struct {
	topicArn  string
	snsClient snsiface.SNSAPI
}

// NewSns is a constructor for the SNSSender object
func NewSns(topicArn string, snsClient snsiface.SNSAPI) *SNSSender {
	return &SNSSender{topicArn: topicArn, snsClient: snsClient}
}

// SendMenu sends a menu to SNS
func (s *SNSSender) SendMenu(m menu.Menu) error {
	message := m.ToString()

	publishInput := sns.PublishInput{
		Message:  &message,
		TopicArn: &s.topicArn,
	}

	_, err := s.snsClient.Publish(&publishInput)
	if err != nil {
		return err
	}

	return nil
}