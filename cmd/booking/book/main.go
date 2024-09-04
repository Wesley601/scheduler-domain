package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"wesley601.com/internal/app"
	"wesley601.com/internal/booking"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dto := booking.CreateBookingDTO{}

	err := json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	err = app.BookingService.Book(context.Background(), dto)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string("ok"),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
