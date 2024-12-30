import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

import {
  GetResultFilePath,
  GetScanRoot,
  GetScanSettingsFilePath,
  SetResultFilePath,
  SetScanRoot,
  SetScanSettingsFilePath,
} from '../../wailsjs/go/main/App';

interface ResultsState {
  scanRoot: string;
  resultsFile: string;
  settingsFile: string;
}

interface ResultsActions {
  setScanRoot: (scanRoot: string) => Promise<void>;
  setResultsFile: (resultsFile: string) => Promise<void>;
  setSettingsFile: (settingsFile: string) => Promise<void>;

  getInitialConfig: () => Promise<void>;
}

type ConfigStore = ResultsState & ResultsActions;

export default create<ConfigStore>()(
  devtools((set) => ({
    scanRoot: '',
    resultsFile: '',
    settingsFile: '',

    setScanRoot: async (scanRoot: string) => {
      await SetScanRoot(scanRoot);
      set({ scanRoot });
    },
    setResultsFile: async (resultsFile: string) => {
      await SetResultFilePath(resultsFile);
      console.log(resultsFile);
      set({ resultsFile });
    },
    setSettingsFile: async (settingsFile: string) => {
      await SetScanSettingsFilePath(settingsFile);
      set({ settingsFile });
    },

    getInitialConfig: async () => {
      const scanRoot = await GetScanRoot();
      const resultsFile = await GetResultFilePath();
      const settingsFile = await GetScanSettingsFilePath();
      set({ scanRoot, resultsFile, settingsFile });
    },
  }))
);
