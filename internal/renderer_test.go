package bdo

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderInstallations(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Errorf("Error creating the renderer: %v", err)
	}

	t.Run("renders empty string when no installations", func(t *testing.T) {
		buffer := bytes.Buffer{}
		renderer.RenderInstallations(&buffer, []*Installation{})
		got := buffer.String()
		want := "\nBrak instalacji spełniających podane kryteria\n\n\n"

		if got != want {
			t.Errorf("got: %q, want %q", got, want)
		}
	})

	t.Run("renders names", func(t *testing.T) {
		collection := []*Installation{
			&Installation{
				Name: "Test1",
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
			},
			&Installation{
				Name: "Test2",
				Address: Address{
					StateCode: "12",
				},
				Capabilities: []Capability{
					Capability{
						WasteCode:   "030303",
						ProcessCode: "R3",
						Quantity:    456,
					},
				},
			},
		}

		buffer := bytes.Buffer{}
		renderer.RenderInstallations(&buffer, collection)
		got := buffer.String()
		want := "Test2"
		if !strings.Contains(got, want) {
			t.Errorf("Expected %q to contain %q", got, want)
		}
		want = "Test1"
		if !strings.Contains(got, want) {
			t.Errorf("Expected %q to contain %q", got, want)
		}
	})
}
