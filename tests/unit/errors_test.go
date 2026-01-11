package unit

import (
	"errors"
	"testing"

	appErr "CountryAPI/pkg/error"
	"CountryAPI/pkg/types"
)

func TestValidationError(t *testing.T) {
	err := appErr.Validation("name", "required")

	if err.Type != types.Validation {
		t.Fatalf("expected %s, got %s", types.Validation, err.Type)
	}
}

func TestGet_FromAppError(t *testing.T) {
	orig := appErr.Validation("name", "required")
	got := appErr.Get(orig)

	if got != orig {
		t.Fatal("expected same error instance")
	}
}

func TestGet_FromStdError(t *testing.T) {
	stdErr := errors.New("oops")
	app := appErr.Get(stdErr)

	if app.Type != types.Internal {
		t.Fatalf("expected %s, got %s", types.Internal, app.Type)
	}
}

func TestError_Unwrap(t *testing.T) {
	root := errors.New("root")
	err := appErr.InternalErr("failed", root)

	if !errors.Is(err, root) {
		t.Fatal("unwrap failed")
	}
}
