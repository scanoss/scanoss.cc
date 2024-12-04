package utils

import (
	"context"
	"fmt"
	"os/exec"

	goRuntime "runtime"

	"github.com/scanoss/scanoss.lui/backend/entities"
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
