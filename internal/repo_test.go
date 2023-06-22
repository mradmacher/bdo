package bdo

import (
    "testing"
    "golang.org/x/exp/slices"
    "github.com/joho/godotenv"
)

func setupSuite(t *testing.T) (func(*testing.T), *DbClient) {
    if err := godotenv.Load("../.env"); err != nil {
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

func setupTest(t *testing.T, db *DbClient) (func(*testing.T), *InstallationRepo) {
    repo := db.NewInstallationRepo()
    repo.Purge()
    repo.Add(
        &Installation{
            Name: "Test1",
            Address: Address{
              StateCode: "11",
            },
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
                Capability{
                    WasteCode: "010101",
                    ProcessCode: "R1",
                    Quantity: 123,
                },
            },
        },
    )
    repo.Add(
        &Installation{
            Name: "Test2",
            Address: Address{
              StateCode: "10",
            },
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
    repo.Add(
        &Installation{
            Name: "Test3",
            Address: Address{
              StateCode: "11",
            },
            Capabilities: []Capability{
                Capability{
                    WasteCode: "030303",
                    ProcessCode: "R1",
                    Quantity: 123,
                },
            },
        },
    )

    return func(t *testing.T) {
        repo.Purge()
    }, repo
}

func TestSearchReturnsAllDocuments(t *testing.T) {
    teardownSuite, db := setupSuite(t)
    defer teardownSuite(t)
    teardownTest, repo := setupTest(t, db)
    defer teardownTest(t)

    testCases := []struct {
      params SearchParams
      want []string
    }{
        {
            SearchParams{
                "process_code": "R1",
            },
            []string{"Test1", "Test2", "Test3"},
        },
        {
            SearchParams{
                "process_code": "R100",
            },
            []string{},
        },
        {
            SearchParams{
                "waste_code": "010101",
            },
            []string{"Test1", "Test2"},
        },
        {
            SearchParams{
                "waste_code": "010203",
            },
            []string{},
        },
        {
            SearchParams{
                "state_code": "11",
            },
            []string{"Test1", "Test3"},
        },
        {
            SearchParams{
                "state_code": "100",
            },
            []string{},
        },
        {
            SearchParams{
                "process_code": "R1",
                "waste_code": "010101",
            },
            []string{"Test1", "Test2"},
        },
        {
            SearchParams{
                "process_code": "R1",
                "waste_code": "010101",
                "state_code": "10",
            },
            []string{"Test2"},
        },
    }

    for _, tc := range testCases {
        results, _ := repo.Search(tc.params)
        got_count := len(results)
        want_count := len(tc.want)
        if got_count != len(tc.want) {
            t.Errorf("len(Search(%v)) = %d; want %d", tc.params, got_count, want_count)
        }
        for _, result := range results {
            if !slices.Contains(tc.want, result.Name) {
                t.Errorf("Search(%v)) returned %q; want %q", tc.params, result.Name, tc.want)
            }
        }
    }
}
