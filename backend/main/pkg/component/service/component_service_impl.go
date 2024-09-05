package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
)

type ComponentUsecaseImpl struct {
	repo repository.ComponentRepository
}

func NewComponentUseCase(repo repository.ComponentRepository) ComponentService {
	return &ComponentUsecaseImpl{
		repo: repo,
	}
}

func (u *ComponentUsecaseImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return u.repo.FindByFilePath(filePath)
}

func (u *ComponentUsecaseImpl) FilterComponent(dto entities.ComponentFilterDTO) error {
	return u.repo.InsertComponentFilter(&dto)
}
