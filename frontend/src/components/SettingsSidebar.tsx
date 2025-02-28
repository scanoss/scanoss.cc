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

import { SETTINGS_OPTIONS } from '@/lib/settings';

import { Button } from './ui/button';

interface SettingsSidebarProps {
  activeSetting: string;
  setActiveSetting: (setting: string) => void;
}

export default function SettingsSidebar({ activeSetting, setActiveSetting }: SettingsSidebarProps) {
  return (
    <div className="flex h-full flex-col">
      <div className="flex h-12 min-h-12 w-full items-center border-b px-4">
        <div className="flex items-center gap-2">
          <Settings className="h-4 w-4" />
          <span className="text-sm">SCANOSS Code Compare</span>
        </div>
      </div>
      <div className="flex flex-col gap-2 p-4">
        {SETTINGS_OPTIONS.map((option) => (
          <Button
            variant={activeSetting === option.id ? 'default' : 'outline'}
            size="sm"
            className="justify-start"
            onClick={() => setActiveSetting(option.id)}
            key={option.id}
          >
            <option.icon className="mr-2 h-4 w-4" />
            {option.title}
          </Button>
        ))}
      </div>
    </div>
  );
}
