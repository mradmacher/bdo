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

	renderer, err := NewRenderer()
	if err != nil {
		t.Errorf("Error creating the renderer: %v", err)
	}

	app, err := NewApp(*renderer)
	if err != nil {
		t.Errorf("Error creating the app: %v", err)
	}
	app.MountHandlers()
	defer app.Stop()

	req := httptest.NewRequest("GET", "/instalacje", nil)

	res := httptest.NewRecorder()
	app.router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected response code %d; got %d\n", http.StatusOK, res.Code)
	}

	if res.Body.String() != "\nBrak instalacji spełniających podane kryteria\n\n\n" {
		t.Errorf("/instalacje => %q, got: %q", "", res.Body.String())
	}
}
