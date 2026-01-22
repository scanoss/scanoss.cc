// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package entities

import (
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// Action defines the type for keyboard actions
type Action string

const (
	ActionUndo        Action = "undo"
	ActionRedo        Action = "redo"
	ActionSave        Action = "save"
	ActionConfirm     Action = "confirm"
	ActionFocusSearch Action = "focusSearch"
	ActionSelectAll   Action = "selectAll"
	ActionMoveUp      Action = "moveUp"
	ActionMoveDown    Action = "moveDown"

	// Filter actions (BOM)
	ActionInclude          Action = "include"
	ActionIncludeWithModal Action = "includeWithModal"
	ActionIncludeFolder    Action = "includeFolder"
	ActionDismiss          Action = "dismiss"
	ActionDismissWithModal Action = "dismissWithModal"
	ActionDismissFolder    Action = "dismissFolder"
	ActionReplace          Action = "replace"
	ActionReplaceFolder    Action = "replaceFolder"
	ActionReplaceComponent Action = "replaceComponent"

	// Skip action (Scan settings) - always opens modal
	ActionSkip          Action = "skip"
	ActionSkipFolder    Action = "skipFolder"
	ActionSkipExtension Action = "skipExtension"

	// View
	ActionToggleSyncScrollPosition   Action = "toggleSyncScrollPosition"
	ActionShowKeyboardShortcutsModal Action = "showKeyboardShortcutsModal"
	ActionOpenSettings               Action = "openSettings"
	// Scan
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
	{ActionInclude, "Include"},
	{ActionIncludeWithModal, "IncludeWithModal"},
	{ActionIncludeFolder, "IncludeFolder"},
	{ActionDismiss, "Dismiss"},
	{ActionDismissWithModal, "DismissWithModal"},
	{ActionDismissFolder, "DismissFolder"},
	{ActionReplace, "Replace"},
	{ActionReplaceFolder, "ReplaceFolder"},
	{ActionReplaceComponent, "ReplaceComponent"},
	{ActionSkip, "Skip"},
	{ActionSkipFolder, "SkipFolder"},
	{ActionSkipExtension, "SkipExtension"},
	{ActionToggleSyncScrollPosition, "ToggleSyncScrollPosition"},
	{ActionShowKeyboardShortcutsModal, "ShowKeyboardShortcutsModal"},
	{ActionScanWithOptions, "ScanWithOptions"},
	{ActionOpenSettings, "OpenSettings"},
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
		Name:                   "Include",
		Description:            "Include file directly",
		Accelerator:            keys.Key("f1"),
		AlternativeAccelerator: keys.Key("i"),
		Keys:                   "i, f1",
		Group:                  GroupActions,
		Action:                 ActionInclude,
	},
	{
		Name:                   "Include (with options)",
		Description:            "Open include dialog for file/folder/component",
		Accelerator:            keys.Shift("f1"),
		AlternativeAccelerator: keys.Shift("i"),
		Keys:                   "shift+i, shift+f1",
		Group:                  GroupActions,
		Action:                 ActionIncludeWithModal,
	},
	{
		Name:        "Include folder",
		Description: "Open include dialog with folder selected",
		Accelerator: keys.Combo("i", keys.OptionOrAltKey, keys.ShiftKey),
		Keys:        "alt+shift+i",
		Group:       GroupActions,
		Action:      ActionIncludeFolder,
	},
	{
		Name:                   "Dismiss",
		Description:            "Dismiss file directly",
		Accelerator:            keys.Key("f2"),
		AlternativeAccelerator: keys.Key("d"),
		Keys:                   "d, f2",
		Group:                  GroupActions,
		Action:                 ActionDismiss,
	},
	{
		Name:                   "Dismiss (with options)",
		Description:            "Open dismiss dialog for file/folder/component",
		Accelerator:            keys.Shift("f2"),
		AlternativeAccelerator: keys.Shift("d"),
		Keys:                   "shift+d, shift+f2",
		Group:                  GroupActions,
		Action:                 ActionDismissWithModal,
	},
	{
		Name:        "Dismiss folder",
		Description: "Open dismiss dialog with folder selected",
		Accelerator: keys.Combo("d", keys.OptionOrAltKey, keys.ShiftKey),
		Keys:        "alt+shift+d",
		Group:       GroupActions,
		Action:      ActionDismissFolder,
	},
	{
		Name:                   "Replace",
		Description:            "Open replace dialog to select replacement component",
		Accelerator:            keys.Key("f3"),
		AlternativeAccelerator: keys.Key("r"),
		Keys:                   "r, f3",
		Group:                  GroupActions,
		Action:                 ActionReplace,
	},
	{
		Name:        "Replace folder",
		Description: "Open replace dialog with folder selected",
		Accelerator: keys.Combo("r", keys.OptionOrAltKey, keys.ShiftKey),
		Keys:        "alt+shift+r",
		Group:       GroupActions,
		Action:      ActionReplaceFolder,
	},
	{
		Name:                   "Skip",
		Description:            "Open skip dialog for file/folder/extension",
		Accelerator:            keys.Key("f4"),
		AlternativeAccelerator: keys.Key("s"),
		Keys:                   "s, f4",
		Group:                  GroupActions,
		Action:                 ActionSkip,
	},
	{
		Name:        "Skip folder",
		Description: "Open skip dialog with folder selected",
		Accelerator: keys.Combo("s", keys.OptionOrAltKey, keys.ShiftKey),
		Keys:        "alt+shift+s",
		Group:       GroupActions,
		Action:      ActionSkipFolder,
	},

	// View
	{
		Name:        "Sync Scroll Position",
		Description: "Sync the scroll position of the editors",
		Accelerator: keys.Combo("e", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:        "shift+mod+e",
		Group:       GroupView,
		Action:      ActionToggleSyncScrollPosition,
	},
	{
		Name:        "Settings",
		Description: "Open the app settings",
		Accelerator: keys.CmdOrCtrl(","),
		Keys:        "mod+,",
		Group:       GroupView,
		Action:      ActionOpenSettings,
	},

	// Scan
	{
		Name:        "Scan With Options",
		Description: "Run a scan with options",
		Accelerator: keys.Combo("c", keys.ShiftKey, keys.CmdOrCtrlKey),
		Keys:        "shift+mod+c",
		Group:       GroupScan,
		Action:      ActionScanWithOptions,
	},
}
