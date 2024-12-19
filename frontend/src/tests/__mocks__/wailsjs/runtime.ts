// src/tests/__mocks__/wailsjs/runtime.ts
import { vi } from 'vitest';

export const runtimeMocks = {
  // Event handling
  EventsOn: vi.fn().mockImplementation(() => {
    // Return a cleanup function
    return vi.fn();
  }),
  EventsOnMultiple: vi.fn().mockImplementation(() => {
    // Return a cleanup function
    return vi.fn();
  }),
  EventsOff: vi.fn(),
  EventsOnce: vi.fn().mockImplementation(() => {
    // Return a cleanup function
    return vi.fn();
  }),
  EventsEmit: vi.fn(),

  // Logging functions
  LogPrint: vi.fn(),
  LogTrace: vi.fn(),
  LogDebug: vi.fn(),
  LogInfo: vi.fn(),
  LogWarning: vi.fn(),
  LogError: vi.fn(),
  LogFatal: vi.fn(),

  // Window management
  WindowReload: vi.fn(),
  WindowReloadApp: vi.fn(),
  WindowSetAlwaysOnTop: vi.fn(),
  WindowSetSystemDefaultTheme: vi.fn(),
  WindowSetLightTheme: vi.fn(),
  WindowSetDarkTheme: vi.fn(),
  WindowCenter: vi.fn(),
  WindowSetTitle: vi.fn(),
  WindowFullscreen: vi.fn(),
  WindowUnfullscreen: vi.fn(),
  WindowIsFullscreen: vi.fn(),
  WindowGetSize: vi.fn(),
  WindowSetSize: vi.fn(),
  WindowSetMaxSize: vi.fn(),
  WindowSetMinSize: vi.fn(),
  WindowSetPosition: vi.fn(),
  WindowGetPosition: vi.fn(),
  WindowHide: vi.fn(),
  WindowShow: vi.fn(),
  WindowMaximise: vi.fn(),
  WindowToggleMaximise: vi.fn(),
  WindowUnmaximise: vi.fn(),
  WindowIsMaximised: vi.fn(),
  WindowMinimise: vi.fn(),
  WindowUnminimise: vi.fn(),
  WindowSetBackgroundColour: vi.fn(),
  WindowIsMinimised: vi.fn(),
  WindowIsNormal: vi.fn(),

  // Screen
  ScreenGetAll: vi.fn(),

  // Browser
  BrowserOpenURL: vi.fn(),

  // Environment and System
  Environment: vi.fn(),
  Quit: vi.fn(),
  Hide: vi.fn(),
  Show: vi.fn(),

  // Clipboard
  ClipboardGetText: vi.fn(),
  ClipboardSetText: vi.fn(),

  // File Operations
  OnFileDrop: vi.fn(),
  OnFileDropOff: vi.fn(),
  CanResolveFilePaths: vi.fn(),
  ResolveFilePaths: vi.fn(),
};
