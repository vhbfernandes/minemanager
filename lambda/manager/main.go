package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, r events.APIGatewayV2HTTPRequest) (Response, error) {
	var buf bytes.Buffer
	//todo replace with something better
	if !auth(r.Headers["authorization"]) {
		fmt.Printf("unauthorized request")
		resp := Response{
			StatusCode: 401,
		}
		return resp, nil
	}

	var msg string
	switch r.RawPath {
	case "/start":
		startInstance()
		msg = "instance starting"
		break
	case "/stop":
		stopInstance()
		msg = "instance stopping"
		break
	case "/ip":
		msg = getInstanceIP()
		break
	}
	fmt.Printf("%s", msg)

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
			"X-MyCompany-Func-Reply": "manager-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func auth(token string) bool {
	return token == os.Getenv("TOKEN")
}
