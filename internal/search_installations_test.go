package bdo

import (
	"fmt"
	"golang.org/x/exp/slices"
	"testing"
)

func testSearchInstallations(t *testing.T, r *Repository) {
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
}
