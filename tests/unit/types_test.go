package unit

import (
	"testing"

	"CountryAPI/pkg/types"
)

func TestSearchRequest(t *testing.T) {
	req := types.SearchRequest{Name: "India"}

	if req.Name != "India" {
		t.Fatal("name mismatch")
	}
}

func TestCountry(t *testing.T) {
	c := types.Country{
		Name:       "India",
		Capital:    "Delhi",
		Currency:   "INR",
		Population: 100,
	}

	if c.Population <= 0 {
		t.Fatal("invalid population")
	}
}
