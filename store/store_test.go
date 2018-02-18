package store

import (
	"reflect"
	"testing"
)

func TestSetGet(t *testing.T) {
	s := &Store{
		Data: map[string]interface{}{},
	}
	tests := []struct {
		Key           string
		Value         interface{}
		BeingSet      bool
		ExpectedError bool
	}{
		{"name", "John Doe", true, false},
		{"hobbies", []string{"web", "sport"}, true, false},
		{"hello world", map[string]string{"programming": "Golang"}, true, false},
		{"unset", "this item will not being set", false, true},
	}

	for _, test := range tests {
		if test.BeingSet {
			s.Set(test.Key, test.Value)
		}
		v, ok := s.Get(test.Key)
		if ok && test.ExpectedError || !ok && !test.ExpectedError {
			t.Errorf("Expected error: %t, received: %v", test.ExpectedError, ok)
		}
		if v != nil && !reflect.DeepEqual(v, test.Value) {
			t.Errorf("Unequal values. Expected: %v, received: %v", test.Value, v)
		}
	}
}
