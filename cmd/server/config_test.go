package main

import (
	"testing"
)

type TestConfigDriver struct {
	Data []byte
}

func (c *TestConfigDriver) Get() (interface{}, error) {
	return c.Data, nil
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		ConfigDriver  *TestConfigDriver
		ExpectedError bool
	}{
		{&TestConfigDriver{[]byte(`{"port": 1000}`)}, false},
		{&TestConfigDriver{[]byte(`Invaid json data`)}, true},
	}

	for _, test := range tests {
		_, err := NewConfig(test.ConfigDriver)
		if err == nil && test.ExpectedError || err != nil && !test.ExpectedError {
			t.Errorf("Expected error: %t, received: %v", test.ExpectedError, err)
		}
	}
}
