package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jagijagijag1/pixela-go-client"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// InputBody is struct for body contents in api request
type InputBody struct {
	ActualTime string `json:"ActualTime"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	// extract env var
	user := os.Getenv("PIXELA_USER")
	token := os.Getenv("PIXELA_TOKEN")
	graph := os.Getenv("PIXELA_GRAPH")

	// parse request.Body and create InputBody
	var ib InputBody
	var err = json.Unmarshal([]byte(request.Body), &ib)
	if err != nil {
		fmt.Println("Error in parsing input body json")
		return Response{
			StatusCode: 400,
		}, err
	}
	fmt.Printf("actual date: %s\n", ib.ActualTime)

	// extract data from toggl
	date, quantity := getDelay(ib.ActualTime)
	if date == "-1" || quantity == "-1" {
		return Response{
			StatusCode: 400,
		}, nil
	}
	fmt.Printf("date: %s, quantity: %s\n", date, quantity)

	// record pixel
	perr := recordPixel(user, token, graph, date, quantity)
	if perr != nil {
		return Response{
			StatusCode: 400,
		}, nil
	}

	resp := Response{
		StatusCode: 200,
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func getDelay(actual string) (string, string) {
	y := time.Now()
	date := y.Format("20060102")

	t := time.Date(y.Year(), y.Month(), y.Day()-1, 22, 30, 0, 0, time.Local) // target
	a, err := time.Parse("2006-01-02T15:04:05-07:00", actual)                // actual
	if err != nil {
		fmt.Println(err)
		return "-1", "-1"
	}

	delay := t.Sub(a)
	quantity := strconv.FormatFloat(delay.Minutes(), 'f', 4, 64)
	fmt.Println(quantity)

	return date, quantity
}

func recordPixel(user, token, graph, date, quantity string) error {
	c := pixela.NewClient(user, token)

	// try to update
	err := c.UpdatePixelQuantity(graph, date, quantity)
	if err == nil {
		fmt.Println("updated")
	}

	return err
}
