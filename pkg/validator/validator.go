package validator

import (
	"fmt"
	"regexp"
	"strings"

	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

type Rule func() error

func Validate(rules ...Rule) error {
	for _, rule := range rules {
		if err := rule(); err != nil {
			return err
		}
	}
	return nil
}

/* =========================
   Common Validation Rules
   ========================= */

func Required(field, value string) Rule {
	return func() error {
		if strings.TrimSpace(value) == "" {
			return appErr.Validation(field, "is required")
		}
		return nil
	}
}

func MinLength(field, value string, min int) Rule {
	return func() error {
		if len(value) < min {
			return appErr.Validation(
				field,
				fmt.Sprintf("must be at least %d characters", min),
			)
		}
		return nil
	}
}

func MaxLength(field, value string, max int) Rule {
	return func() error {
		if len(value) > max {
			return appErr.Validation(
				field,
				fmt.Sprintf("must not exceed %d characters", max),
			)
		}
		return nil
	}
}

/* =========================
   Numeric Rules
   ========================= */


// Min validates that a numeric value is >= min
func Min[T OrderedNumeric](field string, value, min T) Rule {
	return func() error {
		if value < min {
			return appErr.Validation(
				field,
				fmt.Sprintf("must be greater than or equal to %v", min),
			)
		}
		return nil
	}
}
/* =========================
   Pattern-based Rules
   ========================= */

var countryNamePattern = regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s\-'.()]+$`)

func CountryName(field, value string) Rule {
	return func() error {
		if !countryNamePattern.MatchString(value) {
			return appErr.Validation(field, "contains invalid characters")
		}
		return nil
	}
}
