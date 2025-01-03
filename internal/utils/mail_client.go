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

package utils

import (
	"context"
	"fmt"
	"os/exec"

	goRuntime "runtime"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	SCANOSS_SUPPORT_MAILBOX = "support@scanoss.com"
)

func OpenMailClient(to string, subject string, body string) error {
	mailto := "mailto:" + to + "?subject=" + subject + "&body=" + body
	cmd := exec.Command("open", mailto)
	return cmd.Start()
}

func GetIssueReportBody(ctx context.Context) string {
	env := runtime.Environment(ctx)

	return fmt.Sprintf(`Please describe the issue you are facing in detail. 
Detailed information helps us to resolve the issue faster. Do not remove the following information: 

Arch: %s 
Platform: %s 
App Version: %s 
GO Version: %s`, env.Arch, env.Platform, entities.AppVersion, goRuntime.Version())
}
