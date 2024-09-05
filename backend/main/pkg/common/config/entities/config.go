package entities

type Config struct {
	ApiToken             string `json:"apiToken"`
	ApiUrl               string `json:"apiUrl"`
	ResultFilePath       string `json:"resultFilePath:omitempty"`
	ScanRoot             string `json:"scanRoot:omitempty"`
	ScanSettingsFilePath string `json:"scanSettingsFilePath:omitempty"`
}