package service

import "github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/repository"

type ScanossBomServiceImp struct {
	repository repository.ScanossBomRepository
}

func NewScanossBomService(r repository.ScanossBomRepository) *ScanossBomServiceImp {
	return &ScanossBomServiceImp{
		repository: r,
	}
}

func (us *ScanossBomServiceImp) Save() error {
	return us.repository.Save()
}
