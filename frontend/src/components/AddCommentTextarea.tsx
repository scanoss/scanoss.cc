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

import { useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';

import { Textarea } from './ui/textarea';

interface AddCommentTextareaProps {
  onChange: (value: string) => void;
}

export default function AddCommentTextarea({ onChange }: AddCommentTextareaProps) {
  const { modifierKey } = useEnvironment();
  const [inputValue, setInputValue] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputValue(e.target.value);
    onChange(e.target.value);
  };

  return (
    <div>
      <Textarea value={inputValue} onChange={handleChange} />
      <p className="mt-2 text-xs text-muted-foreground">
        Use <span className="rounded bg-primary px-1.5">{modifierKey.label} + return</span> to confirm
      </p>
    </div>
  );
}
