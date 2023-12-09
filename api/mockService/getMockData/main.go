package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mta-hosting-optimizer/lib/constants"
	errorlib "github.com/mta-hosting-optimizer/lib/errorLib"
	s3helper "github.com/mta-hosting-optimizer/lib/s3Helper"
	"github.com/mta-hosting-optimizer/lib/service"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	svc, err := service.NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body: err.Error(),
		}, nil
	}
	ipConfig, svcErr := getIpConfig(svc)
	if svcErr != nil {
		return service.ErrorResponse(svcErr), nil
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusOK, Body: string(ipConfig)}, nil
}

func getIpConfig(svc service.Service) ([]byte, errorlib.Error) {
	// return error if file does not exist in s3
	if !isFileExist(svc) {
		return nil, errorlib.New(errors.New("server information not found"), http.StatusNotFound)
	}
	// get mock data from s3 bucket
	ipConfig, err := s3helper.GetS3Object(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return nil, errorlib.New(err, http.StatusInternalServerError)
	}
	return ipConfig, nil
}

// checks if file exists in S3 bucket
func isFileExist(svc service.Service) bool {
	exist, err := s3helper.KeyExists(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
	}
	return exist
}

func main() {
	lambda.Start(handler)
}
