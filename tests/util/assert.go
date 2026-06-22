package util

import (
	"testing"
)

func AssertEquals(t *testing.T, a any, b any) {
	if a != b {
		t.Errorf("Expected \"%v\" to be equal to \"%v\"", a, b)
		t.FailNow()
	}
}

func AssertErrorReturn(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected to return an error")
		t.FailNow()
	}
}
