package infraestructure

import "integration-git/main/pkg/common/config/domain"

type ConfigRepository interface {
	Read() (domain.Config, error)
	Save(config *domain.Config)
	Init() error
}
