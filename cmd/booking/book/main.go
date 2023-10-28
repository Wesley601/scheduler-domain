package main

import (
	"context"
	"encoding/json"
	"net/http"

	"alinea.com/internal/app"
	"alinea.com/internal/booking"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto := booking.CreateBookingDTO{}

	err := json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	result, err := app.BookingService.Book(context.Background(), dto)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	j, err := result.ToJSON()
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
