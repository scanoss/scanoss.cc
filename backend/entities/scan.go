package entities

type ScanResponse struct {
	Output    string `json:"output,omitempty"`
	ErrOutput string `json:"error_output,omitempty"`
	Error     error  `json:"error,omitempty"`
}
