// ======================================== 
// FILE: internal/repository/rest_repository.go  
// ======================================== 
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/pkg/httpclient"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

type RestCountriesRepo struct {
	client *httpclient.Client
	logger *logger.Logger
}

func NewRestCountriesRepo(client *httpclient.Client, log *logger.Logger) *RestCountriesRepo {
	return &RestCountriesRepo{
		client: client,
		logger: log,
	}
}

type APIResponse []struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Capital    []string `json:"capital"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Population int64 `json:"population"`
}

func (r *RestCountriesRepo) FindByName(ctx context.Context, name string) (*types.Country, error) {
	r.logger.Info("Fetching country from API", 
		logger.String("country", name))

	// ✅ URL encode the country name to handle spaces and special characters
	encodedName := url.QueryEscape(name)
	
	r.logger.Debug("URL encoding applied",
		logger.String("original", name),
		logger.String("encoded", encodedName))

	path := fmt.Sprintf("/name/%s?fullText=false", encodedName)

	body, err := r.client.Get(ctx, path)
	if err != nil {
		r.logger.Error("Failed to fetch from API",
			logger.String("country", name),
			logger.String("error", err.Error()))
		return nil, appErr.External("failed to fetch country data", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		r.logger.Error("Failed to parse API response",
			logger.String("error", err.Error()))
		return nil, appErr.External("failed to parse API response", err)
	}

	if len(apiResp) == 0 {
		return nil, appErr.ResourceNotFound("country")
	}

	countryData := apiResp[0]

	capital := ""
	if len(countryData.Capital) > 0 {
		capital = countryData.Capital[0]
	}

	currency := ""
	for _, curr := range countryData.Currencies {
		currency = curr.Symbol
		break
	}

	country := &types.Country{
		Name:       countryData.Name.Common,
		Capital:    capital,
		Currency:   currency,
		Population: countryData.Population,
	}

	r.logger.Info("Successfully fetched country",
		logger.String("country", country.Name))

	return country, nil
}