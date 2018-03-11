package lunchlambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/larwef/lunchlambda/menusinks"
	"github.com/larwef/lunchlambda/menusources"
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

	menu, err := menusources.NewBraathen(menuURL, time.Now()).GetMenu()
	if err != nil {
		log.Printf("received error from menusource: %s", err)
		return err
	}

	if !menu.IsEmpty() {
		if err := menusinks.NewSlack(hookURL).SendMenu(menu); err != nil {
			log.Printf("received error from menusink: %s", err)
			return err
		}
	}

	log.Println("lunchLambda finished")
	return nil
}

func main() {
	lambda.Start(Handler)
}
