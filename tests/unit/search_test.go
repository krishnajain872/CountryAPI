package unit

import (
	"testing"

	"github.com/krishnajain872/country-search-api-cache/internal/domain"
)

func TestNewSearchRequest_Success(t *testing.T) {
	req, err := domain.NewSearchRequest("India")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "India" {
		t.Fatalf("expected India, got %s", req.Name)
	}
}

func TestNewSearchRequest_Trim(t *testing.T) {
	req, err := domain.NewSearchRequest("  India ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "India" {
		t.Fatalf("expected trimmed name")
	}
}

func TestNewSearchRequest_Invalid(t *testing.T) {
	_, err := domain.NewSearchRequest("I")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCacheKey(t *testing.T) {
	req, _ := domain.NewSearchRequest("India")
	key := domain.CacheKey(req)

	expected := "country:india"
	if key != expected {
		t.Fatalf("expected %s, got %s", expected, key)
	}
}
