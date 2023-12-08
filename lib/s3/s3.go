package s3

import (
	"bytes"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Interface interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	HeadObject(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
}

type service struct {
	s3 s3iface.S3API
	Interface
}

func NewService(sess *session.Session) Interface {
	s3 := s3.New(sess)
	return &service{s3: s3}
}

func (svc *service) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return svc.s3.GetObject(input)
}
func (svc *service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return svc.s3.PutObject(input)
}
func (svc *service) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return svc.s3.HeadObject(input)
}

func (svc *service) PutS3Object(byteData []byte, bucket string, key string) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(byteData)),
	}
	_, err := svc.PutObject(params)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func (svc *service) GetS3Object(bucket string, key string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := svc.GetObject(params)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	defer result.Body.Close()

	// capture all bytes from upload
	byteData, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return byteData, nil

}

func (svc *service) KeyExists(bucket string, key string) (bool, error) {
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}
