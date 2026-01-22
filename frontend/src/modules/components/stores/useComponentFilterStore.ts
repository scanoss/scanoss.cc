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
  filterBy: 'by_file' | 'by_purl' | 'by_folder';
  withComment?: boolean;
  comment?: string;
  license?: string;
  replaceWith?: string;
  folderPath?: string;
  purl?: string;
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
      const { replaceWith, filterBy, action, comment, license, folderPath, purl } = args;

      if (action === FilterAction.Replace && !replaceWith) {
        throw new Error('There was an error replacing the component. Please try again.');
      }

      const selectedResults = useResultsStore.getState().selectedResults;

      // Handle folder-level filtering
      if (filterBy === 'by_folder' && folderPath) {
        // Ensure folder path ends with /
        const normalizedFolderPath = folderPath.endsWith('/') ? folderPath : folderPath + '/';

        // Get all unique purls for files in this folder (from both pending and completed results)
        const { pendingResults, completedResults } = useResultsStore.getState();
        const allResults = [...pendingResults, ...completedResults];
        const purlsInFolder = [...new Set(
          allResults
            .filter(r => r.path.startsWith(normalizedFolderPath))
            .map(r => r.detected_purl)
            .filter((p): p is string => Boolean(p))
        )];

        // Create one DTO per unique purl
        const dto: entities.ComponentFilterDTO[] = purlsInFolder.map(folderPurl => ({
          action,
          comment,
          license,
          purl: folderPurl,
          path: normalizedFolderPath,
          ...(replaceWith && { replace_with: replaceWith }),
        }));

        useResultsStore.getState().moveToNextResult();
        await FilterComponents(dto);
        await get().updateUndoRedoState();
        return;
      }

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
    },

    undo: async () => {
      await Undo();
      await get().updateUndoRedoState();
    },

    redo: async () => {
      await Redo();
      await get().updateUndoRedoState();
    },

    updateUndoRedoState: async () => {
      const [canUndo, canRedo] = await Promise.all([CanUndo(), CanRedo()]);
      set({ canUndo, canRedo }, false, 'UPDATE_UNDO_REDO_STATE');
    },
  }))
);

export default useComponentFilterStore;
