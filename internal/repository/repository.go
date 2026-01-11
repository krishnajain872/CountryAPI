package repository

import (
	"context"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

type CountryRepository interface {
	FindByName(ctx context.Context, name string) (*types.Country, error)
}