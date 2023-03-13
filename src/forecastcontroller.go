package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ForecastController handles the Web Methods for retrieving loadshedding forecast information.
type ForecastController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *ForecastController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/forecast/get").Name("GetCurrentForecast").
		Handler(Logger(c, http.HandlerFunc(c.handleGetCurrent)))
}

// Get the current weather information
func (c *ForecastController) handleGetCurrent(w http.ResponseWriter, r *http.Request) {
	if p, err := c.getForecastProvider(); err != nil {
		c.LogError("Error getting forecast provider. " + err.Error())
		http.Error(w, "Error getting forecast provider. "+err.Error(), 500)
	} else {
		if f, err := p.GetForecast(); err != nil {
			c.LogError("Error getting forecast. " + err.Error())
			http.Error(w, "Error getting forecast. "+err.Error(), 500)
		} else {
			if err := f.WriteTo(w); err != nil {
				c.LogError("Error serializing forecast information. " + err.Error())
				http.Error(w, "Error serializing forecast information. "+err.Error(), 500)
			}
		}
	}
}

func (c *ForecastController) getForecastProvider() (ForecastProvider, error) {
	switch c.Srv.Config.Provider {
	case 0:
		// EskomSePush
		esp := new(EskomSePush)
		esp.SetConfig(c.Srv.Config)
		return esp, nil
	default:
		return nil, errors.New("Invalid Forecast provider")
	}
}

// LogInfo is used to log information messages for this controller.
func (c *ForecastController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("ForecastController: ", a)
}

// LogError is used to log error messages for this controller.
func (c *ForecastController) LogError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("ForecastController: ", a)
}
