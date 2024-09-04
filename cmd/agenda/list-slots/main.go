package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"wesley601.com/internal/app"
	"wesley601.com/internal/core"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	agendaID := request.PathParameters["agenda-id"]
	serviceID := request.PathParameters["service-id"]

	w, err := core.NewWindowFromString(request.QueryStringParameters["from"], request.QueryStringParameters["to"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	result, err := app.AgendaService.ListSlots(context.Background(), agendaID, serviceID, w)
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
