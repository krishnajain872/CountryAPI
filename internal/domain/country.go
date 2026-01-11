package domain

import (
	"fmt"

	"CountryAPI/pkg/validator"
	"CountryAPI/pkg/types"
)

// NewCountry creates a new Country entity
func NewCountry(
	name string,
	capital string,
	currency string,
	population int64,
) *types.Country {
	return &types.Country{
		Name:       name,
		Capital:    capital,
		Currency:   currency,
		Population: population,
	}
}

// ValidateCountry applies business validation
func ValidateCountry(c *types.Country) error {
	return validator.Validate(
		validator.Required("name", c.Name),
		validator.Min("population", c.Population, 1),
	)
}

// String returns string representation
func CountryString(c *types.Country) string {
	return fmt.Sprintf(
		"Country{Name: %s, Capital: %s, Currency: %s, Population: %d}",
		c.Name,
		c.Capital,
		c.Currency,
		c.Population,
	)
}
