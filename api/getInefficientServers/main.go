package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

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
	inefficientServers, svcErr := getInefficientServers(svc)
	if svcErr != nil {
		return service.ErrorResponse(svcErr), nil
	}
	if len(inefficientServers.Hostnames) == 0 {
		svcErr := errorlib.New(errors.New("no inefficient servers found as per threshold"), http.StatusNotFound)
		return service.ErrorResponse(svcErr), nil
	}
	return service.SuccessResponse(inefficientServers), nil
}

func getInefficientServers(svc service.Service) (models.ServerResponse, errorlib.Error) {
	// get server data from s3 bucket
	ipConfig, svcErr := getIpConfigData(svc)
	if svcErr != nil {
		return models.ServerResponse{}, svcErr
	}
	// get servers with active MTA information
	activeIpConfig := makeIpConfigMap(ipConfig)
	// convert threshold to integer
	threshold, err := strconv.ParseInt(os.Getenv(constants.ThresholdKey), 10, 32)
	if err != nil {
		log.Printf("%v", err)
		return models.ServerResponse{}, errorlib.New(errors.New("invalid threshold value"), http.StatusInternalServerError)
	}
	var inefficientHostnames []string
	// get servers whose active MTAs is less than or equal to threshold
	for hostname, activeMTAs := range activeIpConfig {
		if activeMTAs <= int(threshold) {
			inefficientHostnames = append(inefficientHostnames, hostname)
		}
	}
	return models.ServerResponse{
		Hostnames: inefficientHostnames,
	}, nil
}

// make map of server with active MTA information
func makeIpConfigMap(ipConfig []models.IpConfig) map[string]int {
	serverMap := make(map[string]int)
	for _, val := range ipConfig {
		if val.Active {
			serverMap[val.Hostname]++
		} else {
			if _, ok := serverMap[val.Hostname]; !ok {
				serverMap[val.Hostname] = 0
			}
		}
	}
	return serverMap
}

func getIpConfigData(svc service.Service) ([]models.IpConfig, errorlib.Error) {
	// return error if mock data is not present in s3 bucket
	if !isFileExist(svc) {
		return nil, errorlib.New(errors.New("server Information not found"), http.StatusNotFound)
	}
	// get server data from file in s3 bucker
	ipConfig, err := s3helper.GetS3Object(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return nil, errorlib.New(err, http.StatusInternalServerError)
	}
	var ipConfigData []models.IpConfig
	if err := json.Unmarshal(ipConfig, &ipConfigData); err != nil {
		log.Printf("%v", err)
		return nil, errorlib.New(err, http.StatusInternalServerError)
	}
	return ipConfigData, nil
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
