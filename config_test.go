package main

import (
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	for _, test := range []struct {
		fn          string
		expect      Config
		expectError bool
	}{
		{"testdata/nonsuch", Config{}, true},
		{"testdata/valid-config.toml", Config{CheckPTR: true, MapOwners: true, Owners: map[string]Owner{"IlzwX7BAAY1Dll7P4VLVmHUcA/h7mXw=": {Name: "Bill Sticklebricks", Email: "bs@example.com"}, "adsdggfewgrw3rDDdsdUc3/hdmaw=": {Name: "Fatima Fuzzy-Felt", Email: "fff@example.com"}}}, false},
	} {
		t.Run(test.fn, func(t *testing.T) {
			c, err := ParseConfig(test.fn)
			if err != nil && !test.expectError {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && test.expectError {
				t.Error("expected error")
			}

			if !reflect.DeepEqual(test.expect, c) {
				t.Errorf("expected\n%#v\nreceived\n%#v", test.expect, c)
			}
		})
	}
}
