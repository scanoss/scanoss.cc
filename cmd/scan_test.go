//go:build unit

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
