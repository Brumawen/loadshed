package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Provider            int    `json:"provider"`        // Provider type: 0=EskomSePush
	Token               string `json:"token"`           // API Token
	Id                  string `json:"id"`              // Area identifier
	ForecastTimeoutMins int    `json:"forecastTimeout"` // Timeout of current forecast in minutes
}

// ReadFromFile will read the configuration settings from the specified file
func (c *Config) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteToFile will write the configuration settings to the specified file
func (c *Config) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// WriteTo serializes the entity and writes it to the http response
func (c *Config) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (c *Config) Serialize() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (c *Config) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &c)
	c.SetDefaults()
	return err
}

// SetDefaults checks the configuration and makes sure that, if a value is not configured, the default value is set.
func (c *Config) SetDefaults() {
	// Set any defaults required
	if c.ForecastTimeoutMins <= 0 {
		c.ForecastTimeoutMins = 60
	}
}
