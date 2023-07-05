package bdo

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
  "net/http"
)

type App struct {
  *chi.Mux
  db DbClient
}

func NewApp() *App {
  app := App{chi.NewRouter(), DbClient{}}

	err := app.db.Connect()
	if err != nil {
		panic(err)
	}

  return &app
}

func (app *App) MountHandlers() {
	app.Get("/", homeHandler)
	app.Get("/assets/*", staticHandler)
	app.Get("/api/installations", app.searchInstallationsHandler)
}

func (app *App) Start() {
	http.ListenAndServe(":3000", app)
}

func (app *App) Stop() {
  if err := app.db.Disconnect(); err != nil {
    panic(err)
  }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "views/index.html")
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
