package types

import "errors"

var (
	ErrCacheMiss           = errors.New("cache miss")
	ErrInvalidInput        = errors.New("invalid input")
	ErrNotFound            = errors.New("resource not found")
	ErrCache               = errors.New("cache error")
	ErrExternalAPI         = errors.New("external api error")
	ErrTimeout             = errors.New("request timeout")
	ErrInternal            = errors.New("internal server error")
	ErrServiceUnavailable  = errors.New("service unavailable")
)
