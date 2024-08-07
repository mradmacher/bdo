package repo

import (
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
func assertEqualStr(t *testing.T, got, want, name string) {
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}
func assertEqual[K comparable](t *testing.T, got, want K, name string) {
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

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

		want := Installation{
			Name: "ToBeFound",
			Address: Address{
				Line1:     "Address1",
				Line2:     "Address2",
				StateCode: "10",
				Lat:       "12.00",
				Lng:       "10.00",
			},
			Capabilities: []Capability{
				Capability{
					WasteCode:    "010203",
					Dangerous:    true,
					ProcessCode:  "R1",
					ActivityCode: "Z",
					Quantity:     123,
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
		id, err := want.Add(r)
		assertNoError(t, err)

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

		t.Run("existing with metadata", func(t *testing.T) {
			var got Installation

			err := r.Find(id, &got)
			assertNoError(t, err)
			assertEqualStr(t, got.Name, want.Name, "Installation.Name")
			assertEqualStr(t, got.Address.Line1, want.Address.Line1, "Installation.Address.Line1")
			assertEqualStr(t, got.Address.Line2, want.Address.Line2, "Installation.Address.Line1")
			assertEqualStr(t, got.Address.Lng, want.Address.Lng, "Installation.Address.Lng")
			assertEqualStr(t, got.Address.Lat, want.Address.Lat, "Installation.Address.Lat")
		})

		t.Run("existing with capabilities", func(t *testing.T) {
			var got Installation

			err := r.Find(id, &got)
			assertNoError(t, err)
			assertEqual(t, len(got.Capabilities), 3, "len(Capabilities)")
			for _, w := range want.Capabilities {
				found := false
				for _, g := range got.Capabilities {
					if g.WasteCode == w.WasteCode &&
						g.Dangerous == w.Dangerous &&
						g.ProcessCode == w.ProcessCode &&
						g.ActivityCode == w.ActivityCode &&
						g.Quantity == w.Quantity {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected capability not found %v", w)
				}
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
