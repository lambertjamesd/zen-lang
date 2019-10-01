package test

import (
	"testing"
)

func Assert(t *testing.T, value bool, message string) {
	if !value {
		t.Errorf("Failed to assertion %s", message)
	}
}
