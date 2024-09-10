package bdo

import (
	"embed"
	"errors"
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"slices"
	"strings"
)

var (
	//go:embed "templates"
	templates embed.FS

	//go:embed "seeds"
	seeds embed.FS
)

func formatWasteCode(value string, dangerous bool) string {
	formattedCode := value[:2] + " " + value[2:4] + " " + value[4:]
	if dangerous {
		formattedCode += "*"
	}
	return formattedCode
}

func normalizeWasteCode(value string) string {
	return strings.Replace(strings.ReplaceAll(value, " ", ""), "*", "", 1)
}

type MaterialSpecs struct {
	Items []MaterialSpec
}

type MaterialSpec struct {
	Code  string
	Name  string
	Desc  string
	Items []MaterialSpec
}

func (m *MaterialSpecs) findGroup(code string) *MaterialSpec {
	for _, i := range m.Items {
		found := i.findGroup(code)
		if found != nil {
			return found
		}
	}
	return nil
}

func (m *MaterialSpec) findGroup(code string) *MaterialSpec {
	groupCode := code[:4]

	if m.Code == code {
		return m
	}

	for _, i := range m.Items {
		found := i.findGroup(groupCode)
		if found != nil {
			return found
		}
	}

	return nil
}

func loadMaterials() (*MaterialSpecs, error) {
	var m []MaterialSpec
	buffer, err := seeds.ReadFile("seeds/materials.yaml")
	if err != nil {
		return nil, errors.Join(errors.New("Problem reading materials file"), err)
	}
	err = yaml.Unmarshal(buffer, &m)
	if err != nil {
		return nil, errors.Join(errors.New("Problem unmarshaling materials file"), err)
	}

	return &MaterialSpecs{Items: m}, nil
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

type MaterialGroupView struct {
	Code      string
	Name      string
	Desc      string
	Materials []*MaterialView
}

type MaterialView struct {
	Name     string
	Code     string
	Selected bool
}

type CapabilityView struct {
	Id             int64
	WasteCode      string
	ProcessCode    string
	ActivityCode   string
	Quantity       int
	Materials      []MaterialView
	MaterialGroups []MaterialGroupView
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
	capabilitiesSummaryTemplate *template.Template

	materialSpecs *MaterialSpecs
}

func NewRenderer() (*Renderer, error) {
	var err error
	renderer := Renderer{}

	funcMap := template.FuncMap{
		"normalize": normalizeWasteCode,
	}

	renderer.homeTemplate, err = template.ParseFS(templates, "templates/index.html")
	if err != nil {
		return nil, err
	}
	renderer.installationsTemplate, err = template.New("installations.gohtml").Funcs(funcMap).ParseFS(templates, "templates/installations.gohtml")
	if err != nil {
		return nil, err
	}
	renderer.capabilitiesTemplate, err = template.ParseFS(templates, "templates/capabilities.gohtml")
	if err != nil {
		return nil, err
	}
	renderer.capabilitiesSummaryTemplate, err = template.ParseFS(templates, "templates/capabilities_summary.gohtml")
	if err != nil {
		return nil, err
	}
	renderer.materialSpecs, err = loadMaterials()
	if err != nil {
		return nil, err
	}

	return &renderer, nil
}

func (r *Renderer) RenderHome(w io.Writer, data any) error {
	return r.homeTemplate.Execute(w, data)
}

func (r *Renderer) RenderCapabilitiesSummary(w io.Writer, capabilities []Capability) error {
	var result CapabilitiesView

	for _, c := range capabilities {
		cv := CapabilityView{
			WasteCode:   formatWasteCode(c.WasteCode, c.Dangerous),
			ProcessCode: c.ProcessCode,
			Quantity:    c.Quantity,
		}
		result.Capabilities = append(result.Capabilities, cv)
	}

	return r.capabilitiesSummaryTemplate.Execute(w, result)
}

func (r *Renderer) RenderCapabilities(w io.Writer, capabilities []Capability) error {
	var result CapabilitiesView

	for _, c := range capabilities {
		cv := CapabilityView{
			Id:           c.Id,
			WasteCode:    formatWasteCode(c.WasteCode, c.Dangerous),
			ProcessCode:  c.ProcessCode,
			ActivityCode: c.ActivityCode,
			Quantity:     c.Quantity,
		}
		mgroups := make(map[string]*MaterialGroupView)
		for _, m := range c.Materials {
			mg := r.materialSpecs.findGroup(m.Code)
			if mg != nil {
				if mgroups[mg.Code] == nil {
					mgv := MaterialGroupView{Code: mg.Code, Name: mg.Name, Desc: mg.Desc}
					for _, i := range mg.Items {
						mgv.Materials = append(mgv.Materials, &MaterialView{Code: i.Code, Name: i.Name, Selected: false})
					}
					mgroups[mg.Code] = &mgv
				}
				for _, mgm := range mgroups[mg.Code].Materials {
					if mgm.Code == m.Code {
						mgm.Selected = true
					}
				}
			}
			//cv.Materials = append(cv.Materials, MaterialView{Code: m.Code})
		}
		for _, v := range mgroups {
			cv.MaterialGroups = append(cv.MaterialGroups, *v)
		}
		result.Capabilities = append(result.Capabilities, cv)
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
