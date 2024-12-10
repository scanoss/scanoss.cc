package service

import (
	"fmt"
	"os"
	"os/exec"
)

type ScanServicePythonImpl struct {
	cmd string
}

func NewScanServicePythonImpl() *ScanServicePythonImpl {
	return &ScanServicePythonImpl{
		cmd: "scanoss-py",
	}
}

func (s *ScanServicePythonImpl) Scan(args []string) error {
	if err := s.checkDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	cmdArgs := append([]string{"scan"}, args...)

	cmd := exec.Command(s.cmd, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s *ScanServicePythonImpl) checkDependencies() error {
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
