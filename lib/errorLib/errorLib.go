package errorlib

type svcError struct {
	error      error
	statusCode int
}
type Error interface {
	Error() string
	StatusCode() int
}

func New(err error, code int) Error {
	return &svcError{
		error:      err,
		statusCode: code,
	}
}
func (e *svcError) Error() string {
	return e.error.Error()
}

func (e *svcError) StatusCode() int {
	return e.statusCode
}
