package errs

import (
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// Response represents a structured error response with metadata for debugging and tracking.
type Response struct {
	ErrorReference uuid.UUID `json:"error_reference"`
	ErrorCode      ErrorCode `json:"error_code"`
	Code           ErrorCode `json:"-"`
	ErrorType      string    `json:"error_type"`
	Message        string    `json:"message"`
	Err            any       `json:"internal_error_message"`
	StackTrace     string    `json:"-"`
	File           string    `json:"-"`
	Line           int       `json:"-"`
	TimeStamp      string    `json:"-"`
}

// Error returns the formatted error string for the Response, implementing the error interface.
func (e *Response) Error() string {
	return e.Format()
}

// Format returns a detailed string representation of the error, including reference, type, message, file, line, and stack trace.
func (e *Response) Format() string {
	return fmt.Sprintf("%s:%s | %s:%s | %s:%d | stackTrace:%s", e.ErrorReference, e.Err, e.ErrorType, e.Message, e.File, e.Line, e.StackTrace)
}

// Body creates a new error response with the given error code and error, capturing file, line, and timestamp.
func Body(code ErrorCode, err error) error {
	_, file, line, _ := runtime.Caller(1)
	errorResponse := &Response{
		ErrorReference: uuid.New(),
		ErrorCode:      code,
		ErrorType:      errorTypes[code],
		Message:        fmt.Sprintf("error: %s: %v", errorMessages[code], err),
		File:           file,
		Line:           line,
		TimeStamp:      time.Now().Format(time.RFC3339),
	}

	return errorResponse
}

// Message returns the error message string associated with the given error code.
func Message(code ErrorCode) string {
	return errorMessages[code]
}
