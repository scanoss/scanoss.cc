package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/scanoss/scanoss.lui/internal/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScanServicePythonImpl struct {
	cmd string
	ctx context.Context
}

func NewScanServicePythonImpl() *ScanServicePythonImpl {
	return &ScanServicePythonImpl{
		cmd: "scanoss-py",
	}
}

func (s *ScanServicePythonImpl) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ScanServicePythonImpl) Scan(args []string) error {
	cmdArgs := append([]string{"scan"}, args...)

	cmd := exec.Command(s.cmd, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s *ScanServicePythonImpl) ScanStream(args []string) error {
	cmd, stdout, stderr, err := s.executeScanWithPipes(args)
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			s.emitEvent("commandOutput", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s.emitEvent("commandError", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			s.emitEvent("scanFailed", exitErr.Error())
			return exitErr
		}
		return err
	}

	s.emitEvent("scanComplete", nil)
	s.emitEvent("commandOutput", "Scan completed succesfully!")
	return nil
}

func (s *ScanServicePythonImpl) executeScanWithPipes(args []string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	if err := s.CheckDependencies(); err != nil {
		s.emitEvent("scanFailed", err.Error())
		return nil, nil, nil, fmt.Errorf("dependency check failed: %w", err)
	}

	cmdArgs := []string{"scan"}

	defaultArgs, sensitiveArgs := s.GetDefaultScanArgs(), s.GetSensitiveDefaultScanArgs()

	cmdArgs = append(cmdArgs, sensitiveArgs...)

	if len(args) == 0 {
		cmdArgs = append(cmdArgs, defaultArgs...)
	} else {
		cmdArgs = append(cmdArgs, args...)
	}

	cmd := exec.Command(s.cmd, cmdArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	return cmd, stdout, stderr, nil
}

func (s *ScanServicePythonImpl) CheckDependencies() error {
	if err := s.checkPythonInstalled(); err != nil {
		return err
	}

	if err := s.checkScanossPyInstalled(); err != nil {
		return err
	}

	return nil
}

func (s *ScanServicePythonImpl) checkPythonInstalled() error {
	cmd := exec.Command("python", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("python is not installed: %w", err)
	}
	return nil
}

func (s *ScanServicePythonImpl) checkScanossPyInstalled() error {
	cmd := exec.Command(s.cmd, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scanoss-py is not installed or not in PATH: %w", err)
	}
	return nil
}

func (s *ScanServicePythonImpl) GetDefaultScanArgs() []string {
	args := []string{}
	cfg := config.GetInstance()

	if cfg.ResultFilePath != "" {
		relativePath, err := utils.GetRelativePath(cfg.ResultFilePath)
		if err != nil {
			return nil
		}
		args = append(args, "--output", relativePath)
	}

	if cfg.ScanSettingsFilePath != "" {
		relativePath, err := utils.GetRelativePath(cfg.ScanSettingsFilePath)
		if err != nil {
			return nil
		}
		args = append(args, "--settings", relativePath)
	}

	return args
}

func (s *ScanServicePythonImpl) GetSensitiveDefaultScanArgs() []string {
	args := s.GetDefaultScanArgs()
	cfg := config.GetInstance()

	if cfg.ApiToken != "" {
		args = append(args, "--key", cfg.ApiToken)
	}

	if cfg.ApiUrl != "" {
		args = append(args, "--apiurl", fmt.Sprintf("%s/scan/direct", cfg.ApiUrl))
	}

	return args
}

func (s *ScanServicePythonImpl) emitEvent(eventName string, data ...interface{}) {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, eventName, data...)
	}
}
