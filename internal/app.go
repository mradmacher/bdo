package bdo

import (
	"net/http"
	"os"
)

type App struct {
	router  *http.ServeMux
	db       DbClient
	renderer Renderer
}

func NewApp(renderer Renderer) (*App, error) {
	app := App{}
	app.router = http.NewServeMux()
	app.db = DbClient{}
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
	app.router.HandleFunc("GET /api/installations", app.searchInstallationsHandler)
	app.router.HandleFunc("GET /api/installation/{id}", app.showInstallationHandler)
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
}

func (app *App) showInstallationHandler(w http.ResponseWriter, r *http.Request) {
	repo := app.db.NewInstallationRepo()
	installation, err := repo.Find(r.PathValue("id"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = app.renderer.RenderInstallationSummary(w, *installation)
	if err != nil {
		panic(err)
	}
}

func (app *App) searchInstallationsHandler(w http.ResponseWriter, r *http.Request) {
	params := SearchParams{}
	Bind(params, r)
	repo := app.db.NewInstallationRepo()
	installations, err := repo.Search(params)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = app.renderer.RenderInstallations(w, installations)
	if err != nil {
		panic(err)
	}
}
