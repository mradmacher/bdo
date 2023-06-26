package main

import (
  "fmt"
  "net/http"
  "github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, "<h1>Welcome</h1>")
}

func main() {
  router := chi.NewRouter()
  router.Get("/", homeHandler)
  router.NotFound(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Page not found", http.StatusNotFound)
  })

  fmt.Println("Starting the server on :3000...")
  http.ListenAndServe(":3000", router)
}
