package bdo

import (
  "testing"
  "net/http"
  "net/http/httptest"
	"github.com/joho/godotenv"
)

func TestSearch(t *testing.T) {
	if err := godotenv.Load("../.test.env"); err != nil {
		panic("No .test.env file found")
	}

  app := NewApp()
  app.MountHandlers()
  defer app.Stop()

  req := httptest.NewRequest("GET", "/api/installations", nil)

  res := httptest.NewRecorder()
  app.ServeHTTP(res, req)

  if res.Code != http.StatusOK {
    t.Errorf("Expected response code %d; got %d\n", http.StatusOK, res.Code)
  }

  if res.Body.String() != "[]" {
    t.Errorf("/api/installations => []; got: %v", res.Body.String())
  }
}
