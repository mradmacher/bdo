package bdo

import (
	"testing"
)

func testSearchCapabilities(t *testing.T, r *Repository) {
	cap1 := Capability{
		WasteCode:    "010101",
		ProcessCode:  "R3",
		ActivityCode: "Z",
		Quantity:     789,
		Materials: []Material{
			Material{Code: "TSTT01"},
			Material{Code: "TSTT05"},
		},
	}
	cap2 := Capability{
		WasteCode:    "010101",
		ProcessCode:  "R3",
		ActivityCode: "PR",
		Quantity:     789,
		Materials: []Material{
			Material{Code: "TSTT03"},
		},
	}
	cap3 := Capability{
		WasteCode:    "010101",
		ProcessCode:  "R2",
		ActivityCode: "PR",
		Quantity:     987,
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
}
