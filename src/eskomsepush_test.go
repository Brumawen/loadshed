package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestLoadResponse(t *testing.T) {
	r := EspResponse{}

	err := r.ReadFromFile("./test/EspResponse.json")
	if err != nil {
		t.Error(err)
	}
	if r.Info.Name != "Dainfern College (10)" {
		t.Error("Info.Name wrong, ", r.Info.Name)
	}
	if r.Info.Region != "JHB City Power" {
		t.Error("Info.Region wrong, ", r.Info.Region)
	}
	if len(r.Events) != 8 {
		t.Error("Wrong events count, ", len(r.Events))
	}

	evt := r.Events[0]
	s, err := evt.GetStage()
	if err != nil {
		t.Error(err)
	}
	if s != 4 {
		t.Error("Wrong event stage, ", s)
	}
	d := evt.GetDisplay()
	if d != "12:00-14:30" {
		t.Error("Wrong display, ", d)
	}
}

func TestGetForecastFromUrl(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("eskomsepush_test.json")
	if err != nil {
		t.Error(err)
	}

	p := EskomSePush{Config: &c}

	f, err := p.GetForecastFromUrl()
	if err != nil {
		t.Error(err)
	}
	data, err := json.Marshal(&f)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("test.json", data, 0666)
	if err != nil {
		t.Error(err)
	}
}

func TestGetForecast(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("eskomsepush_test.json")
	if err != nil {
		t.Error(err)
	}

	p := EskomSePush{Config: &c}

	f, err := p.GetForecast()
	if err != nil {
		t.Error(err)
	}
	data, err := json.Marshal(&f)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("test.json", data, 0666)
	if err != nil {
		t.Error(err)
	}
}
