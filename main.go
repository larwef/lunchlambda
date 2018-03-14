package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/larwef/lunchlambda/getters"
	"github.com/larwef/lunchlambda/menu"
	"github.com/larwef/lunchlambda/senders"
	"log"
	"os"
	"time"
)

const (
	MenuURL = "MENU_URL"
	HookURL = "HOOK_URL"
)

func Handler() error {
	log.Println("lunchLambda invoked")

	menuURL := os.Getenv(MenuURL)
	hookURL := os.Getenv(HookURL)

	// Sources
	braathen := getters.NewBraathen(menuURL, time.Now())

	// Sinks
	slack := senders.NewSlack(hookURL)

	// Run
	err := menu.NewRunner(braathen).AddSender(slack).Run()

	defer log.Println("lunchLambda finished")
	return err
}

func main() {
	lambda.Start(Handler)
}
