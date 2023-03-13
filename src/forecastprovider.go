package main

type ForecastProvider interface {
	GetForecast() (Forecast, error)
	SetConfig(c *Config)
}
