package bdo

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

type App struct {
	router   *http.ServeMux
	db       Repository
	renderer Renderer
}

func NewApp(renderer Renderer) (*App, error) {
	app := App{}
	app.db = Repository{}
	app.router = http.NewServeMux()
	app.renderer = renderer

	var err error
	err = app.db.Connect()
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (app *App) MountHandlers() {
	app.router.HandleFunc("GET /assets/main.js", staticHandler)
	app.router.HandleFunc("GET /{$}", app.homeHandler)
	app.router.HandleFunc("GET /instalacje", app.searchInstallationsHandler)
	app.router.HandleFunc("GET /instalacje/{id}/mozliwosci", app.searchInstallationCapabilitiesHandler)
	app.router.HandleFunc("GET /mozliwosci", app.searchCapabilitiesHandler)
}

func (app *App) Start() {
	http.ListenAndServe(":3000", app.router)
}

func (app *App) Stop() {
	if err := app.db.Disconnect(); err != nil {
		panic(err)
	}
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	config := struct {
		GoogleMapsApiKey string
	}{os.Getenv("GOOGLE_MAPS_API_KEY")}
	err := app.renderer.RenderHome(w, config)
	if err != nil {
		panic(err)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

func Bind(sp SearchParams, r *http.Request) {
	if r.FormValue("wc") != "" {
		sp["waste_code"] = r.FormValue("wc")
	}
	if r.FormValue("pc") != "" {
		sp["process_code"] = r.FormValue("pc")
	}
	if r.FormValue("sc") != "" {
		sp["state_code"] = r.FormValue("sc")
	}
	if r.FormValue("id") != "" {
		sp["installation_id"] = r.FormValue("id")
	}
}

func redirectNoAjax(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept"), "application/json") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *App) searchInstallationCapabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	redirectNoAjax(w, r)

	var capabilities []Capability
	var err error
	params := SearchParams{}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err == nil {
		Bind(params, r)
		capabilities, err = app.db.SearchCapabilities(id, params)
		if err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = app.renderer.RenderCapabilities(w, capabilities)
	if err != nil {
		panic(err)
	}
}

func (app *App) searchCapabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	redirectNoAjax(w, r)

	var capabilities []Capability
	var err error
	params := SearchParams{}
	Bind(params, r)
	capabilities, err = app.db.Summarize(params)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = app.renderer.RenderCapabilitiesSummary(w, capabilities)
	if err != nil {
		panic(err)
	}
}

func (app *App) searchInstallationsHandler(w http.ResponseWriter, r *http.Request) {
	redirectNoAjax(w, r)

	params := SearchParams{}
	Bind(params, r)
	var installations []*Installation
	var err error
	installations, err = app.db.Search(params)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = app.renderer.RenderInstallations(w, installations)
	if err != nil {
		panic(err)
	}
}
