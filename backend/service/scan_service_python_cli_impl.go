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

func (s *ScanServicePythonImpl) CheckDependencies() error {
	if err := checkPythonInstalled(); err != nil {
		return err
	}

	if err := s.checkScanossPyInstalled(); err != nil {
		return err
	}

	return nil
}

func checkPythonInstalled() error {
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

func (s *ScanServicePythonImpl) Scan(dirPath string, args []string) error {
	cmdArgs := append([]string{"scan", dirPath}, args...)

	cmd := exec.Command(s.cmd, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
