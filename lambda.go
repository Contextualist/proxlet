package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/zeit/now-builders/utils/go/bridge"
	"net/http"
)

func awshandler(q events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	encoding := ""
	if q.IsBase64Encoded {
		encoding = "base64"
	}
	r, err := bridge.Serve(http.HandlerFunc(Handler), &bridge.Request{
		Host:     q.Headers["Host"],
		Path:     q.Path,
		Method:   q.HTTPMethod,
		Headers:  q.Headers,
		Encoding: encoding,
		Body:     q.Body,
	})
	return events.APIGatewayProxyResponse{
		StatusCode:        r.StatusCode,
		MultiValueHeaders: r.Headers,
		Body:              r.Body,
		IsBase64Encoded:   true,
	}, err
}

func main() {
	lambda.Start(awshandler)
}
