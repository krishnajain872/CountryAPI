package unit

import (
	"testing"

	"github.com/krishnajain872/country-search-api-cache/internal/domain"
)

func TestNewCountry(t *testing.T) {
	c := domain.NewCountry("India", "Delhi", "INR", 100)

	if c.Name != "India" {
		t.Fatal("country name mismatch")
	}
}

func TestValidateCountry_Success(t *testing.T) {
	c := domain.NewCountry("India", "Delhi", "INR", 100)

	if err := domain.ValidateCountry(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCountry_InvalidName(t *testing.T) {
	c := domain.NewCountry("", "Delhi", "INR", 100)

	if err := domain.ValidateCountry(c); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateCountry_InvalidPopulation(t *testing.T) {
	c := domain.NewCountry("India", "Delhi", "INR", -1)

	if err := domain.ValidateCountry(c); err == nil {
		t.Fatal("expected population error")
	}
}
