package domain

import (
	"fmt"
	"strings"

	"CountryAPI/pkg/validator"
	"CountryAPI/pkg/types"
)

// NewSearchRequest creates a validated SearchRequest
func NewSearchRequest(name string) (*types.SearchRequest, error) {
	name = strings.TrimSpace(name)

	if err := validator.Validate(
		validator.Required("name", name),
		validator.MinLength("name", name, 2),
		validator.MaxLength("name", name, 100),
		validator.CountryName("name", name),
	); err != nil {
		return nil, err
	}

	return &types.SearchRequest{Name: name}, nil
}

// CacheKey generates cache key
func CacheKey(req *types.SearchRequest) string {
	return fmt.Sprintf("country:%s", strings.ToLower(req.Name))
}
