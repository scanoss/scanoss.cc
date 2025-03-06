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

package service

import (
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
)

type ScanossSettingsServiceImp struct {
	repository repository.ScanossSettingsRepository
}

func NewScanossSettingsServiceImpl(r repository.ScanossSettingsRepository) *ScanossSettingsServiceImp {
	return &ScanossSettingsServiceImp{
		repository: r,
	}
}

func (s *ScanossSettingsServiceImp) Save() error {
	return s.repository.Save()
}

func (s *ScanossSettingsServiceImp) HasUnsavedChanges() (bool, error) {
	return s.repository.HasUnsavedChanges()
}

func (s *ScanossSettingsServiceImp) GetSettings() *entities.SettingsFile {
	return s.repository.GetSettings()
}

func (s *ScanossSettingsServiceImp) AddStagedScanningSkipPattern(pattern string) error {
	return s.repository.AddStagedScanningSkipPattern(pattern)
}

func (s *ScanossSettingsServiceImp) RemoveStagedScanningSkipPattern(pattern string) error {
	return s.repository.RemoveStagedScanningSkipPattern(pattern)
}

func (s *ScanossSettingsServiceImp) CommitStagedSkipPatterns() error {
	return s.repository.CommitStagedSkipPatterns()
}

func (s *ScanossSettingsServiceImp) DiscardStagedSkipPatterns() error {
	return s.repository.DiscardStagedSkipPatterns()
}

func (s *ScanossSettingsServiceImp) HasStagedChanges() bool {
	return s.repository.HasStagedChanges()
}
