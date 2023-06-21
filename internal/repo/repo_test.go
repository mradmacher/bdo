package repo

import (
    "testing"
    "github.com/joho/godotenv"
)

func setupSuite(t *testing.T) (func(*testing.T), *DbClient) {
    if err := godotenv.Load("../../.env"); err != nil {
        t.Fatalf("No .env file found")
    }

    db := DbClient{}
    err := db.Connect()
    if err != nil {
        t.Fatalf("Problem with DB connection: %v", err)
    }
    return func(t *testing.T) {
        db.Disconnect()
    }, &db
}

func TestSearchReturnsAllDocuments(t *testing.T) {
    teardownSuite, db := setupSuite(t)
    defer teardownSuite(t)

    r := db.NewInstallationRepo()
    r.Purge()
    r.Add(&Installation{Name: "Test"})
    results, _ := r.Search()
    want := 1
    got := len(results)
    if want != got {
        t.Fatalf("Expected %d, got %d", want, got)
    }
}
