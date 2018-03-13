package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/larwef/lunchlambda/menu"
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
	braathen := menu.NewBraathen(menuURL, time.Now())

	// Sinks
	slack := menu.NewSlack(hookURL)

	// Run
	err := menu.NewRunner(braathen).AddSender(slack).Run()

	defer log.Println("lunchLambda finished")
	return err
}

func main() {
	lambda.Start(Handler)
}
