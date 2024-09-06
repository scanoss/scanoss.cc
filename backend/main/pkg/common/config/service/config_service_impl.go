package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/repository"
)

type ConfigServiceImpl struct {
	repo repository.ConfigRepository
}

func NewConfigService(repo repository.ConfigRepository) ConfigService {
	return &ConfigServiceImpl{
		repo: repo,
	}
}

func (s *ConfigServiceImpl) ReadConfig() (entities.Config, error) {
	return s.repo.Read()
}
