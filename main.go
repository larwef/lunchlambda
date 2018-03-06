package lunchlambda

import(
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func Handler() (error) {
	log.Println("lunchLambda invoked")
	return nil
}

func main() {
	lambda.Start(Handler)
}