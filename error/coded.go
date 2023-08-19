package error

type Error interface {
	Code() string
	error
}

type codedError struct {
	code string
	msg  string
}

func (e *codedError) Code() string {
	return e.code
}

func (e *codedError) Error() string {
	return e.msg
}

func New(code string, msg string) error {
	return &codedError{code: code, msg: msg}
}

func GetCode(err error) (code string, ok bool) {
	if err == nil {
		return
	}
	if e, ok := err.(Error); ok {
		code = e.Code()
	}
	return
}
