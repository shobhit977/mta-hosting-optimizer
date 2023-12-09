package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mta-hosting-optimizer/lib/aws/s3"
	"github.com/mta-hosting-optimizer/lib/constants"
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
