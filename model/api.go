package model

//APIAuth contains login, password and token for authenticated access/interact with storage
type APIAuth struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

//APIErrors contains all errors responsed by server
type APIErrors struct {
	Errors []*APIMessage `json:"errors"`
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

//APIKeyExpires - struct for request/response expiration time for specific key
type APIKeyExpires struct {
	Expires int64 `json:"expires"`
}
