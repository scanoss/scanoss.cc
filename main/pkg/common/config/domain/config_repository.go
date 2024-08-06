package domain

type ConfigRepository interface {
	Read(path string) (Config, error)
}
