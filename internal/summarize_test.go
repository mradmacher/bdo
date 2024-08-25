package bdo

import (
	"fmt"
	"testing"
)

func testSummarize(t *testing.T, r *Repository) {
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
}
