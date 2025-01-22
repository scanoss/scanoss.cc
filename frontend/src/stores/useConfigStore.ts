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
