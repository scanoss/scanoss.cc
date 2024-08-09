package services

import (
	"integration-git/main/pkg/component/domain"
)

type ComponentService struct {
	r domain.ComponentRepository
}

func NewComponentService(r domain.ComponentRepository) *ComponentService {
	return &ComponentService{
		r: r,
	}
}

func (s *ComponentService) GetComponentByPath(filePath string) (domain.Component, error) {
	return s.r.FindByFilePath(filePath)
}
