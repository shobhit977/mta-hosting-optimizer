package service

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mta-hosting-optimizer/lib/aws/s3"
	"github.com/mta-hosting-optimizer/lib/constants"
	errorlib "github.com/mta-hosting-optimizer/lib/errorLib"
	"github.com/mta-hosting-optimizer/lib/models"
)

type Service struct {
	Sess *session.Session
	S3   s3.Interface
}

func NewService() (Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region)},
	)
	if err != nil {
		return Service{}, err
	}
	s3Client := s3.NewService(sess)
	return Service{
		Sess: sess,
		S3:   s3Client,
	}, nil
}

func SuccessResponse(resp models.ServerResponse) events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusOK,
	}
}

func ErrorResponse(errResp errorlib.Error) events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(models.ErrorResponse{
		Error: errResp.Error(),
		Code:  errResp.StatusCode(),
	})
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: errResp.StatusCode(),
	}
}

func MockSuccessResponse() events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Success",
	})
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusCreated,
	}
}
