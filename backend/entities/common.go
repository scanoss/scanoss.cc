package entities

type ComponentRequest struct {
	Purl        string `json:"purl" validate:"required,min=1"`
	Requirement string `json:"requirement,omitempty"`
}

type StatusResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
