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

import useResultsStore from '@/modules/results/stores/useResultsStore';

import { entities } from '../../../../wailsjs/go/models';
import { CanRedo, CanUndo, FilterComponents, Redo, Undo } from '../../../../wailsjs/go/service/ComponentServiceImpl';
import { FilterAction } from '../domain';

interface ComponentFilterState {
  canRedo: boolean;
  canUndo: boolean;
}

export interface OnFilterComponentArgs {
  action: FilterAction;
  filterBy: 'by_file' | 'by_purl';
  withComment?: boolean;
  comment?: string;
  license?: string;
  replaceWith?: string;
}

interface ComponentFilterActions {
  onFilterComponent: (args: OnFilterComponentArgs) => Promise<void>;
  redo: () => Promise<void>;
  undo: () => Promise<void>;
  updateUndoRedoState: () => Promise<void>;
}

type ComponentFilterStore = ComponentFilterState & ComponentFilterActions;

const useComponentFilterStore = create<ComponentFilterStore>()(
  devtools((set, get) => ({
    canUndo: false,
    canRedo: false,

    onFilterComponent: async (args: OnFilterComponentArgs) => {
      const { replaceWith, filterBy, action, comment, license } = args;

      if (action === FilterAction.Replace && !replaceWith) {
        throw new Error('There was an error replacing the component. Please try again.');
      }

      const selectedResults = useResultsStore.getState().selectedResults;

      const dto: entities.ComponentFilterDTO[] = selectedResults.map((result) => ({
        action,
        comment,
        license,
        purl: result.detected_purl ?? '',
        ...(filterBy === 'by_file' && { path: result.path }),
        ...(replaceWith && {
          replace_with: replaceWith,
        }),
      }));

      useResultsStore.getState().moveToNextResult();

      await FilterComponents(dto);

      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();
    },

    undo: async () => {
      await Undo();
      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();
    },

    redo: async () => {
      await Redo();
      await get().updateUndoRedoState();
      await useResultsStore.getState().fetchResults();
    },

    updateUndoRedoState: async () => {
      const [canUndo, canRedo] = await Promise.all([CanUndo(), CanRedo()]);
      set({ canUndo, canRedo }, false, 'UPDATE_UNDO_REDO_STATE');
    },
  }))
);

export default useComponentFilterStore;
