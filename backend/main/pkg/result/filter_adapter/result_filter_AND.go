package filter_adapter

import (
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
)

type ResultFilterAND struct {
	filters []common.ResultFilter
}

func NewResultFilterAND() *ResultFilterAND {
	return &ResultFilterAND{}
}

func (f *ResultFilterAND) AddFilter(filter common.ResultFilter) {
	f.filters = append(f.filters, filter)
}

func (f *ResultFilterAND) IsValid(result domain.Result) bool {
	for _, filter := range f.filters {
		if !filter.IsValid(result) {
			return false
		}
	}
	return true
}

func (f *ResultFilterAND) HasFilters() bool {
	return len(f.filters) > 0
}
