package main

import (
	"encoding/json"
	"errors"
)

//Config - application-specific configurations
type Config struct {
	SecretKey     string `json:"secretKey"`
	Authorization bool   `json:"authorization"`
	Users         []User `json:"users"`
	DumpFile      string `json:"dumpFile"`
	DumpInterval  int64  `json:"dumpInterval"`
	Port          int    `json:"port"`
}

//User - part of configuration for user auth data
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//ConfigDriver - interface for receiving configuration from some source (file, web source etc)
type ConfigDriver interface {
	Get() (interface{}, error)
}

//NewConfig returns server configuration taken from driver
func NewConfig(driver ConfigDriver) (*Config, error) {
	defaultCfg := &Config{
		Port: 8080,
	}
	v, err := driver.Get()
	if err != nil {
		return defaultCfg, err
	}
	j, ok := v.([]byte)
	if !ok {
		return defaultCfg, errors.New("error assertion interface{} to []byte")
	}
	c := &Config{}
	err = json.Unmarshal(j, c)
	if err != nil {
		return defaultCfg, err
	}
	return c, nil
}
