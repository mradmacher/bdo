package bdo

import (
	"html/template"
	"io"
	"slices"
)

func formatWasteCode(value string, dangerous bool) string {
	formattedCode := value[:2] + " " + value[2:4] + " " + value[4:]
	if dangerous {
		formattedCode += "*"
	}
	return formattedCode
}

type InstallationSummaryView struct {
	Id	string
	Name string
	AddressLat float32
	AddressLng float32
	AddressLine1 string
	AddressLine2 string
	WasteCodes []string
	ProcessCodes []string
}

type CapabilityView struct {
	WasteCode string
	ProcessCode string
	Quantity int
}

type InstallationView struct {
	Id	string
	Name string
	AddressLine1 string
	AddressLine2 string
	Capabilities []CapabilityView
}

type InstallationsView struct {
  Installations []InstallationSummaryView
}


type Renderer struct {
	homeTemplate *template.Template
	installationsTemplate *template.Template
	installationSummaryTemplate *template.Template
}

func NewRenderer(templatesPath string) (*Renderer, error) {
	var err error
	renderer := Renderer{}

	renderer.homeTemplate, err = template.ParseFiles(templatesPath + "/index.html")
	if err != nil {
		return nil, err
	}
	renderer.installationsTemplate, err = template.ParseFiles(templatesPath + "/installations.gohtml")
	if err != nil {
		return nil, err
	}
	renderer.installationSummaryTemplate, err = template.ParseFiles(templatesPath + "/installation_summary.gohtml")
	if err != nil {
		return nil, err
	}

	return &renderer, nil
}

func (r *Renderer) RenderHome(w io.Writer, data any) error {
	return  r.homeTemplate.Execute(w, data)
}

func (r *Renderer) RenderInstallationSummary(w io.Writer, installation Installation) error {
	view := InstallationView {
		Id: installation.ID.Hex(),
		Name: installation.Name,
		AddressLine1: installation.Address.Line1,
		AddressLine2: installation.Address.Line2,
	}
	for _, c := range installation.Capabilities {
		view.Capabilities = append(view.Capabilities, CapabilityView{
			WasteCode: formatWasteCode(c.WasteCode, c.Dangerous),
			ProcessCode: c.ProcessCode,
			Quantity: c.Quantity,
		})
	}

	return r.installationSummaryTemplate.Execute(w, view)
}

func (r *Renderer) RenderInstallations(w io.Writer, installations []Installation) error {
	var result InstallationsView
	for _, installation := range installations {
	  summary := InstallationSummaryView {
		Id: installation.ID.Hex(),
		Name: installation.Name,
		AddressLat: installation.Address.Lat,
		AddressLng: installation.Address.Lng,
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
