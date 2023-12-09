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
	ipConfig, svcErr := getIpConfigData(svc)
	if svcErr != nil {
		return models.ServerResponse{}, svcErr
	}
	activeIpConfig := makeIpConfigMap(ipConfig)
	threshold, err := strconv.ParseInt(os.Getenv(constants.ThresholdKey), 10, 32)
	if err != nil {
		log.Printf("%v", err)
		return models.ServerResponse{}, errorlib.New(errors.New("invalid threshold value"), http.StatusInternalServerError)
	}
	var inefficientHostnames []string
	for k, v := range activeIpConfig {
		if v <= int(threshold) {
			inefficientHostnames = append(inefficientHostnames, k)
		}
	}
	return models.ServerResponse{
		Hostnames: inefficientHostnames,
	}, nil
}

func makeIpConfigMap(ipConfig []models.IpConfig) map[string]int {
	activeIpMap := make(map[string]int)
	for _, v := range ipConfig {
		if v.Active {
			activeIpMap[v.Hostname]++
		}
	}
	return activeIpMap
}

func getIpConfigData(svc service.Service) ([]models.IpConfig, errorlib.Error) {
	if !isFileExist(svc) {
		return nil, errorlib.New(errors.New("server Information not found"), http.StatusNotFound)
	}
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
