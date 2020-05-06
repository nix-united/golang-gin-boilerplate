package service

type errUserAlreadyExists struct {
	message   string
	operation string
}

// NewErrUserAlreadyExists returns the errUserAlreadyExists error
func NewErrUserAlreadyExists(msg, opName string) errUserAlreadyExists {
	return errUserAlreadyExists{
		message:   msg,
		operation: opName,
	}
}

func (e errUserAlreadyExists) Error() string {
	return e.message
}

func (e errUserAlreadyExists) Operation() string {
	return e.operation
}
