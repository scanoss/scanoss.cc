package entities

type Component struct {
	ID          string   `json:"id"`
	Lines       string   `json:"lines,omitempty"`
	OssLines    string   `json:"oss_lines,omitempty"`
	Matched     string   `json:"matched,omitempty"`
	FileHash    string   `json:"file_hash,omitempty"`
	SourceHash  string   `json:"source_hash,omitempty"`
	FileURL     string   `json:"file_url,omitempty"`
	Purl        []string `json:"purl,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Component   string   `json:"component,omitempty"`
	Version     string   `json:"version,omitempty"`
	Latest      string   `json:"latest,omitempty"`
	URL         string   `json:"url,omitempty"`
	Status      string   `json:"status,omitempty"`
	ReleaseDate string   `json:"release_date,omitempty"`
	File        string   `json:"file,omitempty"`
	URLHash     string   `json:"url_hash,omitempty"`
	URLStats    struct{} `json:"url_stats,omitempty"`
	Provenance  string   `json:"provenance,omitempty"`
	Licenses    []struct {
		Name             string `json:"name"`
		PatentHints      string `json:"patent_hints"`
		Copyleft         string `json:"copyleft"`
		ChecklistURL     string `json:"checklist_url"`
		OsadlUpdated     string `json:"osadl_updated"`
		Source           string `json:"source"`
		URL              string `json:"url"`
		IncompatibleWith string `json:"incompatible_with,omitempty"`
	} `json:"licenses,omitempty"`
	Health struct {
		CreationDate string `json:"creation_date"`
		LastUpdate   string `json:"last_update"`
		LastPush     string `json:"last_push"`
		Stars        int    `json:"stars"`
		Issues       int    `json:"issues"`
		Forks        int    `json:"forks"`
	} `json:"health"`
	Dependencies []interface{} `json:"dependencies"`
	Copyrights   []struct {
		Name   string `json:"name"`
		Source string `json:"source"`
	} `json:"copyrights"`
	Vulnerabilities []interface{} `json:"vulnerabilities"`
	Server          struct {
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

type DeclaredComponent struct {
	Name string `json:"name"`
	Purl string `json:"purl"`
}
