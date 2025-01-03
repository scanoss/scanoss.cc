// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
