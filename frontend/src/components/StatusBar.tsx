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

import { useEffect } from 'react';

import useConfigStore from '@/stores/useConfigStore';

import AppSettings from './AppSettings';
import SelectResultsFile from './SelectResultsFile';
import SelectScanRoot from './SelectScanRoot';
import SelectSettingsFile from './SelectSettingsFile';
import UpdateNotification from './UpdateNotification';

export default function StatusBar() {
  const getInitialConfig = useConfigStore((state) => state.getInitialConfig);

  useEffect(() => {
    getInitialConfig();
  }, []);

  return (
    <div className="flex w-full justify-between bg-background px-4 py-1 text-xs text-muted-foreground">
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2">
          <span>Scan Root:</span>
          <SelectScanRoot />
        </div>
        <div className="flex items-center gap-2">
          <span>Results File:</span>
          <SelectResultsFile />
        </div>
        <div className="flex items-center gap-2">
          <span>Settings File:</span>
          <SelectSettingsFile />
        </div>
      </div>
      <div className="flex items-center gap-2">
        <UpdateNotification />
        <AppSettings />
      </div>
    </div>
  );
}
