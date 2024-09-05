package adapter

type ScanInputPathsDTO struct {
	Paths    []string `json:"paths"`
	BasePath string   `json:"basePath"`
}
