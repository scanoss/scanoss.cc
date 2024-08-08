package domain

type Config struct {
	Scanoss struct {
		ApiToken       string `json:"apiToken"`
		ApiUrl         string `json:"apiUrl"`
		ResultFilePath string `json:"resultFilePath"`
	} `json:"scanoss"`
}
