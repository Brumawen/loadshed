package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ConfigController handles the Web Methods for configuring the module.
type ConfigController struct {
	Srv *Server
}

// ConfigPageData holds the data used to write to the configuration page.
type ConfigPageData struct {
	Provider int    // Forecast Provider: 0=EskomSePush
	AreaID   string // Provider Application Identifier
	Token    string // Authentication token
	Timeout  int    // Forecast timeout
}

// AddController adds the controller routes to the router
func (c *ConfigController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Path("/config.html").Handler(http.HandlerFunc(c.handleConfigWebPage))
	router.Methods("GET").Path("/config/get").Name("GetConfig").
		Handler(Logger(c, http.HandlerFunc(c.handleGetConfig)))
	router.Methods("POST").Path("/config/set").Name("SetConfig").
		Handler(Logger(c, http.HandlerFunc(c.handleSetConfig)))
}

func (c *ConfigController) handleConfigWebPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./html/config.html"))

	v := ConfigPageData{
		Provider: c.Srv.Config.Provider,
		AreaID:   c.Srv.Config.Id,
		Token:    c.Srv.Config.Token,
		Timeout:  c.Srv.Config.ForecastTimeoutMins,
	}

	err := t.Execute(w, v)
	if err != nil {
		c.LogError("Failed to execute template. ", err.Error())
	}
}

func (c *ConfigController) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	if err := c.Srv.Config.WriteTo(w); err != nil {
		http.Error(w, "Error serializing configuration. "+err.Error(), 500)
	}
}

func (c *ConfigController) handleSetConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	prv := r.Form.Get("provider")
	aid := r.Form.Get("areaid")
	tkn := r.Form.Get("token")
	tmo := r.Form.Get("timeout")

	if prv == "" {
		http.Error(w, "The Forecast Provider must be selected", 500)
		return
	}
	p, err := strconv.Atoi(prv)
	if err != nil || p < 0 || p > 1 {
		http.Error(w, "Invalid Forecast Provider value", 500)
		return
	}
	if aid == "" {
		http.Error(w, "Invalid Area ID value", 500)
		return
	}
	if tkn == "" {
		http.Error(w, "Invalid Token value", 500)
		return
	}
	if tmo == "" {
		http.Error(w, "Invalid Timeout value", 500)
		return
	}
	tout, err := strconv.Atoi(tmo)
	if err != nil || tout < 0 {
		http.Error(w, "Invalid Timeout value", 500)
		return
	}

	c.LogInfo("Setting new configuration values.")

	c.Srv.Config.Provider = p
	c.Srv.Config.Id = aid
	c.Srv.Config.Token = tkn
	c.Srv.Config.ForecastTimeoutMins = tout

	c.Srv.Config.SetDefaults()

	c.Srv.Config.WriteToFile("config.json")
}

// LogInfo is used to log information messages for this controller.
func (c *ConfigController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("ConfigController: [Inf] ", a)
}

// LogError is used to log error messages for this controller.
func (c *ConfigController) LogError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("ConfigController: [Err] ", a)
}
