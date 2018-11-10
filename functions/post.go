package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the DAO
	fmt.Println("Received body: ", request.Body)
	item, err := Post(request.Body)
	if err != nil {
		fmt.Println("Got error calling post")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	// Log and return result
	fmt.Println("Wrote item:  ", item)
	return events.APIGatewayProxyResponse{Body: "Success\n", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
