package main

import (
  "fmt"
  "net/http"
)

type WebApp struct{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, "<h1>Welcome</h1>")
}

func (app WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.URL.Path {
    case "/":
      homeHandler(w, r)
    default:
      http.Error(w, "Page not found", http.StatusNotFound)
  }
}


func main() {
  var webApp WebApp

  fmt.Println("Starting the server on :3000...")
  http.ListenAndServe(":3000", webApp)
}
