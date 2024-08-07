package use_cases

import "integration-git/main/pkg/common/config/domain"

type ReadConfigFileUseCase interface {
	ReadConfig(path string) (domain.Config, error)
}
