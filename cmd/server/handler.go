package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andreipimenov/kvstore/model"
	"github.com/go-chi/chi"
)

//WriteResponse - common helper function: marshal data to JSON and write into response
func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(j)
}

//WriteErrorResponse - helper function for errors: wrap errs in model.APIErrors, marshal to JSON and write into response
func WriteErrorResponse(w http.ResponseWriter, code int, errs ...*model.APIMessage) {
	j, _ := json.Marshal(&model.APIErrors{
		Errors: errs,
	})
	w.WriteHeader(code)
	w.Write(j)
}

//JSONCtx - setup all requests mime-type to application/json
func JSONCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//Authorization - middelware for checking tokens
func Authorization(s *Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.Split(r.Header.Get("Authorization"), " ")
			if len(auth) != 2 || auth[0] != "Token" {
				WriteErrorResponse(w, http.StatusUnauthorized, &model.APIMessage{
					Code: "Unauthorized", Message: "Unauthorized",
				})
				return
			}
			if !s.ValidToken(auth[1]) {
				WriteErrorResponse(w, http.StatusUnauthorized, &model.APIMessage{
					Code: "Unauthorized", Message: "Invalid token",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

//LoginHandler process authorization
func LoginHandler(c *Config, s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &model.APIAuth{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Cannot decode request body",
			})
			return
		}
		for _, user := range c.Users {
			if req.Login == user.Login && req.Password == user.Password {
				tokenData := fmt.Sprintf("%s%s%s", user.Login, c.SecretKey, user.Password)
				token := fmt.Sprintf("%x", sha256.Sum256([]byte(tokenData)))
				s.AddAuthorizedToken(token)
				WriteResponse(w, http.StatusOK, &model.APIAuth{
					Token: token,
				})
				return
			}
		}
		WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
			Code: "BadRequest", Message: "Invalid login and(or) password",
		})
	})
}

//NotAllowedHandler - handler for "Method Not Allowed" error
func NotAllowedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, &model.APIMessage{
			Code: "NotAllowed", Message: "Method Not Allowed",
		})
	})
}

//NotFoundHandler - handler for "Not Found" error
func NotFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
			Code: "NotFound", Message: "Invalid API endpoint",
		})
	})
}

//PingHandler - health check handler
func PingHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusOK, &model.APIMessage{
			Message: "pong",
		})
	})
}

//SetHandler - set value with key (add or replace)
func SetHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &model.APIKeyValue{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Cannot decode request body",
			})
			return
		}
		if req.Key == "" {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Key must being not-empty string",
			})
			return
		}
		err = s.Set(req.Key, req.Value)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Invalid value",
			})
			return
		}
		WriteResponse(w, http.StatusCreated, &model.APIMessage{
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
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key),
			})
			return
		}
		WriteResponse(w, http.StatusOK, &model.APIKeyValue{
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
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key),
			})
			return
		}
		WriteResponse(w, http.StatusOK, &model.APIMessage{
			Message: "OK",
		})
	})
}

//KeysHandler - get keys by pattern
func KeysHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := chi.URLParam(r, "pattern")
		keys, _ := s.Keys(pattern)
		WriteResponse(w, http.StatusOK, &model.APIKeys{
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
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key),
			})
			return
		}
		switch v := value.(type) {
		case []interface{}:
			i, err := strconv.Atoi(index)
			if err != nil {
				WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
					Code: "BadRequest", Message: "Index type must being int for this key",
				})
				return
			}
			if i < 0 || i > len(v)-1 {
				WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
					Code: "BadRequest", Message: "Index out of range",
				})
				return
			}
			WriteResponse(w, http.StatusOK, &model.APIKeyValue{
				Value: v[i],
			})
			return
		case map[string]interface{}:
			if _, ok := v[index]; ok {
				WriteResponse(w, http.StatusOK, &model.APIKeyValue{
					Value: v[index],
				})
				return
			}
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: fmt.Sprintf("Index %s is not set", index),
			})
			return
		default:
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Value must being []string or map[string]string",
			})
			return
		}
	})
}

//SetExpiresHandler - set expiration time for key
func SetExpiresHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		_, err := s.Get(key)
		if err != nil {
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key),
			})
			return
		}
		req := &model.APIKeyExpires{}
		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Cannot decode request body",
			})
			return
		}
		if req.Expires <= 0 {
			WriteErrorResponse(w, http.StatusBadRequest, &model.APIMessage{
				Code: "BadRequest", Message: "Expiration time must being positive int64 number",
			})
			return
		}
		s.SetExpires(key, req.Expires)
		WriteResponse(w, http.StatusOK, &model.APIMessage{
			Message: "OK",
		})
	})
}

//GetExpiresHandler - get expiration time for key
func GetExpiresHandler(s *Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		_, err := s.Get(key)
		if err != nil {
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Key %s not found", key),
			})
			return
		}
		expires, err := s.GetExpires(key)
		if err != nil {
			WriteErrorResponse(w, http.StatusNotFound, &model.APIMessage{
				Code: "NotFound", Message: fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
		WriteResponse(w, http.StatusOK, &model.APIKeyExpires{
			Expires: expires,
		})
	})
}
