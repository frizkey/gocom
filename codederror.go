package gocom

type CodedError struct {
	Code    int
	Message string
}

func (o *CodedError) Error() string {
	return o.Message
}

func NewError(code int, msg string) *CodedError {

	return &CodedError{
		Code:    code,
		Message: msg,
	}
}
