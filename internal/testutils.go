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

package internal_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
	"github.com/stretchr/testify/mock"
)

func InitValidatorForTests() {
	v := validator.New()
	v.RegisterValidation("valid-purl", utils.ValidatePurl)
	utils.SetValidator(v)
}

func InitializeTestEnvironment(t *testing.T) func() {
	t.Helper()

	InitValidatorForTests()

	cfg := config.GetInstance()
	cfg.ScanRoot = t.TempDir()

	return func() {
		cfg = nil
	}
}

type MockUtils struct {
	mock.Mock
}

func NewMockUtils() *MockUtils { return &MockUtils{} }

func (m *MockUtils) ReadFile(filePath string) ([]byte, error) {
	args := m.Called(filePath)
	return args.Get(0).([]byte), args.Error(1)
}
