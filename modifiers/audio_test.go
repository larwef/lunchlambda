package modifiers

import (
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/testutil"
	"testing"
	"time"
)

type mockPollyClient struct {
	pollyiface.PollyAPI
	handler func(input *polly.SynthesizeSpeechInput)
}

type mockS3UploadManager struct {
	s3manageriface.UploaderAPI
	handler func(input *s3manager.UploadInput)
}

func (m *mockPollyClient) SynthesizeSpeech(*polly.SynthesizeSpeechInput) (*polly.SynthesizeSpeechOutput, error) {
	return &polly.SynthesizeSpeechOutput{}, nil
}

func (m *mockS3UploadManager) Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return &s3manager.UploadOutput{Location: "resourceLocation"}, nil
}

func TestAudio_Modify(t *testing.T) {
	mockPolly := &mockPollyClient{}
	mockS3 := &mockS3UploadManager{}

	mockPolly.handler = func(input *polly.SynthesizeSpeechInput) {
		testutil.AssertEqual(t, input.OutputFormat, "mp3")
		testutil.AssertEqual(t, input.Text, "Some vegetarian alternative. Some main dish. Some soup.")
		testutil.AssertEqual(t, input.TextType, "text")
		testutil.AssertEqual(t, input.VoiceId, "someVoiceId")
	}

	mockS3.handler = func(input *s3manager.UploadInput) {
		testutil.AssertEqual(t, input.ACL, "public-read")
		testutil.AssertEqual(t, input.Bucket, "someBucket")
		testutil.AssertEqual(t, input.ContentType, "audio/mp3")
		testutil.AssertEqual(t, input.Key, "20180307.mp3")
	}

	audioModifier := NewAudio("someVoiceId", "someBucket", mockPolly, mockS3)

	menuItems := []string{"Some vegetarian alternative", "Some main dish", "Some soup"}
	timeStamp, _ := time.Parse(time.RFC3339, "2018-03-07T16:30:03Z")

	m := menu.Menu{
		Timestamp: timeStamp,
		MenuItems: menuItems,
		Source:    "someSource",
	}

	_ = audioModifier.Modify(&m)

	testutil.AssertEqual(t, m.AudioURL, "resourceLocation")
}
