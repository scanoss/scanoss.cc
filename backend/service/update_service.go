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
	"context"

	"github.com/scanoss/scanoss.cc/backend/entities"
)

// UpdateService handles checking for and applying application updates
type UpdateService interface {
	// CheckForUpdate checks if a new version is available
	CheckForUpdate() (*entities.UpdateInfo, error)

	// DownloadUpdate downloads the new version to a temporary location
	DownloadUpdate(updateInfo *entities.UpdateInfo) (string, error)

	// ApplyUpdate applies the downloaded update and restarts the application
	ApplyUpdate(updatePath string) error

	// GetCurrentVersion returns the current application version
	GetCurrentVersion() string

	// SetContext sets the context for the service
	SetContext(ctx context.Context)

	// VerifyUpdateSuccess checks if an update completed successfully and cleans up backup
	VerifyUpdateSuccess() error

	// CheckForFailedUpdate checks if the previous update failed and performs rollback if needed
	CheckForFailedUpdate() error
}
