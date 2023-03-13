package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// EskomSePush is an interface to the EskomSePush internet API
type EskomSePush struct {
	Config *Config // Current Configuration
}

// SetConfig sets the configuration for the provider
func (p *EskomSePush) SetConfig(c *Config) {
	p.Config = c
}

// GetForecast returns the current forecast (within last hour)
func (p *EskomSePush) GetForecast() (Forecast, error) {
	p.logDebug("Checking saved forecast file. Timeout is ", p.Config.ForecastTimeoutMins, " mins")
	f := Forecast{}
	err := f.ReadFromFile("last_esp_forecast.json")
	if err == nil {
		// Check the date
		timeout := f.Created.Add(time.Duration(p.Config.ForecastTimeoutMins) * time.Minute)
		if timeout.After(time.Now()) {
			p.logDebug("Using saved forecast file.  Forecast last read ", f.Created)
			return f, nil
		}
	}

	uf, err := p.GetForecastFromUrl()
	if err != nil {
		return f, err
	}

	p.logDebug("Saving forecast file")
	err = uf.WriteToFile("last_esp_forecast.json")
	if err != nil {
		p.logError("Failed to save last forecast", err.Error())
	}

	return uf, nil
}

// GetForecastFromUrl returns the forecast from EskomSePush API
func (p *EskomSePush) GetForecastFromUrl() (Forecast, error) {
	f := Forecast{}
	if p.Config.Id == "" {
		return f, errors.New("Invalid or missing Id in the configuration")
	}
	if p.Config.Token == "" {
		return f, errors.New("Invalid or missing Token in the configuration")
	}

	p.logDebug("Getting forecast from EskomSePush API")

	url := fmt.Sprintf("https://developer.sepush.co.za/business/2.0/area?id=%s", p.Config.Id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return f, err
	}
	req.Header.Set("Token", p.Config.Token)

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
		resp.Close = true
	}
	if err == nil {
		er := EspResponse{}
		b, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(b, &er)
			if err == nil {
				f.Name = er.Info.Name
				f.Region = er.Info.Region
				f.Created = time.Now()
				for idx, erEvt := range er.Events {
					stage, _ := erEvt.GetStage()
					if idx == 0 {
						f.Stage = stage
					}
					evt := Event{}
					evt.Start = erEvt.Start
					evt.End = erEvt.End
					evt.Display = erEvt.GetDisplay()
					evt.Day = erEvt.Start.Weekday().String()
					evt.Stage = stage
					f.Events = append(f.Events, evt)
				}
			}
		}
	}

	return f, err
}

// logDebug logs a debug message to the logger
func (p *EskomSePush) logDebug(v ...interface{}) {
	a := fmt.Sprint(v...)
	if logger == nil {
		log.Println("EskomSePush: [Dbg] ", a)
	} else {
		logger.Info("EskomSePush: [Dbg] ", a)
	}
}

// logError logs an error message to the logger
func (p *EskomSePush) logError(v ...interface{}) {
	a := fmt.Sprint(v...)
	if logger == nil {
		log.Println("EskomSePush: [Err] ", a)
	} else {
		logger.Error("EskomSePush: [Err] ", a)
	}
}
