import { vi } from 'vitest';

export const serviceMocks = {
  ScanServicePythonImpl: {
    GetDefaultScanArgs: vi.fn().mockResolvedValue([]),
    GetSensitiveDefaultScanArgs: vi.fn().mockResolvedValue([]),
    Scan: vi.fn().mockResolvedValue(undefined),
    ScanStream: vi.fn().mockResolvedValue(undefined),
    CheckDependencies: vi.fn().mockResolvedValue(undefined),
    SetContext: vi.fn(),
  },
};
