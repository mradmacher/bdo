package bdo

import (
	"testing"
)

func testFindInstallation(t *testing.T, r *Repository) {
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
}
