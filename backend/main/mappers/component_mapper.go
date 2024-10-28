package mappers

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type ComponentMapper interface {
	MapToComponentDTO(componentEntity entities.Component) entities.ComponentDTO
}
