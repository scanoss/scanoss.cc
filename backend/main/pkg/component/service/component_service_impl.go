package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
)

type ComponentServiceImpl struct {
	repo repository.ComponentRepository
}

func NewComponentService(repo repository.ComponentRepository) ComponentService {
	return &ComponentServiceImpl{
		repo: repo,
	}
}

func (u *ComponentServiceImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return u.repo.FindByFilePath(filePath)
}

func (u *ComponentServiceImpl) FilterComponent(dto entities.ComponentFilterDTO) error {
	return u.repo.InsertComponentFilter(&dto)
}
