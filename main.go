package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/larwef/lunchlambda/getters"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/senders"
	"log"
	"os"
	"time"
)

// Environment variable names
const (
	MenuURL  = "MENU_URL"
	HookURL  = "HOOK_URL"
	SNSTopic = "SNS_TOPIC"
)

// Handler is the lambda handler function
func Handler() error {
	log.Println("lunchLambda invoked")

	menuURL := os.Getenv(MenuURL)
	hookURL := os.Getenv(HookURL)
	snsTopic := os.Getenv(SNSTopic)

	// Sources
	braathen := getters.NewBraathen(menuURL, time.Now())
	runner := menu.NewRunner(braathen)

	// Sinks
	slack := senders.NewSlack(hookURL)
	runner.AddSender(slack)

	config := aws.Config{Region: aws.String("eu-west-1")}

	if newSession, err := session.NewSession(&config); err == nil {
		newSns := senders.NewSns(snsTopic, "lunchLambda", sns.New(newSession))
		runner.AddSender(newSns)
	} else {
		log.Printf("Error configuring SNS: %v", err)
	}

	// Run
	err := runner.Run()

	defer log.Println("lunchLambda finished")
	return err
}

func main() {
	lambda.Start(Handler)
}
