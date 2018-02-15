package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

//JSONCtx - setup all requests mime-type to application/json
func JSONCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//NotAllowedHandler - handler for "Method Not Allowed" error
func NotAllowedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusMethodNotAllowed, APIErrors{
			[]APIMessage{
				{Code: "NotAllowed", Message: "Method Not Allowed"},
			},
		})
	})
}

//NotFoundHandler - handler for "Not Found" error
func NotFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusNotFound, APIErrors{
			[]APIMessage{
				{Code: "NotFound", Message: "Not Found"},
			},
		})
	})
}

//PingHandler - health check handler
func PingHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusOK, APIMessage{
			Message: "pong",
		})
	})
}

//SetHandler - set value with key (add or replace)
func SetHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &APIKeyValue{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, APIErrors{
				[]APIMessage{
					{Code: "BadRequest", Message: "Cannot decode request body"},
				},
			})
			return
		}
		err = s.Set(req.Key, req.Value)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, APIErrors{
				[]APIMessage{
					{Code: "BadRequest", Message: "Invalid value"},
				},
			})
			return
		}
		WriteResponse(w, http.StatusCreated, APIMessage{
			Message: "OK",
		})
	})
}

//GetHandler - get value by key
func GetHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		value, err := s.Get(key)
		if err != nil {
			WriteResponse(w, http.StatusNotFound, APIErrors{
				[]APIMessage{
					{Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key)},
				},
			})
			return
		}
		WriteResponse(w, http.StatusOK, APIKeyValue{
			Value: value,
		})
	})
}

//RemoveHandler - remove key
func RemoveHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		err := s.Remove(key)
		if err != nil {
			WriteResponse(w, http.StatusNotFound, APIErrors{
				[]APIMessage{
					{Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key)},
				},
			})
			return
		}
		WriteResponse(w, http.StatusOK, APIMessage{
			Message: "OK",
		})
	})
}

//KeysHandler - get keys by pattern
func KeysHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := chi.URLParam(r, "pattern")
		keys, _ := s.Keys(pattern)
		WriteResponse(w, http.StatusOK, APIKeys{
			Keys: keys,
		})
	})
}

//GetIndexHandler - get value by key
func GetIndexHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		index := chi.URLParam(r, "index")
		value, err := s.Get(key)
		if err != nil {
			WriteResponse(w, http.StatusNotFound, APIErrors{
				[]APIMessage{
					{Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key)},
				},
			})
			return
		}
		switch v := value.(type) {
		case []interface{}:
			i, err := strconv.Atoi(index)
			if err != nil {
				WriteResponse(w, http.StatusBadRequest, APIErrors{
					[]APIMessage{
						{Code: "BadRequest", Message: "Index type must being int for this key"},
					},
				})
				return
			}
			if i < 0 || i > len(v)-1 {
				WriteResponse(w, http.StatusBadRequest, APIErrors{
					[]APIMessage{
						{Code: "BadRequest", Message: "Index out of range"},
					},
				})
				return
			}
			WriteResponse(w, http.StatusOK, APIKeyValue{
				Value: v[i],
			})
			return
		case map[string]interface{}:
			if _, ok := v[index]; ok {
				WriteResponse(w, http.StatusOK, APIKeyValue{
					Value: v[index],
				})
				return
			}
			WriteResponse(w, http.StatusBadRequest, APIErrors{
				[]APIMessage{
					{Code: "BadRequest", Message: fmt.Sprintf("Index %s is not set", index)},
				},
			})
			return
		default:
			WriteResponse(w, http.StatusBadRequest, APIErrors{
				[]APIMessage{
					{Code: "BadRequest", Message: "Value must being []string or map[string]string"},
				},
			})
			return
		}
	})
}
