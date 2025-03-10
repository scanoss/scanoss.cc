// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
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

import { Settings } from 'lucide-react';
import { useEffect, useState } from 'react';

import { entities } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import SettingsDialog from './SettingsDialog';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';

export default function AppSettings() {
  const [showSettingsModal, setShowSettingsModal] = useState(false);

  const handleOpenSettingsModal = () => setShowSettingsModal(true);
  const handleCloseSettingsModal = () => setShowSettingsModal(false);

  useEffect(() => {
    EventsOn(entities.Action.OpenSettings, () => {
      setShowSettingsModal(true);
    });
  }, []);

  return (
    <>
      <Tooltip>
        <TooltipTrigger>
          <div className="cursor-pointer rounded-full p-1 hover:bg-muted" onClick={handleOpenSettingsModal}>
            <Settings className="h-4 w-4" />
          </div>
        </TooltipTrigger>
        <TooltipContent>App Settings</TooltipContent>
      </Tooltip>
      <SettingsDialog open={showSettingsModal} onOpenChange={handleCloseSettingsModal} />
    </>
  );
}
