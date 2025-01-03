//go:build unit

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

package cmd_test

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/service/mocks"
	"github.com/scanoss/scanoss.cc/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScanCommand(t *testing.T) {
	t.Run("successful scan folder with flags", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)

		mockService.EXPECT().CheckDependencies().Return(nil)
		// We need to sort the arguments to match the ones we set below
		mockService.EXPECT().Scan(mock.MatchedBy(func(args []string) bool {
			expectedArgs := []string{
				"/test/path",
				"--quiet",
				"--output", "results.json",
				"--no-wfp-output",
				"--threads", "10",
				"--format", "json",
			}
			sort.Strings(args)
			sort.Strings(expectedArgs)
			return reflect.DeepEqual(args, expectedArgs)
		})).Return(nil)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"/test/path",
			"--output=results.json",
			"--no-wfp-output",
			"--threads=10",
			"--format=json",
		})

		err := cmd.Execute()

		assert.NoError(t, err)
	})

	t.Run("successful scan with only scan path", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)
		mockService.EXPECT().CheckDependencies().Return(nil)
		mockService.EXPECT().Scan(mock.MatchedBy(func(args []string) bool {
			expectedArgs := []string{
				"/test/path",
				"--quiet",
			}
			sort.Strings(args)
			sort.Strings(expectedArgs)
			return reflect.DeepEqual(args, expectedArgs)
		})).Return(nil)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{"/test/path"})

		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("successful scan with file list", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)
		mockService.EXPECT().CheckDependencies().Return(nil)
		mockService.EXPECT().Scan(mock.MatchedBy(func(args []string) bool {
			expectedArgs := []string{
				"--files", "file1.go,file2.go",
				"--quiet",
			}
			sort.Strings(args)
			sort.Strings(expectedArgs)
			return reflect.DeepEqual(args, expectedArgs)
		})).Return(nil)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"--files", "file1.go,file2.go",
		})

		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("fails with invalid flag value", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"/test/path",
			"--threads", "invalid",
		})

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid argument")
	})

	t.Run("handles multiple boolean flags", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)
		mockService.EXPECT().CheckDependencies().Return(nil)
		mockService.EXPECT().Scan(mock.MatchedBy(func(args []string) bool {
			expectedArgs := []string{
				"/test/path",
				"--quiet",
				"--no-wfp-output",
				"--dependencies",
				"--debug",
				"--trace",
			}
			sort.Strings(args)
			sort.Strings(expectedArgs)
			return reflect.DeepEqual(args, expectedArgs)
		})).Return(nil)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"/test/path",
			"--no-wfp-output",
			"--dependencies",
			"--debug",
			"--trace",
		})

		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("fails when dependency check fails", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)
		mockService.EXPECT().CheckDependencies().Return(
			errors.New("dependency check failed"),
		)
		mockService.EXPECT().Scan(mock.Anything).Maybe()

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"/test/path",
		})

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dependency check failed")
	})
}
