package s3helper

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

func Test_PutS3Object_Success(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyPutObject: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		Sess: sess,
	}
	err := PutS3Object(svc, []byte("testdata"), "dummy", "dummy")
	assert.Nil(t, err)

}
func Test_PutS3Object_Failure(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyPutObject: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, errors.New("put object failure")
			},
		},
		Sess: sess,
	}
	err := PutS3Object(svc, []byte("testdata"), "dummy", "dummy")
	assert.Equal(t, err.Error(), "put object failure")

}

func Test_GetS3Object_Success(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBufferString("Dummy Data")),
				}, nil
			},
		},
		Sess: sess,
	}
	result, err := GetS3Object(svc, "dummy", "dummy")
	assert.Nil(t, err)
	assert.Equal(t, result, []byte("Dummy Data"))

}
func Test_GetS3Object_Failure(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyGetObject: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, errors.New("get s3 object fail")
			},
		},
		Sess: sess,
	}
	result, err := GetS3Object(svc, "dummy", "dummy")
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "get s3 object fail")

}

func Test_KeyExists_Success(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
		},
		Sess: sess,
	}
	result, err := KeyExists(svc, "dummy", "dummy")
	assert.Nil(t, err)
	assert.Equal(t, result, true)

}
func Test_KeyExists_Failure(t *testing.T) {
	sess, _ := session.NewSession()
	svc := service.Service{
		S3: dummyS3.S3Interface{
			DummyHeadObject: func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, errors.New("NotFound")
			},
		},
		Sess: sess,
	}
	result, err := KeyExists(svc, "dummy", "dummy")
	assert.Equal(t, err.Error(), "NotFound")
	assert.Equal(t, result, false)

}
