// Mock declarations must be at the top
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import ScanDialog from '@/components/ScanDialog';

import { GetWorkingDir, SelectDirectory, SetScanRoot } from '../../wailsjs/go/main/App';
import { GetDefaultScanArgs, ScanStream } from '../../wailsjs/go/service/ScanServicePythonImpl';
import { EventsEmit, EventsOn } from '../../wailsjs/runtime/runtime';

vi.mock('../../wailsjs/go/main/App');
vi.mock('../../wailsjs/go/service/ScanServicePythonImpl');
vi.mock('../../wailsjs/runtime/runtime');
vi.mock('@/components/ui/use-toast');
vi.mock('@/components/ui/scroll-area');
vi.mock('@/modules/results/stores/useResultsStore');

const defaultArgs = ['-n', '100'];

// Setup mocks
const mockToast = vi.fn();
const mockFetchResults = vi.fn().mockResolvedValue(undefined);

// Configure mocks
vi.mocked(SelectDirectory).mockImplementation(vi.fn());
vi.mocked(SetScanRoot).mockImplementation(vi.fn());
vi.mocked(GetWorkingDir).mockResolvedValue('/default/path');
vi.mocked(GetDefaultScanArgs).mockResolvedValue(defaultArgs);
vi.mocked(ScanStream).mockImplementation(vi.fn());
vi.mocked(EventsOn).mockReturnValue(() => {});
vi.mocked(EventsEmit).mockImplementation(vi.fn());

vi.mock('@/components/ui/use-toast', () => ({
  useToast: () => ({
    toast: mockToast,
  }),
}));

vi.mock('@/components/ui/scroll-area', () => ({
  ScrollArea: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

vi.mock('@/modules/results/stores/useResultsStore', () => ({
  default: () => ({
    fetchResults: mockFetchResults,
  }),
}));

describe('ScanDialog Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('loads default scan arguments on mount', async () => {
    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    await waitFor(() => {
      expect(GetDefaultScanArgs).toHaveBeenCalled();
    });

    expect(EventsOn).toHaveBeenCalledWith('commandOutput', expect.any(Function));
    expect(EventsOn).toHaveBeenCalledWith('commandError', expect.any(Function));
    expect(EventsOn).toHaveBeenCalledWith('scanComplete', expect.any(Function));
    expect(EventsOn).toHaveBeenCalledWith('scanFailed', expect.any(Function));
  });

  it('allows directory selection', async () => {
    const selectedDir = '/selected/path';
    vi.mocked(SelectDirectory).mockResolvedValue(selectedDir);

    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    const selectButton = screen.getByRole('button', { name: /select directory/i });
    await act(async () => {
      fireEvent.click(selectButton);
    });

    expect(SelectDirectory).toHaveBeenCalled();

    await waitFor(() => {
      expect(screen.getByDisplayValue(selectedDir)).toBeInTheDocument();
    });
  });

  it('handles successful scan process', async () => {
    const selectedDir = '/test/path';
    vi.mocked(SelectDirectory).mockResolvedValue(selectedDir);
    vi.mocked(ScanStream).mockResolvedValue(undefined);

    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    // Select directory
    const selectButton = screen.getByRole('button', { name: /select directory/i });
    await act(async () => {
      fireEvent.click(selectButton);
    });

    // Wait for directory to be set
    await waitFor(() => {
      expect(screen.getByDisplayValue(selectedDir)).toBeInTheDocument();
    });

    // Start scan
    const scanButton = screen.getByRole('button', { name: /start scan/i });
    await act(async () => {
      fireEvent.click(scanButton);
    });

    // vi.spyOn(mockFetchResults, 'mockClear');

    await waitFor(() => {
      expect(ScanStream).toHaveBeenCalledWith([selectedDir, ...defaultArgs]);
      expect(SetScanRoot).toHaveBeenCalledWith(selectedDir);
    });
  });

  it('handles scan failure gracefully', async () => {
    const selectedDir = '/test/path';
    const errorMessage = 'Scan failed';
    vi.mocked(SelectDirectory).mockResolvedValue(selectedDir);
    vi.mocked(ScanStream).mockRejectedValue(new Error(errorMessage));

    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    // Select directory
    const selectButton = screen.getByRole('button', { name: /select directory/i });
    await act(async () => {
      fireEvent.click(selectButton);
    });

    // Wait for directory to be set
    await waitFor(() => {
      expect(screen.getByDisplayValue(selectedDir)).toBeInTheDocument();
    });

    // Start scan
    const scanButton = screen.getByRole('button', { name: /start scan/i });
    await act(async () => {
      fireEvent.click(scanButton);
    });

    await waitFor(() => {
      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Error',
          description: expect.stringContaining('error occurred while scanning'),
          variant: 'destructive',
        })
      );
    });
  });

  it('displays scan output correctly', async () => {
    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    // Get the event handler for commandOutput
    const commandOutputHandler = (EventsOn as jest.Mock).mock.calls.find((call) => call[0] === 'commandOutput')?.[1];

    const commandErrorHandler = (EventsOn as jest.Mock).mock.calls.find((call) => call[0] === 'commandError')?.[1];

    // Simulate output events
    await act(async () => {
      commandOutputHandler?.('Scanning file 1');
      commandOutputHandler?.('Scanning file 2');
      commandErrorHandler?.('Warning: file not found');
    });

    expect(screen.getByText(/scanning file 1/i)).toBeInTheDocument();
    expect(screen.getByText(/scanning file 2/i)).toBeInTheDocument();
    expect(screen.getByText(/warning: file not found/i)).toBeInTheDocument();
  });

  it('handles directory selection error', async () => {
    vi.mocked(SelectDirectory).mockRejectedValue(new Error('Failed to select directory'));

    render(<ScanDialog onOpenChange={() => {}} withOptions={true} />);

    const selectButton = screen.getByRole('button', { name: /select directory/i });
    await act(async () => {
      fireEvent.click(selectButton);
    });

    await waitFor(() => {
      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Error',
          description: expect.stringContaining('error occurred while selecting the directory'),
          variant: 'destructive',
        })
      );
    });
  });
});
