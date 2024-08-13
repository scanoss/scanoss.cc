package cmd

import (
	"integration-git/main/pkg/component"
	"integration-git/main/pkg/file"
	"integration-git/main/pkg/git"
	"integration-git/main/pkg/result"
	"integration-git/main/pkg/scan"
)

type Server struct {
	FileModule      *file.Module
	GitModule       *git.Module
	ResultModule    *result.Module
	ComponentModule *component.Module
	ScanModule      *scan.Module
}

// NewApp creates a new Backend container with all modules
func NewServer() *Server {
	return &Server{
		FileModule:      file.NewModule(),
		GitModule:       git.NewModule(),
		ResultModule:    result.NewModule(),
		ComponentModule: component.NewModule(),
		ScanModule:      scan.NewModule(),
	}
}
