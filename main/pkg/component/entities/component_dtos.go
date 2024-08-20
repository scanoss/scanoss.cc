package entities

type License interface {
}

type LicenseDTO struct {
	Name             string `json:"name"`
	PatentHints      string `json:"patent_hints"`
	Copyleft         string `json:"copyleft"`
	ChecklistURL     string `json:"checklist_url"`
	OsadlUpdated     string `json:"osadl_updated"`
	Source           string `json:"source"`
	URL              string `json:"url"`
	IncompatibleWith string `json:"incompatible_with,omitempty"`
}

func NewLicenseDTO(name string, patentHints string, copyleft string, checklistURL string, osadlUpdated string, source string, url string, incompatibleWith string) *LicenseDTO {
	return &LicenseDTO{
		Name:             name,
		URL:              url,
		ChecklistURL:     checklistURL,
		Copyleft:         copyleft,
		PatentHints:      patentHints,
		IncompatibleWith: incompatibleWith,
		OsadlUpdated:     osadlUpdated,
		Source:           source,
	}
}

type ComponentDTO struct {
	ID          string    `json:"id"`
	Lines       string    `json:"lines,omitempty"`
	OssLines    string    `json:"oss_lines,omitempty"`
	Matched     string    `json:"matched,omitempty"`
	FileHash    string    `json:"file_hash,omitempty"`
	SourceHash  string    `json:"source_hash,omitempty"`
	FileURL     string    `json:"file_url,omitempty"`
	Purl        []string  `json:"purl,omitempty"`
	Vendor      string    `json:"vendor,omitempty"`
	Component   string    `json:"component,omitempty"`
	Version     string    `json:"version,omitempty"`
	Latest      string    `json:"latest,omitempty"`
	URL         string    `json:"url,omitempty"`
	Status      string    `json:"status,omitempty"`
	ReleaseDate string    `json:"release_date,omitempty"`
	File        string    `json:"file,omitempty"`
	URLHash     string    `json:"url_hash,omitempty"`
	URLStats    struct{}  `json:"url_stats,omitempty"`
	Provenance  string    `json:"provenance,omitempty"`
	Licenses    []License `json:"licenses,omitempty"`
	Server      struct {
		Version   string `json:"version,omitempty"`
		KbVersion struct {
			Monthly string `json:"monthly,omitempty"`
			Daily   string `json:"daily,omitempty"`
		} `json:"kb_version"`
		Hostname string `json:"hostname,omitempty"`
		Flags    string `json:"flags,omitempty"`
		Elapsed  string `json:"elapsed,omitempty"`
	} `json:"server"`
}

type FilterAction string

const (
	Include FilterAction = "include"
	Remove  FilterAction = "remove"
)

type ComponentFilterDTO struct {
	Path    string       `json:"path"`
	Purl    string       `json:"purl"`
	Usage   string       `json:"usage,omitempty"`
	Version string       `json:"version"`
	Action  FilterAction `json:"action"`
}
