package unit

import (
	"testing"

	"github.com/krishnajain872/country-search-api-cache/pkg/validator"
)

func TestValidate_Success(t *testing.T) {
	err := validator.Validate(
		validator.Required("name", "India"),
		validator.MinLength("name", "India", 2),
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidate_Required(t *testing.T) {
	err := validator.Validate(
		validator.Required("name", ""),
	)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidate_MinLength(t *testing.T) {
	err := validator.Validate(
		validator.MinLength("name", "I", 2),
	)
	if err == nil {
		t.Fatal("expected min length error")
	}
}

func TestValidate_MaxLength(t *testing.T) {
	long := make([]byte, 101)
	for i := range long {
		long[i] = 'a'
	}
	err := validator.Validate(
		validator.MaxLength("name", string(long), 100),
	)
	if err == nil {
		t.Fatal("expected max length error")
	}
}

func TestValidate_CountryName(t *testing.T) {
	err := validator.Validate(
		validator.CountryName("name", "Ind!a"),
	)
	if err == nil {
		t.Fatal("expected invalid character error")
	}
}

func TestValidate_Min(t *testing.T) {
	err := validator.Validate(
		validator.Min("population", int64(-1), int64(0)),
	)
	if err == nil {
		t.Fatal("expected min int error")
	}
}
