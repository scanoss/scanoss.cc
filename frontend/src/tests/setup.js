import '@testing-library/jest-dom/vitest';

import { cleanup } from '@testing-library/react';
import { afterEach, vi } from 'vitest';

import { runtimeMocks } from './__mocks__/wailsjs/runtime';
import { serviceMocks } from './__mocks__/wailsjs/services';

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

vi.mock('../../wailsjs/runtime/runtime', () => runtimeMocks);
vi.mock('../../wailsjs/go/service/ScanServicePythonImpl', () => serviceMocks.ScanServicePythonImpl);
