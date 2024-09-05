package service

import "github.com/scanoss/scanoss.lui/backend/main/internal/common/config/entities"

type ConfigService interface {
	ReadConfig() (entities.Config, error)
}
