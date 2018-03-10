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
	MenuUrl = "MENU_URL"
	HookUrl = "HOOK_URL"
)

func Handler() error {
	log.Println("lunchLambda invoked")

	menuUrl := os.Getenv(MenuUrl)

	menu, err := menusources.NewBraathen(menuUrl, time.Now()).GetMenu()
	if err != nil {
		log.Printf("received error from menusource: %s", err)
		return err
	}

	if !menu.IsEmpty() {
		if err := menusinks.NewSlack(HookUrl).SendMenu(menu); err != nil {
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
