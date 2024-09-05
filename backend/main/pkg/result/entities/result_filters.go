package entities

type ResultFilterAND struct {
	filters []ResultFilter
}

func NewResultFilterAND() *ResultFilterAND {
	return &ResultFilterAND{}
}

func (f *ResultFilterAND) AddFilter(filter ResultFilter) {
	f.filters = append(f.filters, filter)
}

func (f *ResultFilterAND) IsValid(result Result) bool {
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

type ResultFilterMatchType struct {
	matchType string
}

func NewResultFilterMatchType(matchType string) *ResultFilterMatchType {
	return &ResultFilterMatchType{
		matchType: matchType,
	}
}

func (f *ResultFilterMatchType) IsValid(result Result) bool {
	return f.matchType == result.GetMatchType()
}

type ResultFilterFactory struct {
}

func NewResultFilterFactory() *ResultFilterFactory {
	return &ResultFilterFactory{}
}

func (f *ResultFilterFactory) Create(dto *RequestResultDTO) ResultFilter {
	if dto == nil {
		return nil
	}
	filterAND := NewResultFilterAND()

	if dto.MatchType != "" {
		filterAND.AddFilter(NewResultFilterMatchType(string(dto.MatchType)))
	}

	return filterAND
}
