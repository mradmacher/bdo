package bdo

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"os"
)

type App struct {
	router   *chi.Mux
	db       DbClient
	template *template.Template
}

func NewApp(templatesPath string) (*App, error) {
	app := App{}
	app.router = chi.NewRouter()
	app.db = DbClient{}

	var err error

	app.template, err = template.ParseFiles(templatesPath + "/index.html")
	if err != nil {
		return nil, err
	}

	err = app.db.Connect()
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (app *App) MountHandlers() {
	app.router.Get("/", app.homeHandler)
	app.router.Get("/assets/*", staticHandler)
	app.router.Get("/api/installations", app.searchInstallationsHandler)
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
	err := app.template.Execute(w, config)
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
}

func (app *App) searchInstallationsHandler(w http.ResponseWriter, r *http.Request) {
	params := SearchParams{}
	Bind(params, r)
	repo := app.db.NewInstallationRepo()
	results, err := repo.Search(params)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(results) != 0 {
		jsonBlob, err := json.Marshal(&results)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(w, string(jsonBlob[:]))
	} else {
		fmt.Fprint(w, "[]")
	}
}
