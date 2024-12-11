package entities

import (
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// Action defines the type for keyboard actions
type Action string

const (
	ActionUndo                            Action = "undo"
	ActionRedo                            Action = "redo"
	ActionSave                            Action = "save"
	ActionConfirm                         Action = "confirm"
	ActionFocusSearch                     Action = "focusSearch"
	ActionSelectAll                       Action = "selectAll"
	ActionMoveUp                          Action = "moveUp"
	ActionMoveDown                        Action = "moveDown"
	ActionIncludeFileWithoutComments      Action = "includeFileWithoutComments"
	ActionIncludeFileWithComments         Action = "includeFileWithComments"
	ActionIncludeComponentWithoutComments Action = "includeComponentWithoutComments"
	ActionIncludeComponentWithComments    Action = "includeComponentWithComments"
	ActionDismissFileWithoutComments      Action = "dismissFileWithoutComments"
	ActionDismissFileWithComments         Action = "dismissFileWithComments"
	ActionDismissComponentWithoutComments Action = "dismissComponentWithoutComments"
	ActionDismissComponentWithComments    Action = "dismissComponentWithComments"
	ActionReplaceFileWithoutComments      Action = "replaceFileWithoutComments"
	ActionReplaceFileWithComments         Action = "replaceFileWithComments"
	ActionReplaceComponentWithoutComments Action = "replaceComponentWithoutComments"
	ActionReplaceComponentWithComments    Action = "replaceComponentWithComments"

	// View
	ActionSyncScrollPosition         Action = "toggleSyncScrollPosition"
	ActionShowKeyboardShortcutsModal Action = "showKeyboardShortcutsModal"

	// Scan
	ActionScanCwd         Action = "scanCurrentDirectory"
	ActionScanWithOptions Action = "scanWithOptions"
)

type Shortcut struct {
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	Accelerator            *keys.Accelerator `json:"accelerator"`
	AlternativeAccelerator *keys.Accelerator `json:"alternativeAccelerator,omitempty"`
	Keys                   string            `json:"keys"`
	Group                  Group             `json:"group,omitempty"`
	Action                 Action            `json:"action,omitempty"`
}

type Group string

const (
	GroupGlobal     Group = "Global"
	GroupNavigation Group = "Navigation"
	GroupActions    Group = "Actions"
	GroupView       Group = "View"
	GroupScan       Group = "Scan"
)

// This is necessary to bind the enum in main.go
var AllShortcutActions = []struct {
	Value  Action
	TSName string
}{
	{ActionUndo, "Undo"},
	{ActionRedo, "Redo"},
	{ActionSave, "Save"},
	{ActionConfirm, "Confirm"},
	{ActionFocusSearch, "FocusSearch"},
	{ActionSelectAll, "SelectAll"},
	{ActionMoveUp, "MoveUp"},
	{ActionMoveDown, "MoveDown"},
	{ActionIncludeFileWithoutComments, "IncludeFileWithoutComments"},
	{ActionIncludeFileWithComments, "IncludeFileWithComments"},
	{ActionIncludeComponentWithoutComments, "IncludeComponentWithoutComments"},
	{ActionIncludeComponentWithComments, "IncludeComponentWithComments"},
	{ActionDismissFileWithoutComments, "DismissFileWithoutComments"},
	{ActionDismissFileWithComments, "DismissFileWithComments"},
	{ActionDismissComponentWithoutComments, "DismissComponentWithoutComments"},
	{ActionDismissComponentWithComments, "DismissComponentWithComments"},
	{ActionReplaceFileWithoutComments, "ReplaceFileWithoutComments"},
	{ActionReplaceFileWithComments, "ReplaceFileWithComments"},
	{ActionReplaceComponentWithoutComments, "ReplaceComponentWithoutComments"},
	{ActionReplaceComponentWithComments, "ReplaceComponentWithComments"},
	{ActionSyncScrollPosition, "ToggleSyncScrollPosition"},
	{ActionShowKeyboardShortcutsModal, "ShowKeyboardShortcutsModal"},
	{ActionScanCwd, "ScanCurrentDirectory"},
	{ActionScanWithOptions, "ScanWithOptions"},
}

var DefaultShortcuts = []Shortcut{
	// Global
	{
		Name:        "Undo",
		Description: "Undo the last action",
		Accelerator: keys.CmdOrCtrl("z"),
		Keys:        "mod+z",
		Group:       GroupGlobal,
		Action:      ActionUndo,
	},
	{
		Name:        "Redo",
		Description: "Redo the last action",
		Accelerator: keys.Combo("z", keys.CmdOrCtrlKey, keys.ShiftKey),
		Keys:        "mod+shift+z",
		Group:       GroupGlobal,
		Action:      ActionRedo,
	},
	{
		Name:        "Save",
		Description: "Save the current changes",
		Accelerator: keys.CmdOrCtrl("s"),
		Keys:        "mod+s",
		Group:       GroupGlobal,
		Action:      ActionSave,
	},
	{
		Name:        "Confirm",
		Description: "Confirm the current action",
		Accelerator: keys.CmdOrCtrl("enter"),
		Keys:        "mod+enter",
		Group:       GroupGlobal,
		Action:      ActionConfirm,
	},
	{
		Name:        "Focus Search",
		Description: "Focus the search bar",
		Accelerator: keys.CmdOrCtrl("f"),
		Keys:        "mod+f",
		Group:       GroupGlobal,
		Action:      ActionFocusSearch,
	},
	{
		Name:        "Select All",
		Description: "Select all text",
		Accelerator: keys.CmdOrCtrl("a"),
		Keys:        "mod+a",
		Group:       GroupGlobal,
		Action:      ActionSelectAll,
	},

	// Navigation
	{
		Name:                   "Move Up",
		Description:            "Move the selection up",
		Accelerator:            keys.Key("up"),
		AlternativeAccelerator: keys.Key("k"),
		Keys:                   "k, up",
		Group:                  GroupNavigation,
		Action:                 ActionMoveUp,
	},
	{
		Name:                   "Move Down",
		Description:            "Move the selection down",
		Accelerator:            keys.Key("down"),
		AlternativeAccelerator: keys.Key("j"),
		Keys:                   "j, down",
		Group:                  GroupNavigation,
		Action:                 ActionMoveDown,
	},

	// Actions
	{
		Name:                   "Include File Without Comments",
		Description:            "Include the file without comments",
		Accelerator:            keys.Key("i"),
		AlternativeAccelerator: keys.Key("f1"),
		Keys:                   "i, f1",
		Group:                  GroupActions,
		Action:                 ActionIncludeFileWithoutComments,
	},
	{
		Name:                   "Include File With Comments",
		Description:            "Include the file with comments",
		Accelerator:            keys.Shift("i"),
		AlternativeAccelerator: keys.Shift("f1"),
		Keys:                   "shift+i, shift+f1",
		Group:                  GroupActions,
		Action:                 ActionIncludeFileWithComments,
	},
	{
		Name:                   "Include Component Without Comments",
		Description:            "Include the component without comments",
		Accelerator:            keys.Key("c"),
		AlternativeAccelerator: keys.CmdOrCtrl("f1"),
		Keys:                   "c, mod+f1",
		Group:                  GroupActions,
		Action:                 ActionIncludeComponentWithoutComments,
	},
	{
		Name:                   "Include Component With Comments",
		Description:            "Include the component with comments",
		Accelerator:            keys.Shift("c"),
		AlternativeAccelerator: keys.Combo("c", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:                   "shift+c, shift+mod+f1",
		Group:                  GroupActions,
		Action:                 ActionIncludeComponentWithComments,
	},
	{
		Name:                   "Dismiss File Without Comments",
		Description:            "Dismiss the file without comments",
		Accelerator:            keys.Key("d"),
		AlternativeAccelerator: keys.Key("f2"),
		Keys:                   "d, f2",
		Group:                  GroupActions,
		Action:                 ActionDismissFileWithoutComments,
	},
	{
		Name:                   "Dismiss File With Comments",
		Description:            "Dismiss the file with comments",
		Accelerator:            keys.Shift("d"),
		AlternativeAccelerator: keys.Shift("f2"),
		Keys:                   "shift+d, shift+f2",
		Group:                  GroupActions,
		Action:                 ActionDismissFileWithComments,
	},
	{
		Name:                   "Dismiss Component Without Comments",
		Description:            "Dismiss the component without comments",
		Accelerator:            keys.Key("x"),
		AlternativeAccelerator: keys.CmdOrCtrl("f2"),
		Keys:                   "x, mod+f2",
		Group:                  GroupActions,
		Action:                 ActionDismissComponentWithoutComments,
	},
	{
		Name:                   "Dismiss Component With Comments",
		Description:            "Dismiss the component with comments",
		Accelerator:            keys.Shift("x"),
		AlternativeAccelerator: keys.Combo("f2", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:                   "shift+x, shift+mod+f2",
		Group:                  GroupActions,
		Action:                 ActionDismissComponentWithComments,
	},
	{
		Name:                   "Replace File Without Comments",
		Description:            "Replace the file without comments",
		Accelerator:            keys.Key("r"),
		AlternativeAccelerator: keys.Key("f3"),
		Keys:                   "r, f3",
		Group:                  GroupActions,
		Action:                 ActionReplaceFileWithoutComments,
	},
	{
		Name:                   "Replace File With Comments",
		Description:            "Replace the file with comments",
		Accelerator:            keys.Shift("r"),
		AlternativeAccelerator: keys.Shift("f3"),
		Keys:                   "shift+r, shift+f3",
		Group:                  GroupActions,
		Action:                 ActionReplaceFileWithComments,
	},
	{
		Name:                   "Replace Component Without Comments",
		Description:            "Replace the component without comments",
		Accelerator:            keys.Key("e"),
		AlternativeAccelerator: keys.CmdOrCtrl("f3"),
		Keys:                   "e, mod+f3",
		Group:                  GroupActions,
		Action:                 ActionReplaceComponentWithoutComments,
	},
	{
		Name:                   "Replace Component With Comments",
		Description:            "Replace the component with comments",
		Accelerator:            keys.Shift("e"),
		AlternativeAccelerator: keys.Combo("f3", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:                   "shift+e, shift+mod+f3",
		Group:                  GroupActions,
		Action:                 ActionReplaceComponentWithComments,
	},

	{
		Name:        "Sync Scroll Position",
		Description: "Sync the scroll position of the editors",
		Accelerator: keys.Combo("e", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:        "shift+mod+e",
		Group:       GroupActions,
		Action:      ActionReplaceComponentWithComments,
	},

	// Scan
	{
		Name:        "Scan Current Directory",
		Description: "Scan the current directory",
		Accelerator: keys.Combo("b", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:        "shift+mod+b",
		Group:       GroupScan,
		Action:      ActionScanCwd,
	},
	{
		Name:        "Scan With Options",
		Description: "Run a scan with options",
		Accelerator: keys.Combo("c", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:        "shift+mod+c",
		Group:       GroupScan,
		Action:      ActionScanWithOptions,
	},
}
