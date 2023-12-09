package s3

import (
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
