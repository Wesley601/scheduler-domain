package main

import (
	"context"
	"encoding/json"
	"net/http"

	"alinea.com/internal/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, err := app.AgendaService.List(context.Background())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(j),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
