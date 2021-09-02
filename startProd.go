// +build prod

package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func startup() {
	// Create a new router
	log.Println("[ Lambdalith ] Lambda Mux Cold Start")
	r := mux.NewRouter()
	registerRoutes(r)
	log.Println("Routes registered")
	muxLambda = gorillamux.New(r)
	lambda.Start(Handler)
}
