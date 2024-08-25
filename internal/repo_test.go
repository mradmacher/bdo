package bdo

import (
	"fmt"
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

func assertHasCapability(t *testing.T, inst Installation, want Capability) {
	found := false
	for _, c := range inst.Capabilities {
		if c.WasteCode == want.WasteCode &&
			c.Dangerous == want.Dangerous &&
			c.ProcessCode == want.ProcessCode &&
			c.ActivityCode == want.ActivityCode &&
			c.Quantity == want.Quantity {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected %q to has capability %v", inst.Name, want)
	}
}
func assertIncludeCapability(t *testing.T, collection []Capability, want Capability) {
	found := false
	for _, c := range collection {
		if c.WasteCode == want.WasteCode &&
			c.Dangerous == want.Dangerous &&
			c.ProcessCode == want.ProcessCode &&
			c.ActivityCode == want.ActivityCode &&
			c.Quantity == want.Quantity {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected %v to has capability %v", collection, want)
	}
}

func setupSuite(t *testing.T) (func(*testing.T), *Repository) {
	if err := godotenv.Load("../.test.env"); err != nil {
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
	return func(t *testing.T) {
		r.Purge()
	}
}

func TestRepo(t *testing.T) {
	teardownSuite, r := setupSuite(t)
	defer teardownSuite(t)

	t.Run("Summarize", func(t *testing.T) {
		teardownTest := setupTest(t, r)
		defer teardownTest(t)

		inst1 := Installation{
			Name: "Test1",
			Address: Address{
				StateCode: "10",
			},
			Capabilities: []Capability{
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R1",
					Quantity:    100,
				},
				Capability{
					WasteCode:   "020202",
					Dangerous:   true,
					ProcessCode: "D2",
					Quantity:    100,
				},
			},
		}

		inst2 := Installation{
			Name: "Test2",
			Address: Address{
				StateCode: "10",
			},
			Capabilities: []Capability{
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R1",
					Quantity:    100,
				},
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R2",
					Quantity:    100,
				},
				Capability{
					WasteCode:   "020202",
					Dangerous:   true,
					ProcessCode: "D3",
					Quantity:    100,
				},
			},
		}

		inst3 := Installation{
			Name: "Test3",
			Address: Address{
				StateCode: "11",
			},
			Capabilities: []Capability{
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R1",
					Quantity:    100,
				},
				Capability{
					WasteCode:   "010101",
					ProcessCode: "R2",
					Quantity:    100,
				},
				Capability{
					WasteCode:   "040404",
					ProcessCode: "R10",
					Quantity:    100,
				},
			},
		}

		inst1.Add(r)
		inst2.Add(r)
		inst3.Add(r)

		testCases := []struct {
			params SearchParams
			want   []Capability
		}{
			{
				SearchParams{
					"waste_code":   "010101",
					"process_code": "R1",
				},
				[]Capability{
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R1",
						Quantity:    300,
					},
				},
			},
			{
				SearchParams{
					"process_code": "R1",
				},
				[]Capability{
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R1",
						Quantity:    300,
					},
				},
			},
			{
				SearchParams{
					"waste_code": "010101",
				},
				[]Capability{
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R1",
						Quantity:    300,
					},
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R2",
						Quantity:    200,
					},
				},
			},
			{
				SearchParams{
					"waste_code": "010101",
					"state_code": "10",
				},
				[]Capability{
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R1",
						Quantity:    200,
					},
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R2",
						Quantity:    100,
					},
				},
			},
			{
				SearchParams{
					"waste_code": "020202",
				},
				[]Capability{
					Capability{
						WasteCode:   "020202",
						Dangerous:   true,
						ProcessCode: "D2",
						Quantity:    100,
					},
					Capability{
						WasteCode:   "020202",
						Dangerous:   true,
						ProcessCode: "D3",
						Quantity:    100,
					},
				},
			},
			{
				SearchParams{
					"state_code": "10",
				},
				[]Capability{
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R1",
						Quantity:    200,
					},
					Capability{
						WasteCode:   "010101",
						Dangerous:   false,
						ProcessCode: "R2",
						Quantity:    100,
					},
					Capability{
						WasteCode:   "020202",
						Dangerous:   true,
						ProcessCode: "D2",
						Quantity:    100,
					},
					Capability{
						WasteCode:   "020202",
						Dangerous:   true,
						ProcessCode: "D3",
						Quantity:    100,
					},
				},
			},
			{
				SearchParams{
					"waste_code": "040404",
					"state_code": "10",
				},
				[]Capability{},
			},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%v", tc.params), func(t *testing.T) {
				var results []Capability
				var err error

				results, err = r.Summarize(tc.params)
				assertNoError(t, err)
				t.Run("returns summarized capabilities", func(t *testing.T) {
					assertEqual(t, len(results), len(tc.want), "len(Capabilites)")
					for _, want := range tc.want {
						assertIncludeCapability(t, results, want)
					}
				})
			})
		}
	})

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
				assertHasCapability(t, got, w)
			}
		})
	})

	t.Run("SearchCapabilities", func(t *testing.T) {
		teardownTest := setupTest(t, r)
		defer teardownTest(t)

		cap1 := Capability{
			WasteCode:   "010101",
			ProcessCode: "R3",
			ActivityCode: "Z",
			Quantity:    789,
		}
		cap2 := Capability{
			WasteCode:   "010101",
			ProcessCode: "R3",
			ActivityCode: "PR",
			Quantity:    789,
		}
		cap3 := Capability{
			WasteCode:   "010101",
			ProcessCode: "R2",
			ActivityCode: "PR",
			Quantity:    987,
		}

		inst := Installation{
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
				cap1,
				cap2,
				cap3,
			},
		}
		id, err := inst.Add(r)
		assertNoError(t, err)

		t.Run("not existing installation", func(t *testing.T) {
			var got []Capability
			sp := SearchParams{
				"waste_code": "010101",
			}
			got, err := r.SearchCapabilities(0, sp)
			assertNoError(t, err)
			if len(got) > 0 {
				t.Errorf("Expected empty collection but got %v", got)
			}
		})

		t.Run("not existing code", func(t *testing.T) {
			var got []Capability
			sp := SearchParams{
				"waste_code": "101010",
			}
			got, err := r.SearchCapabilities(id, sp)
			assertNoError(t, err)
			if len(got) > 0 {
				t.Errorf("Expected empty collection but got %v", got)
			}
		})

		t.Run("existing with capabilities", func(t *testing.T) {
			var got []Capability
			want := []Capability{cap1, cap2, cap3}
			sp := SearchParams{
				"waste_code": "010101",
			}
			got, err := r.SearchCapabilities(id, sp)

			assertNoError(t, err)
			assertEqual(t, len(got), 3, "len(Capabilities)")
			for _, w := range want {
				assertIncludeCapability(t, got, w)
			}
		})
	})

	t.Run("Search", func(t *testing.T) {
		teardownTest := setupTest(t, r)
		defer teardownTest(t)

		inst1 := Installation{
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
					Dangerous:   true,
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
		inst1.Add(r)

		inst2 := Installation{
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
		inst2.Add(r)

		inst3 := Installation{
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
		inst3.Add(r)

		installations := map[string]*Installation{
			"Test1": &inst1,
			"Test2": &inst2,
			"Test3": &inst3,
		}

		testCases := []struct {
			params SearchParams
			want   []*Installation
		}{
			{
				SearchParams{
					"process_code": "R1",
				},
				[]*Installation{&inst1, &inst2, &inst3},
			},
			{
				SearchParams{
					"process_code": "R100",
				},
				[]*Installation{},
			},
			{
				SearchParams{
					"waste_code": "010101",
				},
				[]*Installation{&inst1, &inst2},
			},
			{
				SearchParams{
					"waste_code": "010203",
				},
				[]*Installation{},
			},
			{
				SearchParams{
					"state_code": "11",
				},
				[]*Installation{&inst1, &inst3},
			},
			{
				SearchParams{
					"state_code": "100",
				},
				[]*Installation{},
			},
			{
				SearchParams{
					"process_code": "R1",
					"waste_code":   "010101",
				},
				[]*Installation{&inst1, &inst2},
			},
			{
				SearchParams{
					"process_code": "R1",
					"waste_code":   "010101",
					"state_code":   "10",
				},
				[]*Installation{&inst2},
			},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%v", tc.params), func(t *testing.T) {
				var results []*Installation
				var err error

				results, err = r.Search(tc.params)
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				t.Run("returns correct installations", func(t *testing.T) {
					got_count := len(results)
					want_count := len(tc.want)
					if got_count != len(tc.want) {
						t.Errorf("len(Search(%v)) = %d; want %d", tc.params, got_count, want_count)
					}

					for _, result := range results {
						var wantedNames []string
						for _, inst := range tc.want {
							wantedNames = append(wantedNames, inst.Name)
						}
						if !slices.Contains(wantedNames, result.Name) {
							t.Errorf("Search(%v)) returned %q; want %q", tc.params, result.Name, wantedNames)
						}
					}
				})

				t.Run("returns installation capabilities", func(t *testing.T) {
					for _, result := range results {
						t.Logf("%v", result)
						got := len(result.Capabilities)
						want := len(installations[result.Name].Capabilities)
						if got != want {
							t.Errorf("len(%s.Capabilites) = %d; want %d", result.Name, got, want)
						}
						for _, want := range installations[result.Name].Capabilities {
							assertHasCapability(t, *result, want)
						}
					}
				})
			})
		}
	})
}
