package domain

type Config struct {
	Scanoss struct {
		ApiToken string `json:"apiToken"`
		ApiUrl   string `json:"apiUrl"`
	} `json:"scanoss"`
}
