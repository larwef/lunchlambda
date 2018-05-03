package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/larwef/lunchlambda/getters"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/modifiers"
	"github.com/larwef/lunchlambda/senders"
	"log"
	"os"
	"strings"
	"time"
)

// Environment variable names
const (
	MenuURL      = "MENU_URL"
	HookURL      = "HOOK_URL"
	SNSTopic     = "SNS_TOPIC"
	AudioEnabled = "AUDIO_ENABLED"
	AudioVoice   = "AUDIO_VOICE_ID"
	AudioBucket  = "AUDIO_BUCKET"
	AWSRegion    = "CURRENT_AWS_REGION"
)

// Handler is the lambda handler function
func Handler() error {
	log.Println("lunchLambda invoked")

	menuURL := os.Getenv(MenuURL)
	hookURL := os.Getenv(HookURL)
	snsTopic := os.Getenv(SNSTopic)
	audioEnabled := os.Getenv(AudioEnabled)
	audioVoice := os.Getenv(AudioVoice)
	audioBucket := os.Getenv(AudioBucket)
	awsRegion := os.Getenv(AWSRegion)

	braathen := getters.NewBraathen(menuURL, time.Now())
	runner := menu.NewRunner(braathen)

	config := aws.Config{Region: aws.String(awsRegion)}
	if newSession, err := session.NewSession(&config); err == nil {
		if strings.ToLower(audioEnabled) == "true" {
			modifier := modifiers.NewAudio(audioVoice, audioBucket, polly.New(newSession), s3manager.NewUploader(newSession))
			runner.AddModifier(modifier)
		}
		newSns := senders.NewSns(snsTopic, "lunchLambda", sns.New(newSession))
		runner.AddSender(newSns)
	} else {
		log.Printf("error getting AWS session: %v. Dependent objects are not added", err)
	}

	runner.AddSender(senders.NewSlack(hookURL))

	// Run
	err := runner.Run()

	defer log.Println("lunchLambda finished")
	return err
}

func main() {
	lambda.Start(Handler)
}
