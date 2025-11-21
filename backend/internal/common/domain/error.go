package domain

import (
	"bytes"
	"fmt"
)

// --- 1. Strong Types ---

// Code defines the category of the error.
type Code string

// Op defines the logical operation (usually "struct.method").
type Op string

// Predefined Domain Error Codes.
const (
	EInternal  Code = "internal"
	EInvalid   Code = "invalid"
	ENotFound  Code = "not_found"
	EConflict  Code = "conflict"
	EForbidden Code = "forbidden"
)

// --- 2. The Error Struct ---

// Error represents a domain error.
// We export the fields so standard JSON marshaling works if needed,
// but you should use the Builder methods to construct it.
type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Op      Op     `json:"op,omitempty"`
	Err     error  `json:"-"` // Don't serialize the raw internal error to clients
}

// --- 3. The Builder Constructor ---

// NewError creates a new domain error with a Code and a Message.
// This is the entry point.
func NewError(code Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// At attaches the operation context to the error.
// Usage: domain.New(...).At("UserService.Create")
func (e *Error) At(op Op) *Error {
	e.Op = op
	return e
}

// CausedBy attaches the underlying technical error (e.g., database error).
// Usage: domain.New(...).CausedBy(err)
func (e *Error) CausedBy(err error) *Error {
	e.Err = err
	return e
}

// --- 4. Interface Implementations ---

// Error satisfies the stdlib error interface.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Format: <Op>: <Code>: <Message>
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	if e.Code != "" {
		fmt.Fprintf(&buf, "<%s> ", e.Code)
	}

	buf.WriteString(e.Message)

	// We do NOT print e.Err here to keep the message clean for end-users.
	// Use Unwrap() or logging middleware to see the underlying error.
	return buf.String()
}

// Unwrap allows standard library errors.Is/As to work.
func (e *Error) Unwrap() error {
	return e.Err
}

// Helper: ErrorCode safely extracts the code from any error.
func ErrorCode(err error) Code {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return EInternal
}
