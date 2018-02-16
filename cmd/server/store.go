package main

import (
	"fmt"
	"sync"
)

//Store - key-value storage implementation
type Store struct {
	sync.RWMutex
	AuthorizedTokens []string
	Driver           StoreDriver
}

//StoreDriver - interface for store
type StoreDriver interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Remove(string)
	Keys(string) []string
	SetExpires(string, int64)
	GetExpires(string) (int64, bool)
}

//NewStore creates store with specific driver
func NewStore(driver StoreDriver) *Store {
	return &Store{
		AuthorizedTokens: []string{},
		Driver:           driver,
	}
}

//ValidToken return true if token authorized successfully
func (s *Store) ValidToken(token string) bool {
	s.RLock()
	defer s.RUnlock()
	for _, t := range s.AuthorizedTokens {
		if token == t {
			return true
		}
	}
	return false
}

//AddAuthorizedToken adds new token
func (s *Store) AddAuthorizedToken(token string) {
	s.Lock()
	s.AuthorizedTokens = append(s.AuthorizedTokens, token)
	s.Unlock()
}

//ValidValue returns true if value is string, slice of strings or map of strings by strings
func (s *Store) ValidValue(value interface{}) bool {
	switch x := value.(type) {
	case string:
		return true
	case []interface{}:
		for _, v := range x {
			if _, ok := v.(string); !ok {
				return false
			}
		}
		return true
	case map[string]interface{}:
		for _, v := range x {
			if _, ok := v.(string); !ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

//Set - set key with value
func (s *Store) Set(key string, value interface{}) error {
	if !s.ValidValue(value) {
		return fmt.Errorf("type of value must being string, []string or map[string]string")
	}
	s.Driver.Set(key, value)
	return nil
}

//Get - get value by key
func (s *Store) Get(key string) (interface{}, error) {
	if value, ok := s.Driver.Get(key); ok {
		return value, nil
	}
	return nil, fmt.Errorf("key %s not found", key)
}

//Remove - remove key
func (s *Store) Remove(key string) error {
	if _, ok := s.Driver.Get(key); ok {
		s.Driver.Remove(key)
		return nil
	}
	return fmt.Errorf("key %s not found", key)
}

//Keys - get keys by glob pattern
func (s *Store) Keys(pattern string) ([]string, error) {
	keys := s.Driver.Keys(pattern)
	if len(keys) > 0 {
		return keys, nil
	}
	return keys, fmt.Errorf("keys not found by pattern: %s", pattern)
}

//SetExpires set expiration time in seconds for key
func (s *Store) SetExpires(key string, expires int64) {
	s.Driver.SetExpires(key, expires)
}

//GetExpires returns expiration time in seconds for key
func (s *Store) GetExpires(key string) (int64, error) {
	if expires, ok := s.Driver.GetExpires(key); ok {
		return expires, nil
	}
	return 0, fmt.Errorf("expiration time for key %s is not set", key)
}
