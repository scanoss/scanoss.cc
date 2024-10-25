package keyboard

type KeyboardServiceInMemoryImpl struct {
}

func NewKeyboardServiceInMemoryImpl() KeyboardService {
	return &KeyboardServiceInMemoryImpl{}
}

func (k *KeyboardServiceInMemoryImpl) GetShortcuts() []Shortcut {
	return DefaultShortcuts
}

func (k *KeyboardServiceInMemoryImpl) GetGroupedShortcuts() map[Group][]Shortcut {
	groupedShortcuts := make(map[Group][]Shortcut)
	for _, shortcut := range DefaultShortcuts {
		groupedShortcuts[shortcut.Group] = append(groupedShortcuts[shortcut.Group], shortcut)
	}
	return groupedShortcuts
}
