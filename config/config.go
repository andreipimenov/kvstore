package config

import (
	"io/ioutil"
	"os"
)

//Config - driver which reads config file
type Config struct {
	File string
}

//New creates new Config
func New(file string) *Config {
	return &Config{
		File: file,
	}
}

//Get - reads file and return raw data
func (c *Config) Get() (interface{}, error) {
	_, err := os.Stat(c.File)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadFile(c.File)
	if err != nil {
		return nil, err
	}
	return f, nil
}
