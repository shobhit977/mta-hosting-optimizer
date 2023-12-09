package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mta-hosting-optimizer/lib/constants"
	errorlib "github.com/mta-hosting-optimizer/lib/errorLib"
	"github.com/mta-hosting-optimizer/lib/models"
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
	svcErr := addIpConfig(svc, req)
	if svcErr != nil {
		return service.ErrorResponse(svcErr), nil
	}

	return service.MockSuccessResponse(), nil
}

// return server data in bytes
func generateIpConfigOutput(svc service.Service, IpConfigData models.IpConfig, existingInfo []models.IpConfig) []byte {
	ipConfigData := append(existingInfo, IpConfigData)
	allIpConfigBytes, _ := json.Marshal(ipConfigData)
	return allIpConfigBytes
}

func addIpConfig(svc service.Service, req events.APIGatewayV2HTTPRequest) errorlib.Error {
	//return error if body is empty
	if req.Body == "" {
		return errorlib.New(errors.New("request body cannot be empty. Please provide valid data"), http.StatusBadRequest)
	}
	//convert request body to go struct
	var request models.IpConfig
	var ipConfigBytes []byte
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		log.Printf("%v", err)
		return errorlib.New(err, http.StatusInternalServerError)
	}
	// if file exist, append server data to existing file in S3, else create a new file
	if isFileExist(svc) {
		existingInfo, err := getExistingIpConfigData(svc)
		if err != nil {
			return err
		}
		ipConfigBytes = generateIpConfigOutput(svc, request, existingInfo)
	} else {
		ipConfigBytes = generateIpConfigOutput(svc, request, nil)
	}
	// add server data to s3 bucket
	err := s3helper.PutS3Object(svc, ipConfigBytes, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return errorlib.New(err, http.StatusInternalServerError)
	}

	return nil
}

// checks if file exists in S3 bucket
func isFileExist(svc service.Service) bool {
	exist, err := s3helper.KeyExists(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
	}
	return exist
}

// get existing server data from S3
func getExistingIpConfigData(svc service.Service) ([]models.IpConfig, errorlib.Error) {
	existingInfo, err := s3helper.GetS3Object(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return nil, errorlib.New(err, http.StatusInternalServerError)
	}
	var ipConfig []models.IpConfig
	if err := json.Unmarshal(existingInfo, &ipConfig); err != nil {
		log.Printf("%v", err)
		return nil, errorlib.New(err, http.StatusInternalServerError)
	}
	return ipConfig, nil
}

func main() {
	lambda.Start(handler)
}
