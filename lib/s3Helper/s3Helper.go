package s3helper

import (
	"bytes"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	s3Svc "github.com/aws/aws-sdk-go/service/s3"
	"github.com/mta-hosting-optimizer/lib/service"
)

func PutS3Object(svc service.Service, byteData []byte, bucket string, key string) error {
	params := &s3Svc.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(byteData)),
	}
	_, err := svc.S3.PutObject(params)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func GetS3Object(svc service.Service, bucket string, key string) ([]byte, error) {
	params := &s3Svc.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := svc.S3.GetObject(params)
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

func KeyExists(svc service.Service, bucket string, key string) (bool, error) {
	_, err := svc.S3.HeadObject(&s3Svc.HeadObjectInput{
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
