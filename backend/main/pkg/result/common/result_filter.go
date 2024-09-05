package common

import "integration-git/main/pkg/result/domain"

type ResultFilter interface {
	IsValid(result domain.Result) bool
}
