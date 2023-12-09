package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	dummyS3 "github.com/mta-hosting-optimizer/lib/aws/s3/dummy"
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

func Test_getIPConfig_EmptyBucket_Fail(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, errors.New("NotFound")
			},
		},
		Sess: sess,
	}
	result, err := getIpConfig(svc)
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "server information not found")
	assert.Equal(t, err.StatusCode(), 404)
}
func Test_getIPConfig_GetS3Object_Fail(t *testing.T) {
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
	result, err := getIpConfig(svc)
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "get object failed")
	assert.Equal(t, err.StatusCode(), 500)
}

func Test_getIPConfig_Success(t *testing.T) {
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
	result, err := getIpConfig(svc)
	assert.Equal(t, result, []byte(mockServerJsonData))
	assert.Nil(t, err)
}
