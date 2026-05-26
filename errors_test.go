package jsonrpc

import (
	"errors"
	"testing"
)

func TestValidationErrors(t *testing.T) {
	var v ValidationErrors

	if v.HasErrors() {
		t.Errorf("Expected HasErrors to be false, got true")
	}

	err1 := errors.New("error 1")
	v.Append(err1)

	if !v.HasErrors() {
		t.Errorf("Expected HasErrors to be true, got false")
	}

	if len(v.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(v.Errors))
	}

	err2 := errors.New("error 2")
	v.Append(err2)

	if len(v.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(v.Errors))
	}

	expectedMsg := "error 1, error 2"
	if v.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, v.Error())
	}

	unwrapped := v.Unwrap()
	if len(unwrapped) != 2 {
		t.Errorf("Expected 2 unwrapped errors, got %d", len(unwrapped))
	}
	if unwrapped[0] != err1 || unwrapped[1] != err2 {
		t.Errorf("Unwrapped errors do not match original errors")
	}
}
