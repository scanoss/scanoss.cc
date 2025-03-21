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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScanServicePythonImpl struct {
	cmd        string
	ctx        context.Context
	currentCmd *exec.Cmd
	cmdLock    sync.Mutex
	cancelFunc context.CancelFunc
}

func NewScanServicePythonImpl() *ScanServicePythonImpl {
	return &ScanServicePythonImpl{
		cmd:        "scanoss-py",
		currentCmd: nil,
		cancelFunc: nil,
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
	scanCtx, cancel := context.WithCancel(context.Background())

	s.cmdLock.Lock()
	s.cancelFunc = cancel
	s.cmdLock.Unlock()

	defer func() {
		cancel()
		s.cmdLock.Lock()
		s.cancelFunc = nil
		s.cmdLock.Unlock()
	}()

	cmd, stdout, stderr, err := s.executeScanWithPipes(args, scanCtx)
	if err != nil {
		return err
	}

	s.cmdLock.Lock()
	s.currentCmd = cmd
	s.cmdLock.Unlock()

	stdoutChan := make(chan string, 100)
	stderrChan := make(chan string, 100)
	done := make(chan error, 1)

	outputCtx, outputCancel := context.WithCancel(context.Background())
	defer outputCancel()

	go s.processStreamOutput(scanCtx, stdout, stdoutChan)
	go s.processStreamOutput(scanCtx, stderr, stderrChan)

	go s.emitOutputEvents(outputCtx, stdoutChan, "commandOutput")
	go s.emitOutputEvents(outputCtx, stderrChan, "commandError")

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		s.cmdLock.Lock()
		s.currentCmd = nil
		s.cmdLock.Unlock()

		// Wait a short time for any remaining output to be processed
		time.Sleep(100 * time.Millisecond)

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ExitCode() == -1 {
					s.emitEvent("scanAborted", "Scan was aborted")
					return fmt.Errorf("scan aborted: %w", exitErr)
				}
				s.emitEvent("scanFailed", exitErr.Error())
				return exitErr
			}
			return err
		}

		s.emitEvent("scanComplete", nil)
		s.emitEvent("commandOutput", "Scan completed successfully!")
		return nil

	case <-scanCtx.Done():
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		s.cmdLock.Lock()
		s.currentCmd = nil
		s.cmdLock.Unlock()

		s.emitEvent("scanAborted", "Scan was aborted")
		return nil
	}
}

func (s *ScanServicePythonImpl) AbortScan() error {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}

	if s.currentCmd == nil || s.currentCmd.Process == nil {
		return nil
	}

	err := s.currentCmd.Process.Signal(os.Interrupt)
	if err != nil {
		err := s.currentCmd.Process.Kill()
		if err != nil {
			log.Error().Err(err).Msg("Failed to abort scan")
			return fmt.Errorf("failed to abort scan: %w", err)
		}
	}

	s.emitEvent("scanAborted", nil)
	s.emitEvent("commandOutput", "Scan aborted by user")

	s.currentCmd = nil

	return nil
}

func (s *ScanServicePythonImpl) executeScanWithPipes(args []string, scanCtx context.Context) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	if err := s.CheckDependencies(); err != nil {
		s.emitEvent("scanFailed", err.Error())
		return nil, nil, nil, fmt.Errorf("dependency check failed: %w", err)
	}

	cmdArgs := []string{"scan"}

	defaultArgs, sensitiveArgs := s.GetDefaultScanArgs(), s.GetSensitiveDefaultScanArgs()

	if len(args) == 0 {
		cmdArgs = append(cmdArgs, ".") // scan current directory by default
		cmdArgs = append(cmdArgs, defaultArgs...)
	} else {
		cmdArgs = append(cmdArgs, args...)
	}

	cmdArgs = append(cmdArgs, sensitiveArgs...)

	// If the output folder does not exist, create it. This should be handled by the python cli
	s.maybeCreateOutputFolder(args)

	// This is to prevent we don't see anything on screen while scanning small directories
	env := os.Environ()
	env = append(env, "PYTHONUNBUFFERED=1")

	cmd := exec.CommandContext(scanCtx, s.cmd, cmdArgs...)
	cmd.Env = env

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

func (s *ScanServicePythonImpl) processStreamOutput(ctx context.Context, reader io.Reader, output chan<- string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			output <- scanner.Text()
		}
	}
}

func (s *ScanServicePythonImpl) emitOutputEvents(ctx context.Context, input <-chan string, eventName string) {
	for {
		select {
		case <-ctx.Done():
			return
		case text, ok := <-input:
			if !ok {
				return
			}
			s.emitEvent(eventName, text)
		}
	}
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
	pythonCommands := []string{"python3", "python"}

	for _, cmd := range pythonCommands {
		if err := exec.Command(cmd, "--version").Run(); err == nil {
			return nil
		}

		commonPaths := []string{
			"/usr/bin/",
			"/usr/local/bin/",
			"/opt/homebrew/bin/",
		}

		for _, path := range commonPaths {
			fullPath := path + cmd
			if _, err := os.Stat(fullPath); err == nil {
				if err := exec.Command(fullPath, "--version").Run(); err == nil {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("python is not installed or not found in PATH or common locations")
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

	if cfg.GetResultFilePath() != "" {
		args = append(args, "--output", cfg.GetResultFilePath())
	}

	if cfg.GetScanSettingsFilePath() != "" {
		args = append(args, "--settings", cfg.GetScanSettingsFilePath())
	}

	return args
}

func (s *ScanServicePythonImpl) GetSensitiveDefaultScanArgs() []string {
	args := make([]string, 0)
	cfg := config.GetInstance()

	if cfg.GetApiToken() != "" {
		args = append(args, "--key", cfg.GetApiToken())
	}

	if cfg.GetApiUrl() != "" {
		args = append(args, "--apiurl", fmt.Sprintf("%s/scan/direct", cfg.GetApiUrl()))
	}

	return args
}

func (s *ScanServicePythonImpl) emitEvent(eventName string, data ...any) {
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, eventName, data...)
	}
}

func (s *ScanServicePythonImpl) GetScanArgs() []entities.ScanArgDef {
	return entities.ScanArguments
}

func (s *ScanServicePythonImpl) maybeCreateOutputFolder(args []string) {
	outputPath := s.getOutputPathFromArgs(args)
	outputFolder := filepath.Dir(outputPath)
	if outputFolder != "" {
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			log.Info().Msgf("The provided output path does not exist. Creating it: %s", outputFolder)
			if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
				log.Error().Err(err).Msgf("Failed to create output folder: %s", outputFolder)
			}
		}
	}
}

func (s *ScanServicePythonImpl) getOutputPathFromArgs(args []string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--output" {
			return args[i+1]
		}
	}
	return ""
}
