package senders

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/testutil"
	"testing"
	"time"
)

type mockSNSClient struct {
	snsiface.SNSAPI
	handler func(input *sns.PublishInput)
}

func (m *mockSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	m.handler(input)
	messageID := "testMessageId"
	return &sns.PublishOutput{MessageId: &messageID}, nil
}

func TestSNSSender_SendMenu(t *testing.T) {
	mockSvc := &mockSNSClient{}

	expectedMessage := "Menu 07.03.2018\n- Some vegetarian alternative\n- Some main dish\n- Some soup\nSource: someSource\nNB: Menu may vary from what's presented"
	mockSvc.handler = func(input *sns.PublishInput) {
		testutil.AssertEqual(t, *input.TopicArn, "testTopic")
		testutil.AssertEqual(t, *input.Subject, "testSubject")
		testutil.AssertEqual(t, *input.Message, expectedMessage)
	}

	newSns := NewSns("testTopic", "testSubject", mockSvc)

	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	m := menu.Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
		Source:    "someSource",
	}

	_ = newSns.SendMenu(m)
}
