package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/gobwas/glob"
)

//Store implements cache: in-memory store with ttl and http-api
type Store struct {
	sync.RWMutex `json:"-"`
	Data         map[string]interface{} `json:"data"`
	Expires      map[string]int64       `json:"expires"`
	DumpFile     string                 `json:"-"`
	DumpInterval int64                  `json:"-"`
}

//New creates Store, runs workers for removing expired keys and autosave storage into file
func New(dumpFile string, dumpInterval int64) *Store {
	s := &Store{
		Data:         map[string]interface{}{},
		Expires:      map[string]int64{},
		DumpFile:     dumpFile,
		DumpInterval: dumpInterval,
	}
	if dumpInterval > 0 {
		fileData, err := ioutil.ReadFile(dumpFile)
		if err != nil {
			log.Printf("Error loading storage dump: %s\n", err.Error())
		} else {
			json.Unmarshal(fileData, s)
		}
		go s.dumpWorker()
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
				if currentTime >= value {
					s.Remove(key)
				}
			}
		}()
	}
}

//dumpWorker saves store to file
func (s *Store) dumpWorker() {
	for {
		<-time.After(time.Duration(s.DumpInterval) * time.Second)
		j, _ := json.Marshal(s)
		ioutil.WriteFile(s.DumpFile, j, 0644)
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
		return value - time.Now().Unix(), true
	}
	return 0, false
}
