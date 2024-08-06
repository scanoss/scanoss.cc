package config

type Config struct {
	Scanoss struct {
		ApiToken string `json:"apiToken"`
		ApiUrl   string `json:"apiToken"`
	}
}

func newFromJson(data []byte) {

}
