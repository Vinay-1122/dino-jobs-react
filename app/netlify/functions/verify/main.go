package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	token := request.QueryStringParameters["token"]
	fmt.Println(token)
	isAuth, email, _, err := verifyToken(token)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
			},
			Body: `{"error": "Error verifying token"}`,
		}, nil
	}

	client := connect()
	if isAuth {
		updateEmailVerified(client, email)
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
			},
			Body: fmt.Sprintf(`{"res": "Email verified for %s"}`, email),
		}, nil
	}
	defer disconnect(client)
	return &events.APIGatewayProxyResponse{
		StatusCode: 401,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		},
		Body: `{"error": "Unauthorized Access"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
