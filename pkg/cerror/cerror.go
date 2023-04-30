package cerror

import "net/http"

type (
	CError interface {
		SetStatusCode(code int)
		GetStatusCode() int
		IsValidationError() bool
		IsSystemError() bool
		Error() string
	}

	cError struct {
		s          string
		statusCode int
	}
)

// SetStatusCode status code setter
func (v *cError) SetStatusCode(code int) {
	v.statusCode = code
}

// GetStatusCode status code getter
func (v *cError) GetStatusCode() int {
	return v.statusCode
}

// Error return error message
func (e *cError) Error() string {
	return e.s
}

// IsValidationError return true if validation error
func (e *cError) IsValidationError() bool {
	return e.GetStatusCode() >= http.StatusBadRequest && e.GetStatusCode() < http.StatusInternalServerError
}

// IsSystemError return true if system error
func (e *cError) IsSystemError() bool {
	return e.GetStatusCode() >= http.StatusInternalServerError
}

// NewWithStatusCode used to initialize error with code
func NewWithStatusCode(msg string, code int) CError {
	return &cError{
		s:          msg,
		statusCode: code,
	}
}

// NewWithStatusCode used to initialize error validation
func NewValidationError(msg string) CError {
	return &cError{
		s:          msg,
		statusCode: http.StatusBadRequest,
	}
}

// NewWithStatusCode used to initialize error system
func NewSystemError(msg string) CError {
	return &cError{
		s:          msg,
		statusCode: http.StatusInternalServerError,
	}
}
