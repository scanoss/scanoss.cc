package services

import (
	"fmt"
	"integration-git/main/pkg/common/config/domain"
)

type ConfigService struct {
	repository domain.ConfigRepository
}

func NewConfigService(r domain.ConfigRepository) *ConfigService {
	return &ConfigService{
		repository: r,
	}
}

func (s ConfigService) ReadFromJson(path string) (domain.Config, error) {
	fmt.Printf("Path: %s\n", path)
	return s.repository.Read(path)
}
