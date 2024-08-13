package filter_adapter

import "integration-git/main/pkg/result/domain"

type ResultFilterMatchType struct {
	matchType string
}

func NewResultFilterMatchType(matchType string) *ResultFilterMatchType {
	return &ResultFilterMatchType{
		matchType: matchType,
	}
}

func (f *ResultFilterMatchType) IsValid(result domain.Result) bool {
	return f.matchType == result.GetMatchType()
}
