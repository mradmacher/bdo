package bdo

import (
	"embed"
	"html/template"
	"io"
	"slices"
)

var (
	//go:embed "templates"
	templates embed.FS
)

func formatWasteCode(value string, dangerous bool) string {
	formattedCode := value[:2] + " " + value[2:4] + " " + value[4:]
	if dangerous {
		formattedCode += "*"
	}
	return formattedCode
}

type InstallationSummaryView struct {
	Id           int64
	Name         string
	AddressLat   string
	AddressLng   string
	AddressLine1 string
	AddressLine2 string
	WasteCodes   []string
	ProcessCodes []string
}

type CapabilityView struct {
	WasteCode    string
	ProcessCode  string
	ActivityCode string
	Quantity     int
}

type InstallationView struct {
	Id           int64
	Name         string
	AddressLine1 string
	AddressLine2 string
	Capabilities []CapabilityView
}

type InstallationsView struct {
	Installations []InstallationSummaryView
}

type CapabilitiesView struct {
	Capabilities []CapabilityView
}

type Renderer struct {
	homeTemplate                *template.Template
	installationsTemplate       *template.Template
	capabilitiesTemplate        *template.Template
}

func NewRenderer() (*Renderer, error) {
	var err error
	renderer := Renderer{}

	renderer.homeTemplate, err = template.ParseFS(templates, "templates/index.html")
	if err != nil {
		return nil, err
	}
	renderer.installationsTemplate, err = template.ParseFS(templates, "templates/installations.gohtml")
	if err != nil {
		return nil, err
	}
	renderer.capabilitiesTemplate, err = template.ParseFS(templates, "templates/capabilities.gohtml")
	if err != nil {
		return nil, err
	}

	return &renderer, nil
}

func (r *Renderer) RenderHome(w io.Writer, data any) error {
	return r.homeTemplate.Execute(w, data)
}

func (r *Renderer) RenderCapabilities(w io.Writer, capabilities []Capability) error {
	var result CapabilitiesView

	for _, c := range capabilities {
		result.Capabilities = append(result.Capabilities, CapabilityView{
			WasteCode:    formatWasteCode(c.WasteCode, c.Dangerous),
			ProcessCode:  c.ProcessCode,
			ActivityCode: c.ActivityCode,
			Quantity:     c.Quantity,
		})
	}

	return r.capabilitiesTemplate.Execute(w, result)
}

func (r *Renderer) RenderInstallations(w io.Writer, installations []*Installation) error {
	var result InstallationsView
	for _, installation := range installations {
		summary := InstallationSummaryView{
			Id:           installation.Id,
			Name:         installation.Name,
			AddressLat:   installation.Address.Lat,
			AddressLng:   installation.Address.Lng,
			AddressLine1: installation.Address.Line1,
			AddressLine2: installation.Address.Line2,
		}

		for _, c := range installation.Capabilities {
			formattedCode := formatWasteCode(c.WasteCode, c.Dangerous)
			if !slices.Contains(summary.WasteCodes, formattedCode) {
				summary.WasteCodes = append(summary.WasteCodes, formattedCode)
			}
			if !slices.Contains(summary.ProcessCodes, c.ProcessCode) {
				summary.ProcessCodes = append(summary.ProcessCodes, c.ProcessCode)
			}
		}
		result.Installations = append(result.Installations, summary)
	}
	return r.installationsTemplate.Execute(w, result)
}
