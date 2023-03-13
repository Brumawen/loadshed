package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// EspResponse holds the response from the EskomSePush API call
type EspResponse struct {
	Events []EspEvent `json:"events"`
	Info   struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}
}

// EspEvent holds information about a forecasted loadshedding event
type EspEvent struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Note  string    `json:"note"`
}

// ReadFromFile will read the response from the specified file
func (r *EspResponse) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &r)
		}
	}
	return err
}

func (e *EspEvent) GetStage() (int, error) {
	if strings.HasPrefix(e.Note, "Stage ") {
		return strconv.Atoi(e.Note[6:])
	}
	return -1, errors.New("Stage not specified")
}

func (e *EspEvent) GetDisplay() string {
	return fmt.Sprintf("%02d:%02d-%02d:%02d", e.Start.Hour(), e.Start.Minute(), e.End.Hour(), e.End.Minute())
}
