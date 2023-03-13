package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Forecast struct {
	Name    string    `json:"name"`    // Area Name
	Region  string    `json:"region"`  // Region Name
	Stage   int       `json:"stage"`   // Current Stage
	Events  []Event   `json:"events"`  // Event information
	Created time.Time `json:"created"` // Date and time forecast was created
}

type Event struct {
	Start   time.Time `json:"start"` // Start date and time
	End     time.Time `json:"end"`   // End date and time
	Day     string    `json:"day"`   // The day of the event (e.g. Mon, Tue etc)
	Display string    `json:"note"`  // Display information
	Stage   int       `json:"stage"` // Stage
}

// ReadFromFile will read the forecast information from the specified file
func (c *Forecast) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	return err
}

// WriteToFile will write the forecast information to the specified file
func (c *Forecast) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// WriteTo serializes the entity and writes it to the http response
func (c *Forecast) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}
