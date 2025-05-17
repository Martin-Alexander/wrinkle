//go:build integration
// +build integration

package integration

import (
	"testing"
)

func TestNothing(t *testing.T) {
	actual := 1
	expected := 1

	if actual != expected {
		t.Errorf("Expected %d, but got %d", expected, actual)
	}
}
