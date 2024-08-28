package infraestructure

import "integration-git/main/pkg/common/scanoss_bom/application/entities"

type ScanossBomRepository interface {
	Save(bomFile entities.BomFile) error
	Read() (entities.BomFile, error)
	Init() (entities.BomFile, error)
}
