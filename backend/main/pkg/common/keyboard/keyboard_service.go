package keyboard

type KeyboardService interface {
	GetShortcuts() []Shortcut
	GetGroupedShortcuts() map[Group][]Shortcut
}
