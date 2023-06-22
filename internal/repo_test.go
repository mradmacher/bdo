package bdo

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
    r.Add(
        &Installation{
            Name: "Test1",
            Capabilities: []Capability{
                Capability{
                    WasteCode: "020202",
                    ProcessCode: "R2",
                    Quantity: 123,
                },
                Capability{
                    WasteCode: "010101",
                    ProcessCode: "R2",
                    Quantity: 123,
                },
            },
        },
    )
    r.Add(
        &Installation{
            Name: "Test2",
            Capabilities: []Capability{
                Capability{
                    WasteCode: "010101",
                    ProcessCode: "R1",
                    Quantity: 123,
                },
                Capability{
                    WasteCode: "020202",
                    ProcessCode: "R2",
                    Quantity: 123,
                },
            },
        },
    )
    r.Add(
        &Installation{
            Name: "Test3",
            Capabilities: []Capability{
                Capability{
                    WasteCode: "030303",
                    ProcessCode: "R3",
                    Quantity: 123,
                },
            },
        },
    )
    params := Params{
        "process_code": "R1",
        "waste_code": "010101",
    }
    results, _ := r.Search(params)
    got := len(results)
    if got != 1 {
        t.Fatalf("Expected %d, got %d", 1, got)
    }
    result := results[0]
    want := "Test2"
    if result.Name != want {
        t.Fatalf("Expected %q, got %q", want, result.Name)
    }
}
