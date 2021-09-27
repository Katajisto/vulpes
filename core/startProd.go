// +build prod

package core

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

// StartProd starts the API Gateway proxy in production mode.
func Startup() {
	// Create a new router
	log.Println("[ Lambdalith ] Lambda Mux Cold Start")
	r := mux.NewRouter()
	registerRoutes(r)
	log.Println("Routes registered")
	muxLambda = gorillamux.New(r)
	// Anything we do before this is sometimes kept in memory between lambda execs.
	// This means, that because we open up our DB before this point here, it will be kept open
	// between lambda execs. We could also consider doing some in memory caching here.
	lambda.Start(Handler)
}
