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
