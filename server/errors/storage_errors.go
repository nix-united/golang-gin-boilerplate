package error

type ErrInvalidStorageOperation interface {
	Error() string
	Operation() string
}
