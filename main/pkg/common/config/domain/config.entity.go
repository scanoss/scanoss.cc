package domain

type Config struct {
	Scanoss ScanossConfig `json:"scanoss"`
}

type ScanossConfig struct {
	ApiToken       string `json:"apiToken"`
	ApiUrl         string `json:"apiUrl"`
	ResultFilePath string `json:"resultFilePath"`
	ScanRoot       string `json:"scanRoot"`
}
