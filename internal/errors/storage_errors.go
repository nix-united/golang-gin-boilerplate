package errors

type ErrInvalidStorageOperation interface {
	Error() string
	Operation() string
}
