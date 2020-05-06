package error

// ErrInvalidStorageOperation occurs when an unacceptable operation
// is happened during interaction with a storage. For instance,
// if we try to save an entity with a unique parameter
// that already exists in the storage.
type ErrInvalidStorageOperation interface {
	// Error returns an error message
	Error() string

	// Operation returns the name of an operation, which is failed
	Operation() string
}
