package service

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type KeyboardService interface {
	GetShortcuts() []entities.Shortcut
	GetGroupedShortcuts() map[entities.Group][]entities.Shortcut
}
