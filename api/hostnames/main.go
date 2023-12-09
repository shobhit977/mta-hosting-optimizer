package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

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
	hostnames, err := getInefficientHostnames(svc, req)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return service.SuccessResponse(hostnames), nil
}

func getInefficientHostnames(svc service.Service, req events.APIGatewayV2HTTPRequest) (hostNamesResp models.HostnameResponse, err error) {
	ipConfig, err := getIpConfigData(svc)
	if err != nil {
		return hostNamesResp, err
	}
	activeIpConfig := makeIpConfigMap(ipConfig)
	threshold, err := strconv.ParseInt(os.Getenv(constants.ThresholdKey), 10, 32)
	if err != nil {
		log.Printf("%v", err)
		return hostNamesResp, errors.New("invalid threshold value")
	}
	var inefficientHostnames []string
	for k, v := range activeIpConfig {
		if v <= int(threshold) {
			inefficientHostnames = append(inefficientHostnames, k)
		}
	}
	hostNamesResp.Hostnames = inefficientHostnames
	return hostNamesResp, nil
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

func getIpConfigData(svc service.Service) ([]models.IpConfig, error) {
	if !isFileExist(svc) {
		return nil, errors.New("ip configuration data not found")
	}
	ipConfig, err := s3helper.GetS3Object(svc, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	var ipConfigData []models.IpConfig
	if err := json.Unmarshal(ipConfig, &ipConfigData); err != nil {
		log.Printf("%v", err)
		return []models.IpConfig{}, err
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
