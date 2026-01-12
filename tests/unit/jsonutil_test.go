package unit

import (
	"testing"

	jsonutil "github.com/krishnajain872/country-search-api-cache/pkg/utils/json"
)

func TestMarshal_Unmarshal(t *testing.T) {
	type sample struct {
		Name string
		Age  int
	}

	orig := sample{Name: "Krishna", Age: 25}

	jsonStr, err := jsonutil.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	decodedPtr, err := jsonutil.Unmarshal[sample](jsonStr)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	decoded := *decodedPtr

	if decoded != orig {
		t.Fatal("decoded does not match original")
	}
}
