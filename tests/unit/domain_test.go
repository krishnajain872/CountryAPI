// ========================================
// FILE: tests/unit/domain_test.go
// ========================================
package unit

import (
	"testing"

	"github.com/krishnajain872/country-search-api-cache/internal/domain"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

func TestNewCountry(t *testing.T) {
	c := domain.NewCountry("India", "Delhi", "INR", 100)

	if c.Name != "India" {
		t.Fatal("country name mismatch")
	}
}

func TestValidateCountry_Success(t *testing.T) {
	c := &types.Country{
		Name:       "India",
		Capital:    "Delhi",
		Currency:   "INR",
		Population: 100,
	}

	if err := domain.ValidateCountry(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCountry_InvalidName(t *testing.T) {
	c := &types.Country{
		Name:       "",
		Capital:    "Delhi",
		Currency:   "INR",
		Population: 100,
	}

	if err := domain.ValidateCountry(c); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateCountry_InvalidPopulation(t *testing.T) {
	c := &types.Country{
		Name:       "India",
		Capital:    "Delhi",
		Currency:   "INR",
		Population: -1,
	}

	if err := domain.ValidateCountry(c); err == nil {
		t.Fatal("expected population error")
	}
}
