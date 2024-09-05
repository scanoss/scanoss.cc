package repository

import "github.com/scanoss/scanoss.lui/backend/main/internal/common/config/entities"

type ConfigRepository interface {
	Read() (entities.Config, error)
	Save(config *entities.Config)
	Init()
}
