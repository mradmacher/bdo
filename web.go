package main

import (
  "fmt"
  "net/http"
  "github.com/joho/godotenv"
  "encoding/json"
  "github.com/go-chi/chi/v5"
  "github.com/mradmacher/bdo/internal"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  http.ServeFile(w, r, "views/index.html")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "." + r.URL.Path)
}

func searchInstallationsHandler(w http.ResponseWriter, r *http.Request) {
  db := bdo.DbClient{}
  err := db.Connect()
  if err != nil { panic(err) }

  defer func() {
      if err := db.Disconnect(); err != nil {
          panic(err)
      }
  }()

  params := bdo.SearchParams{}
  if r.FormValue("wc") != "" {
    params["waste_code"] = r.FormValue("wc")
  }
  if r.FormValue("pc") != "" {
    params["process_code"] = r.FormValue("pc")
  }
  if r.FormValue("sc") != "" {
    params["state_code"] = r.FormValue("sc")
  }
  repo := db.NewInstallationRepo()
  results, err := repo.Search(params)
  if err != nil { panic(err) }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  if len(results) != 0 {
    jsonBlob, err := json.Marshal(&results)
    if err != nil { panic(err) }

    fmt.Fprint(w, string(jsonBlob[:]))
  } else {
    fmt.Fprint(w, "[]")
  }
}

func main() {
  if err := godotenv.Load(); err != nil {
      panic("No .env file found")
  }


  router := chi.NewRouter()
  router.Get("/", homeHandler)
  router.Get("/assets/*", staticHandler)
  router.Get("/api/installations", searchInstallationsHandler)

  fmt.Println("Starting the server on :3000...")
  http.ListenAndServe(":3000", router)
}
