package service

import "github.com/scanoss/scanoss.cc/backend/entities"

type KeyboardService interface {
	GetShortcuts() []entities.Shortcut
	GetGroupedShortcuts() map[entities.Group][]entities.Shortcut
}
