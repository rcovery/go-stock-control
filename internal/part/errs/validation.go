package errs

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}

func NewValidationError(msg string) error {
	return ValidationError{Message: msg}
}
