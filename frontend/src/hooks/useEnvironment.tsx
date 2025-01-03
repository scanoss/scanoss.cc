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

import { useEffect, useState } from 'react';

import { Environment, EnvironmentInfo } from '../../wailsjs/runtime/runtime';

interface UseEnvironmentReturnType {
  environment: EnvironmentInfo | undefined;
  isMac: boolean;
  modifierKey: {
    label: string;
    keyCode: string;
  };
}

export default function useEnvironment(): UseEnvironmentReturnType {
  const [environment, setEnvironment] = useState<EnvironmentInfo>();

  const isMac = environment?.platform === 'darwin';
  const modifierKey = isMac
    ? {
        label: '⌘',
        keyCode: 'Meta',
      }
    : {
        label: 'Ctrl',
        keyCode: 'Control',
      };

  useEffect(() => {
    async function fetchEnvironment() {
      const environment = await Environment();

      setEnvironment(environment);
    }

    fetchEnvironment();
  }, []);

  return {
    environment,
    isMac,
    modifierKey,
  };
}
