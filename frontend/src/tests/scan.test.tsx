import { render, screen } from '@testing-library/react';
import { beforeEach, describe, it, vi } from 'vitest';

import ScanDialog from '@/components/ScanDialog';

describe('ScanDialog', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the ScanDialog component', () => {
    render(<ScanDialog onOpenChange={() => {}} withOptions />);
    screen.debug();
  });
});
