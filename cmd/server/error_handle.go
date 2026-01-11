package main

import (
    "fmt"
	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

func handleError(err error) string {
	e := appErr.Get(err)

	return fmt.Sprintf(
		"[ERROR] type=%s status=%d message=%s context=%v",
		e.Type,
		e.StatusCode,
		e.Message,
		e.Context,
	)
}
