package repository

import "github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/entities"

type ScanossBomRepository interface {
	Save() error
	Read() (entities.BomFile, error)
	Init() (entities.BomFile, error)
}
