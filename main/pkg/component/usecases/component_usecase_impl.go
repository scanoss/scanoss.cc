package usecases

import (
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/component/repositories"
)

type ComponentUsecaseImpl struct {
	repo repositories.ComponentRepository
}

func NewComponentUseCase(repo repositories.ComponentRepository) ComponentUsecase {
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
