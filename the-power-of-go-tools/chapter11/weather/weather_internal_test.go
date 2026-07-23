// This is testing unexported things from weather.go which is in the weather
// package.
package weather

import (
	"testing"
)

func Test_BasicCountryCodeValidationWorks(t *testing.T) {
	t.Parallel()

	err := validateCountryCode("A")
	if err == nil {
		t.Fatal("expected error but nil")
	}

	err = validateCountryCode("AU")
	if err != nil {
		t.Fatal("expected no error but got %w", err)
	}

	err = validateCountryCode("AUS")
	if err != nil {
		t.Fatal("expected no error but got %w", err)
	}

	err = validateCountryCode("AUST")
	if err == nil {
		t.Fatal("expected error but nil")
	}
}
