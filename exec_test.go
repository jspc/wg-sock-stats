package main

import (
	"testing"
)

func TestWGDump(t *testing.T) {
	oldBinary := binary
	defer func() {
		binary = oldBinary
	}()

	for _, test := range []struct {
		binary      string
		expect      string
		expectError bool
	}{
		{"this-command-doesnt-exist", "", true},
		{"echo", "show all dump\n", false},
	} {
		t.Run(test.binary, func(t *testing.T) {
			binary = test.binary
			received, err := WGDump()
			if err != nil && !test.expectError {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && test.expectError {
				t.Error("expected error")
			}

			if test.expect != string(received) {
				t.Errorf("expected %q, received %q", test.expect, received)
			}
		})
	}
}
