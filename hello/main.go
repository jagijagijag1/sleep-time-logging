package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// InputBody is struct for body contents in api request
type InputBody struct {
	Test string `json:"test"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	// parse request.Body and create InputBody
	var ib InputBody
	var err = json.Unmarshal([]byte(request.Body), &ib)
	if err != nil {
		var em = "Error in parsing input body json"
		fmt.Println(em)
		resp := Response{
			StatusCode: 400,
			Body:       em,
		}
		return resp, err
	}

	// if parsed, create return message
	var msg = "input is  " + ib.Test
	body, err := json.Marshal(map[string]interface{}{
		"message": msg,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
