package modifiers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/larwef/lunchlambda/menu"
	"strings"
)

// Audio modifier used AWS Polly to synthesize speech from MenuItems in a Menu object, adds the resulting .mp3 file to S3
// and adds a resource link in the Menu object
type Audio struct {
	voiceID         string
	bucket          string
	pollyClient     pollyiface.PollyAPI
	s3UploadManager s3manageriface.UploaderAPI
}

// NewAudio is a constructor for Audio modifier
func NewAudio(voiceID string, bucket string, pollyClient pollyiface.PollyAPI, s3UploadManager s3manageriface.UploaderAPI) *Audio {
	return &Audio{
		voiceID:         voiceID,
		bucket:          bucket,
		pollyClient:     pollyClient,
		s3UploadManager: s3UploadManager,
	}
}

// Modify modifies the menu object
func (a *Audio) Modify(m *menu.Menu) error {
	synthesizeSpeechInput := polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(polly.OutputFormatMp3),
		Text:         aws.String(strings.Join(m.MenuItems, ". ") + "."),
		TextType:     aws.String(polly.TextTypeText),
		VoiceId:      &a.voiceID,
	}

	synthesizeSpeechOutput, err := a.pollyClient.SynthesizeSpeech(&synthesizeSpeechInput)
	if err != nil {
		return fmt.Errorf("error synthesizing speech: %v", err)
	}

	s3UploadInput := s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(a.bucket),
		ContentType: aws.String("audio/mp3"),
		Key:         aws.String(fmt.Sprintf("%04d%02d%02d.mp3", m.Timestamp.Year(), m.Timestamp.Month(), m.Timestamp.Day())),
		Body:        synthesizeSpeechOutput.AudioStream,
	}

	uploadOutput, err := a.s3UploadManager.Upload(&s3UploadInput)
	if err != nil {
		return fmt.Errorf("error uloading to s3: %v", err)
	}

	m.AudioURL = uploadOutput.Location

	return nil
}
