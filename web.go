package main

import (
  "fmt"
  "net/http"
  "github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  http.ServeFile(w, r, "views/index.html")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "." + r.URL.Path)
}

func main() {
  router := chi.NewRouter()
  router.Get("/", homeHandler)
  router.Get("/assets/*", staticHandler)

  fmt.Println("Starting the server on :3000...")
  http.ListenAndServe(":3000", router)
}
