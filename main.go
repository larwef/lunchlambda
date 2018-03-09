package lunchlambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/larwef/lunchlambda/lunchsources"
	"log"
	"os"
	"time"
)

const (
	LunchUrl = "LUNCH_URL"
	HookUrl  = "HOOK_URL"
)

func Handler() error {
	log.Println("lunchLambda invoked")

	lunchUrl := os.Getenv(LunchUrl)

	menus, err := lunchsources.NewBraathen(lunchUrl, time.Now()).GetMenus()
	if err != nil {
		log.Printf("received error from lunchgetter:%s", err)
		return err
	}

	for _, menu := range menus {
		log.Println("\n" + menu.ToString())
	}

	log.Println("lunchLambda finished")
	return nil
}

func main() {
	lambda.Start(Handler)
}
