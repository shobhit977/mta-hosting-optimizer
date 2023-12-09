package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mta-hosting-optimizer/lib/constants"
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
	err = addIpConfig(svc, req)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 201, Body: "Success"}, nil
}

func generateIpConfigOutput(svc service.Service, IpConfigData models.IpConfig, existingInfo []models.IpConfig) []byte {
	ipConfigData := append(existingInfo, IpConfigData)
	allIpConfigBytes, _ := json.Marshal(ipConfigData)
	return allIpConfigBytes
}

func addIpConfig(svc service.Service, req events.APIGatewayV2HTTPRequest) (err error) {
	//return error if body is empty
	if req.Body == "" {
		return errors.New("request body cannot be empty. Please provide valid data")
	}
	var request models.IpConfig
	var ipConfigBytes []byte
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		log.Printf("%v", err)
		return err
	}
	// add validation
	if isFileExist(svc) {
		existingInfo, err := getExistingIpConfigData(svc)
		if err != nil {
			return err
		}
		ipConfigBytes = generateIpConfigOutput(svc, request, existingInfo)
	} else {
		ipConfigBytes = generateIpConfigOutput(svc, request, nil)
	}
	err = s3helper.PutS3Object(svc, ipConfigBytes, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func isFileExist(svc service.Service) bool {
	exist, err := s3helper.KeyExists(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
	}
	return exist
}

func getExistingIpConfigData(svc service.Service) ([]models.IpConfig, error) {
	existingInfo, err := s3helper.GetS3Object(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return []models.IpConfig{}, err
	}
	var ipConfig []models.IpConfig
	if err := json.Unmarshal(existingInfo, &ipConfig); err != nil {
		log.Printf("%v", err)
		return []models.IpConfig{}, err
	}
	return ipConfig, nil
}

func main() {
	lambda.Start(handler)
}
