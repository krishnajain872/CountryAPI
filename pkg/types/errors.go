package types

import (
	"fmt"
	"time"
)

/* =====================
   Error Types
   ===================== */

type ErrorType string

const (
	Validation  ErrorType = "VALIDATION_ERROR"
	NotFound    ErrorType = "NOT_FOUND"
	External    ErrorType = "EXTERNAL_API_ERROR"
	Cache       ErrorType = "CACHE_ERROR"
	Internal    ErrorType = "INTERNAL_ERROR"
	Timeout     ErrorType = "TIMEOUT_ERROR"
	Unavailable ErrorType = "SERVICE_UNAVAILABLE"
)

/* =====================
   AppError
   ===================== */

type AppError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Err        error                  `json:"-"`
	Context    map[string]any         `json:"context,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

/* =====================
   Methods (MUST BE HERE)
   ===================== */

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}
