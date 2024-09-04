package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"wesley601.com/internal/agenda"
	"wesley601.com/internal/app"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto := agenda.CreateAgendaDTO{}

	err := json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	result, err := app.AgendaService.Create(context.Background(), dto)
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
