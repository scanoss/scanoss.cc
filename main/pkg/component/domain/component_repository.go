package domain

type ComponentRepository interface {
	FindByFilePath(path string) (Component, error)
}
