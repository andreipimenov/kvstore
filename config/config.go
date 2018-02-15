package config

import (
	"io/ioutil"
	"os"
)

//Config implements configuration from file
type Config struct {
	File string
}

//New creates new Config struct
func New(file string) *Config {
	return &Config{
		File: file,
	}
}

//Get reads file and return its data
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
