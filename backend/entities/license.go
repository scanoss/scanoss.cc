package entities

type LicenseFile struct {
	LicenseListVersion string    `json:"licenseListVersion"`
	Licenses           []License `json:"licenses"`
	ReleaseDate        string    `json:"releaseDate"`
}

type License struct {
	Name      string `json:"name"`
	LicenseId string `json:"licenseId"`
	Reference string `json:"reference"`
}
