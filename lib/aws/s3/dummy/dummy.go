package dummy

import (
	s3Svc "github.com/aws/aws-sdk-go/service/s3"
	"github.com/mta-hosting-optimizer/lib/aws/s3"
)

type S3Interface struct {
	DummyGetObject  func(*s3Svc.GetObjectInput) (*s3Svc.GetObjectOutput, error)
	DummyPutObject  func(*s3Svc.PutObjectInput) (*s3Svc.PutObjectOutput, error)
	DummyHeadObject func(*s3Svc.HeadObjectInput) (*s3Svc.HeadObjectOutput, error)
}

var _ s3.Interface = &S3Interface{}

func (d S3Interface) GetObject(input *s3Svc.GetObjectInput) (*s3Svc.GetObjectOutput, error) {
	return d.DummyGetObject(input)
}
func (d S3Interface) PutObject(input *s3Svc.PutObjectInput) (*s3Svc.PutObjectOutput, error) {
	return d.DummyPutObject(input)
}
func (d S3Interface) HeadObject(input *s3Svc.HeadObjectInput) (*s3Svc.HeadObjectOutput, error) {
	return d.DummyHeadObject(input)
}
