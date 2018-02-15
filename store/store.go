package store

import (
	"sync"
	"time"

	"github.com/gobwas/glob"
)

//Store implements cache: in-memory store with ttl and http-api
type Store struct {
	sync.RWMutex
	Data    map[string]interface{}
	Expires map[string]int64
}

//New creates Store and runs worker for removing expired keys
func New() *Store {
	s := &Store{
		Data:    map[string]interface{}{},
		Expires: map[string]int64{},
	}
	go s.expiresWorker()
	return s
}

//expiresWorker removes expired keys
func (s *Store) expiresWorker() {
	for {
		<-time.After(1 * time.Second)
		go func() {
			currentTime := time.Now().Unix()
			for key, value := range s.Expires {
				if currentTime > value {
					s.Remove(key)
				}
			}
		}()
	}
}

//Set value associated with key
func (s *Store) Set(key string, value interface{}) {
	s.Lock()
	s.Data[key] = value
	s.Unlock()
}

//Get value by key
func (s *Store) Get(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	if value, ok := s.Data[key]; ok {
		return value, true
	}
	return nil, false
}

//Remove key
func (s *Store) Remove(key string) {
	s.Lock()
	delete(s.Data, key)
	delete(s.Expires, key)
	s.Unlock()
}

//Keys returns all keys by glob pattern
func (s *Store) Keys(pattern string) []string {
	e := glob.MustCompile(pattern)
	keys := []string{}
	for key := range s.Data {
		if e.Match(key) {
			keys = append(keys, key)
		}
	}
	return keys
}

//SetExpires setup expires time for key in seconds
func (s *Store) SetExpires(key string, expires int64) {
	s.Lock()
	s.Expires[key] = time.Now().Unix() + expires
	s.Unlock()
}

//GetExpires returns seconds until key expires
func (s *Store) GetExpires(key string) (int64, bool) {
	s.RLock()
	defer s.RUnlock()
	if value, ok := s.Expires[key]; ok {
		return time.Now().Unix() - value, true
	}
	return 0, false
}
