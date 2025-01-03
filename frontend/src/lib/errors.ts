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

/* eslint-disable @typescript-eslint/no-explicit-any */

interface WithErrorHandlingParams {
  asyncFn: (...args: any[]) => Promise<void>;
  onError: (error: unknown) => void;
  onFinish?: () => void;
  onSuccess?: () => Promise<void>;
}

export const withErrorHandling = ({ asyncFn, onError, onFinish, onSuccess }: WithErrorHandlingParams) => {
  return async (...args: any[]) => {
    try {
      await asyncFn(...args);
      await onSuccess?.();
    } catch (error) {
      onError(error);
    } finally {
      onFinish?.();
    }
  };
};
