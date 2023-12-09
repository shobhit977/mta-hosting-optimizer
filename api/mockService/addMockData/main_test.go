package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	dummyS3 "github.com/mta-hosting-optimizer/lib/aws/s3/dummy"
	"github.com/mta-hosting-optimizer/lib/models"
	"github.com/mta-hosting-optimizer/lib/service"
	"github.com/stretchr/testify/assert"
)

var mockServerJsonData = `[{
	"ip":"DummyIP1",
	"hostname":"DummyHostname1",
	"active": true
},{
	"ip":"DummyIP2",
	"hostname":"DummyHostname2",
	"active": false
}]`

func Test_addIpConfig_EmptyBody_Fail(t *testing.T) {
	svc := service.Service{}
	req := events.APIGatewayV2HTTPRequest{}
	err := addIpConfig(svc, req)
	assert.Equal(t, err.Error(), "request body cannot be empty. Please provide valid data")
	assert.Equal(t, err.StatusCode(), 400)
}

func Test_addIpConfig_InvalidRequest_Fail(t *testing.T) {
	svc := service.Service{}
	req := events.APIGatewayV2HTTPRequest{Body: "Invalid Data"}
	err := addIpConfig(svc, req)
	assert.Equal(t, err.StatusCode(), 500)

}
func Test_addIpConfig_GetMockData_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, errors.New("get object failed")
			},
		},
		Sess: sess,
	}
	req := events.APIGatewayV2HTTPRequest{Body: `{
		"ip":"DummyIP1",
		"hostname":"DummyHostname1",
		"active": true
	}`}
	err := addIpConfig(svc, req)
	assert.Equal(t, err.Error(), "get object failed")
	assert.Equal(t, err.StatusCode(), 500)

}
func Test_getExistingIpConfigData_InvalidMockData_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBufferString("InvalidData")),
				}, nil
			},
		},
		Sess: sess,
	}

	result, err := getExistingIpConfigData(svc)
	assert.Equal(t, result, []models.IpConfig(nil))
	assert.Equal(t, err.StatusCode(), 500)

}
func Test_addIPConfig_PutObject_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBufferString(mockServerJsonData)),
				}, nil
			},
			DummyPutObject: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, errors.New("put s3 object failed")
			},
		},
		Sess: sess,
	}
	req := events.APIGatewayV2HTTPRequest{Body: `{
		"ip":"DummyIP1",
		"hostname":"DummyHostname1",
		"active": true
	}`}
	err := addIpConfig(svc, req)
	assert.Equal(t, err.Error(), "put s3 object failed")
	assert.Equal(t, err.StatusCode(), 500)

}

func Test_addIPConfig_EmptyS3Bucket_Success(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, errors.New("NotFound")
			},
			DummyPutObject: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		Sess: sess,
	}
	req := events.APIGatewayV2HTTPRequest{Body: `{
		"ip":"DummyIP1",
		"hostname":"DummyHostname1",
		"active": true
	}`}
	err := addIpConfig(svc, req)
	assert.Nil(t, err)

}

func Test_addIPConfig_ExistingDataInS3Bucket_Success(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBufferString(mockServerJsonData)),
				}, nil
			},
			DummyPutObject: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		Sess: sess,
	}
	req := events.APIGatewayV2HTTPRequest{Body: `{
		"ip":"DummyIP1",
		"hostname":"DummyHostname1",
		"active": true
	}`}
	err := addIpConfig(svc, req)
	assert.Nil(t, err)

}
