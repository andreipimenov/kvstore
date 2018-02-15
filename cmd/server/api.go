package main

import (
	"encoding/json"
	"net/http"
)

//APIErrors contains all errors responsed by server
type APIErrors struct {
	Errors []APIMessage `json:"errors"`
}

//APIMessage - common server response with code and message
type APIMessage struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//APIKeyValue - common server request/response with key and(or) value
type APIKeyValue struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

//APIKeys - server response for multiple APIKeys
type APIKeys struct {
	Keys []string `json:"keys"`
}

//WriteResponse - marshal to JSON and write data into response
func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(j)
}
