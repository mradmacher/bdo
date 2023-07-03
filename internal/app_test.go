package bdo

import (
  "testing"
  "net/http/httptest"
  "io"
	"github.com/joho/godotenv"
)

func TestSearch(t *testing.T) {
	if err := godotenv.Load("../.test.env"); err != nil {
		panic("No .test.env file found")
	}

  app := NewApp()
  defer app.Close()

  ts := httptest.NewServer(app)
  defer ts.Close()

  res, err := ts.Client().Get(ts.URL + "/api/installations")
  if err != nil {
    t.Fatal(err)
  }
  json, err := io.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    t.Fatal(err)
  }
  if string(json) != "[]" {
    t.Errorf("/api/installations => []; got: %v", string(json))
  }
}
