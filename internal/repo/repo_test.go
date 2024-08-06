package repo

import (
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"testing"
)

func setupSuite(t *testing.T) (func(*testing.T), *Repository) {
	if err := godotenv.Load("../../.test.env"); err != nil {
		t.Fatalf("No .env file found")
	}

	r := Repository{}
	err := r.Connect()
	if err != nil {
		t.Fatalf("Problem with DB connection: %v", err)
	}
	return func(t *testing.T) {
		r.Disconnect()
	}, &r
}

func setupTest(t *testing.T, r *Repository) func(*testing.T) {
	var inst Installation

	inst = Installation{
		Name: "Test1",
		Address: Address{
			StateCode: "11",
		},
		Capabilities: []Capability{
			Capability{
				WasteCode:   "020202",
				ProcessCode: "R2",
				Quantity:    123,
			},
			Capability{
				WasteCode:   "010101",
				ProcessCode: "R2",
				Quantity:    123,
			},
			Capability{
				WasteCode:   "010101",
				ProcessCode: "R1",
				Quantity:    123,
			},
		},
	}
	inst.Add(r)

	inst = Installation{
		Name: "Test2",
		Address: Address{
			StateCode: "10",
		},
		Capabilities: []Capability{
			Capability{
				WasteCode:   "010101",
				ProcessCode: "R1",
				Quantity:    123,
			},
			Capability{
				WasteCode:   "020202",
				ProcessCode: "R2",
				Quantity:    123,
			},
		},
	}
	inst.Add(r)

	inst = Installation{
		Name: "Test3",
		Address: Address{
			StateCode: "11",
		},
		Capabilities: []Capability{
			Capability{
				WasteCode:   "030303",
				ProcessCode: "R1",
				Quantity:    123,
			},
		},
	}
	inst.Add(r)

	return func(t *testing.T) {
		r.Purge()
	}
}

func TestRepo(t *testing.T) {
	teardownSuite, r := setupSuite(t)
	defer teardownSuite(t)

	t.Run("Find", func(t *testing.T) {
		teardownTest := setupTest(t, r)
		defer teardownTest(t)

		inst := Installation{
			Name: "ToBeFound",
			Address: Address{
				StateCode: "10",
			},
			Capabilities: []Capability{
				Capability{
					WasteCode:   "010203",
					ProcessCode: "R1",
					Quantity:    123,
				},
				Capability{
					WasteCode:   "030201",
					ProcessCode: "R2",
					Quantity:    456,
				},
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R3",
					Quantity:    789,
				},
			},
		}
		id, err := inst.Add(r)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		t.Run("not existing", func(t *testing.T) {
			var got Installation
			err := r.Find(0, &got)
			if err == nil {
				t.Errorf("Expected error")
			}
			if !IsRecordNotFound(err) {
				t.Errorf("Expected no data could be found but got %v error", err)
			}
		})

		t.Run("existing", func(t *testing.T) {
			var got Installation

			err := r.Find(id, &got)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	})

	t.Run("Search", func(t *testing.T) {
		teardownTest := setupTest(t, r)
		defer teardownTest(t)

		testCases := []struct {
			params SearchParams
			want   []string
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
					"waste_code":   "010101",
				},
				[]string{"Test1", "Test2"},
			},
			{
				SearchParams{
					"process_code": "R1",
					"waste_code":   "010101",
					"state_code":   "10",
				},
				[]string{"Test2"},
			},
		}

		for _, tc := range testCases {
			var results []Installation
			err := r.Search(tc.params, &results)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
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
	})
}
