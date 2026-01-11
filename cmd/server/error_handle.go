package main

import (
	"log"

	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

func handleError(err error) {
	e := appErr.Get(err)

	log.Printf(
		"[ERROR] type=%s status=%d message=%s context=%v",
		e.Type,
		e.StatusCode,
		e.Message,
		e.Context,
	)
}
