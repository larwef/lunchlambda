package senders

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/larwef/lunchlambda/menu"
	"log"
)

// SNSSender is used to send a menu to AWS SNS
type SNSSender struct {
	topicArn  string
	subject   string
	snsClient snsiface.SNSAPI
}

// NewSns is a constructor for the SNSSender object
func NewSns(topicArn string, subject string, snsClient snsiface.SNSAPI) *SNSSender {
	return &SNSSender{topicArn: topicArn, subject: subject, snsClient: snsClient}
}

// SendMenu sends a menu to SNS
func (s *SNSSender) SendMenu(m menu.Menu) error {
	log.Printf("Publishing menu to SNS topic: %s", s.topicArn)
	message := m.ToString()

	publishInput := sns.PublishInput{
		Message:  &message,
		Subject:  &s.subject,
		TopicArn: &s.topicArn,
	}

	_, err := s.snsClient.Publish(&publishInput)
	return err
}
