package mappers

import "github.com/scanoss/scanoss.cc/backend/entities"

type ComponentMapper interface {
	MapToComponentDTO(componentEntity entities.Component) entities.ComponentDTO
}
