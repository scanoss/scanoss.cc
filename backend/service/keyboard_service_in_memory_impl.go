package service

import "github.com/scanoss/scanoss.lui/backend/entities"

type KeyboardServiceInMemoryImpl struct {
}

func NewKeyboardServiceInMemoryImpl() KeyboardService {
	return &KeyboardServiceInMemoryImpl{}
}

func (k *KeyboardServiceInMemoryImpl) GetShortcuts() []entities.Shortcut {
	return entities.DefaultShortcuts
}

func (k *KeyboardServiceInMemoryImpl) GetGroupedShortcuts() map[entities.Group][]entities.Shortcut {
	groupedShortcuts := make(map[entities.Group][]entities.Shortcut)
	for _, shortcut := range entities.DefaultShortcuts {
		groupedShortcuts[shortcut.Group] = append(groupedShortcuts[shortcut.Group], shortcut)
	}
	return groupedShortcuts
}
