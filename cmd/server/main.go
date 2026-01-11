package main

import (
	"fmt"
	"github.com/krishnajain872/country-search-api-cache/internal/domain"
)

func main() {
	// 1️⃣ Create validated search request
	searchReq, err := domain.NewSearchRequest("India")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println("Search Request:", searchReq)
	fmt.Println("Cache Key:", domain.CacheKey(searchReq))

	// 2️⃣ Create country entity
	country := domain.NewCountry(
		"India",
		"New Delhi",
		"INR",
		1_428_000_000,
	)

	// 3️⃣ Validate domain entity
	if err := domain.ValidateCountry(country); err != nil {
		handleError(err)
		return
	}

	fmt.Println("Country:", domain.CountryString(country))
}
