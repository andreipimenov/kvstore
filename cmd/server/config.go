package main

import (
	"encoding/json"
	"errors"
)

type Config struct {
	SecretKey     string `json:"secretKey"`
	Authorization bool   `json:"authorization"`
	Users         []User `json:"users"`
	DumpFile      string `json:"dumpFile"`
	DumpInterval  int64  `json:"dumpInterval"`
	Port          int    `json:"port"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//ConfigDriver - interface for receiving configuration data from some source (file, other server etc)
type ConfigDriver interface {
	Get() (interface{}, error)
}

//NewConfig receives config data and create config structure
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
	cfg := &Config{}
	err = json.Unmarshal(j, cfg)
	if err != nil {
		return defaultCfg, err
	}
	return cfg, nil
}
