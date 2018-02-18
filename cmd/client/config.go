package main

import (
	"encoding/json"
	"errors"
)

//Config - application-specific configurations
type Config struct {
	ServerHost string `json:"serverHost"`
	ServerPort int    `json:"serverPort"`
	Port       int    `json:"port"`
}

//ConfigDriver - interface for receiving configuration from some source (file, web source etc)
type ConfigDriver interface {
	Get() (interface{}, error)
}

//NewConfig returns server configuration taken from driver
func NewConfig(driver ConfigDriver) (*Config, error) {
	defaultCfg := &Config{
		Port: 8090,
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
