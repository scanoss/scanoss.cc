package filter_adapter

import (
	"integration-git/main/pkg/result/common"
)

type ResultFilterFactory struct {
}

func NewResultFilterFactory() *ResultFilterFactory {
	return &ResultFilterFactory{}
}

func (f *ResultFilterFactory) Create(dto *common.RequestResultDTO) common.ResultFilter {
	if dto == nil {
		return nil
	}
	filterAND := NewResultFilterAND()

	if dto.MatchType != "" {
		filterAND.AddFilter(NewResultFilterMatchType(string(dto.MatchType)))
	}

	return filterAND
}
