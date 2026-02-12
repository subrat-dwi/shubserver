package utils

// ValidationError represents an error that occurs when input validation fails
type ValidationError struct {
	Message string
}

// Error returns the error message for the ValidationError
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new instance of ValidationError with the provided message
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

// IsValidationError checks if the given error is of type ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
