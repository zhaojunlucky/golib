package err

import (
	"errors"
	"testing"
)

func TestCheckErr_WithError(t *testing.T) {
	testErr := errors.New("test error")
	var capturedErr error

	// Define error function that captures the error
	errFun := func(err error) {
		capturedErr = err
	}

	// Call CheckErr with an error
	CheckErr(testErr, errFun)

	// Verify the error function was called with the correct error
	if !errors.Is(capturedErr, testErr) {
		t.Errorf("Expected error function to be called with %v, got %v", testErr, capturedErr)
	}
}

func TestCheckErr_WithNilError(t *testing.T) {
	var functionCalled bool

	// Define error function that sets a flag
	errFun := func(err error) {
		functionCalled = true
	}

	// Call CheckErr with nil error
	CheckErr(nil, errFun)

	// Verify the error function was NOT called
	if functionCalled {
		t.Error("Error function should not be called when error is nil")
	}
}

func TestCheckErr_WithPanicFunction(t *testing.T) {
	testErr := errors.New("panic error")

	// Define error function that panics
	errFun := func(err error) {
		panic("error occurred: " + err.Error())
	}

	// Test that CheckErr properly calls the panic function
	defer func() {
		if r := recover(); r != nil {
			expectedPanic := "error occurred: panic error"
			if r != expectedPanic {
				t.Errorf("Expected panic message '%s', got '%v'", expectedPanic, r)
			}
		} else {
			t.Error("Expected function to panic, but it didn't")
		}
	}()

	CheckErr(testErr, errFun)
}

func TestCheckErr_WithLoggingFunction(t *testing.T) {
	testErr := errors.New("logging test error")
	var loggedMessage string

	// Define error function that logs
	errFun := func(err error) {
		loggedMessage = "Error logged: " + err.Error()
	}

	// Call CheckErr
	CheckErr(testErr, errFun)

	// Verify the logging happened
	expectedMessage := "Error logged: logging test error"
	if loggedMessage != expectedMessage {
		t.Errorf("Expected logged message '%s', got '%s'", expectedMessage, loggedMessage)
	}
}

func TestCheckErr_WithMultipleErrors(t *testing.T) {
	errors := []error{
		errors.New("first error"),
		nil,
		errors.New("second error"),
		nil,
	}

	var capturedErrors []error

	// Define error function that collects errors
	errFun := func(err error) {
		capturedErrors = append(capturedErrors, err)
	}

	// Call CheckErr for each error
	for _, err := range errors {
		CheckErr(err, errFun)
	}

	// Verify only non-nil errors were captured
	expectedCount := 2
	if len(capturedErrors) != expectedCount {
		t.Errorf("Expected %d errors to be captured, got %d", expectedCount, len(capturedErrors))
	}

	if capturedErrors[0].Error() != "first error" {
		t.Errorf("Expected first captured error to be 'first error', got '%s'", capturedErrors[0].Error())
	}

	if capturedErrors[1].Error() != "second error" {
		t.Errorf("Expected second captured error to be 'second error', got '%s'", capturedErrors[1].Error())
	}
}

func TestCheckErr_WithNilErrorFunction(t *testing.T) {
	testErr := errors.New("test error")

	// This should panic when error function is nil and there's an error
	defer func() {
		if r := recover(); r == nil {
			t.Error("CheckErr should panic when error function is nil and there's an error")
		}
	}()

	CheckErr(testErr, nil)
}

func TestCheckErr_WithNilErrorFunction_NoError(t *testing.T) {
	// This should not panic when error is nil, even if error function is nil
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CheckErr should not panic when error is nil, even with nil error function, but panicked: %v", r)
		}
	}()

	CheckErr(nil, nil)
}

func TestCheckErr_ErrorFunctionReceivesCorrectError(t *testing.T) {
	originalErr := errors.New("original error message")
	var receivedErr error

	errFun := func(err error) {
		receivedErr = err
	}

	CheckErr(originalErr, errFun)

	// Verify the exact same error instance is passed
	if receivedErr != originalErr {
		t.Error("Error function should receive the exact same error instance")
	}

	// Verify error message is preserved
	if receivedErr.Error() != "original error message" {
		t.Errorf("Expected error message 'original error message', got '%s'", receivedErr.Error())
	}
}
