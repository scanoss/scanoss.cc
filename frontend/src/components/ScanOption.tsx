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

import { Folder } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

interface ScanOptionProps {
  name: string;
  type: 'string' | 'int' | 'bool' | 'stringSlice';
  value: string | number | boolean | string[];
  defaultValue: string | number | boolean | string[];
  usage: string;
  onChange: (value: string | number | boolean | string[]) => void;
  onSelectFile?: () => void;
  isFileSelector?: boolean;
}

export default function ScanOption({ name, type, value, defaultValue, usage, onChange, onSelectFile, isFileSelector }: ScanOptionProps) {
  const displayName = name
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');

  const renderInput = () => {
    switch (type) {
      case 'bool':
        return (
          <div className="flex items-center space-x-2">
            <Checkbox id={name} checked={value as boolean} onCheckedChange={(checked) => onChange(checked || false)} />
            <Label htmlFor={name} className="cursor-help">
              {displayName}
            </Label>
          </div>
        );

      case 'int':
        return (
          <div className="space-y-2">
            <Label htmlFor={name} className="text-sm text-muted-foreground">
              {displayName}
            </Label>
            <Input
              id={name}
              type="number"
              value={value as number}
              onChange={(e) => {
                const val = parseInt(e.target.value);
                if (!isNaN(val)) {
                  onChange(val);
                }
              }}
              className="font-mono text-sm"
            />
          </div>
        );

      case 'stringSlice':
        return (
          <div className="space-y-2">
            <Label htmlFor={name}>{displayName}</Label>
            <div className="flex gap-2">
              <Input
                id={name}
                value={(value as string[]).join(',')}
                onChange={(e) => onChange(e.target.value.split(','))}
                placeholder={`${(defaultValue as string[]).join(',') || `[Enter ${displayName}]`}`}
                className="text-sm"
              />
              {isFileSelector && (
                <Button type="button" onClick={onSelectFile} variant="secondary" size="icon">
                  <Folder className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
        );

      case 'string':
      default:
        return (
          <div className="space-y-2">
            <Label htmlFor={name}>{displayName}</Label>
            <div className="flex gap-2">
              <Input
                id={name}
                value={value as string}
                onChange={(e) => onChange(e.target.value)}
                placeholder={`${(defaultValue as string) || `[Enter ${displayName}]`}`}
                className="text-sm"
              />
              {isFileSelector && (
                <Button type="button" onClick={onSelectFile} variant="secondary" size="icon">
                  <Folder className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
        );
    }
  };

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className="w-full">{renderInput()}</div>
      </TooltipTrigger>
      <TooltipContent>
        <p>{usage}</p>
      </TooltipContent>
    </Tooltip>
  );
}
