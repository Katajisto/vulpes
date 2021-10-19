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
	r := mux.NewRouter()
	registerRoutes(r)

	// Serve static assets from the static directory. In production serve these from S3 through API Gateway.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("s3/static"))))

	// Bind the router to a port
	http.ListenAndServe(":1337", r)
}
