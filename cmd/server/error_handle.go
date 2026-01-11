package main

import (
	"log"

	appErr "CountryAPI/pkg/error"
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
