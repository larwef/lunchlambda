package lunchlambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/larwef/lunchlambda/lunchgetters"
	"log"
	"os"
)

const (
	LunchUrl = "LUNCH_URL"
	HookUrl  = "HOOK_URL"
)

func Handler() error {
	log.Println("lunchLambda invoked")

	lunchUrl := os.Getenv(LunchUrl)

	menus, err := lunchgetters.NewBraathenLunchGetter().GetLunches(lunchUrl)
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
