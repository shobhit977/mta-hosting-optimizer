package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	dummyS3 "github.com/mta-hosting-optimizer/lib/aws/s3/dummy"
	"github.com/mta-hosting-optimizer/lib/constants"
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

func Test_getInefficientServers_Success(t *testing.T) {
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
		},
		Sess: sess,
	}
	os.Setenv(constants.ThresholdKey, "1")
	expectedResult := models.ServerResponse{
		Hostnames: []string{"DummyHostname1"},
	}
	result, err := getInefficientServers(svc)
	assert.Nil(t, err)
	assert.Equal(t, result, expectedResult)

}

func Test_getInefficientServers_MockDataNotFound_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, errors.New("NotFound")
			},
		},
		Sess: sess,
	}
	result, err := getInefficientServers(svc)
	assert.Equal(t, result, models.ServerResponse{})
	assert.Equal(t, err.Error(), "server Information not found")

}
func Test_getInefficientServers_InvalidThreshold_Fail(t *testing.T) {
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
		},
		Sess: sess,
	}
	os.Setenv(constants.ThresholdKey, "dummy")
	result, err := getInefficientServers(svc)
	assert.Equal(t, result, models.ServerResponse{})
	assert.Equal(t, err.Error(), "invalid threshold value")
	assert.Equal(t, err.StatusCode(), 500)

}

func Test_getInefficientServers_GetS3ObjectError_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, errors.New("get s3 object fail")
			},
		},
		Sess: sess,
	}
	result, err := getInefficientServers(svc)
	assert.Equal(t, result, models.ServerResponse{})
	assert.Equal(t, err.Error(), "get s3 object fail")
	assert.Equal(t, err.StatusCode(), 500)

}
func Test_getIpConfigData_InvalidModelStructure_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBufferString("Invalid Data")),
				}, nil
			},
		},
		Sess: sess,
	}
	result, err := getIpConfigData(svc)
	assert.Equal(t, result, []models.IpConfig(nil))
	assert.Error(t, err)
	assert.Equal(t, err.StatusCode(), 500)

}
