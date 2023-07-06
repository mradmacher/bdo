package bdo

import (
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	if err := godotenv.Load("../.test.env"); err != nil {
		panic("No .test.env file found")
	}

	app, err := NewApp("../views")
	if err != nil {
		t.Errorf("Error creating the app: %v", err)
	}
	app.MountHandlers()
	defer app.Stop()

	req := httptest.NewRequest("GET", "/api/installations", nil)

	res := httptest.NewRecorder()
	app.router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected response code %d; got %d\n", http.StatusOK, res.Code)
	}

	if res.Body.String() != "[]" {
		t.Errorf("/api/installations => []; got: %v", res.Body.String())
	}
}
