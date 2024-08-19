package application

import (
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/scan/domain"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ErrScanossPyNotDetected = errors.New("scanner python not detected\n")
	ErrCreatingTmpDirectory = errors.New("creating tmp directory\n")
	ErrCopyingFiles         = errors.New("copying files\n")
	ErrScanFailed           = errors.New("scanning failed\n")
	ErrStoringResult        = errors.New("storing result failed\n")
)

type ScanService struct{}

func NewScanService() *ScanService {
	return &ScanService{}
}

func (ss *ScanService) Scan(input domain.ScanInputPaths) (domain.ScanResult, error) {
	//Test scanoss-py exist in the system.
	cmd := exec.Command("scanoss-py", "scan")
	err := cmd.Run()
	if err != nil {
		return domain.ScanResult{}, ErrScanossPyNotDetected
	}

	tempDir, err := os.MkdirTemp("", "scanoss_git_*")
	if err != nil {
		return domain.ScanResult{}, ErrCreatingTmpDirectory
	}
	defer os.RemoveAll(tempDir)

	err = ss.CopyFilesToNewBase(input.GetRelativePaths(), input.GetBasePath(), tempDir)
	if err != nil {
		return domain.ScanResult{}, ErrCopyingFiles
	}

	cmd2 := ss.GenerateScanCommand(tempDir)
	scanResult, err := cmd2.Output()
	if err != nil {
		return domain.ScanResult{}, ErrScanFailed
	}

	err = os.WriteFile(config.Get().ResultFilePath, scanResult, 0777)
	if err != nil {
		return domain.ScanResult{}, ErrStoringResult
	}

	return domain.NewScanResult(scanResult), nil
}

func (ss *ScanService) GenerateScanCommand(folder string) *exec.Cmd {
	cmd := exec.Command("scanoss-py", "scan", "--no-wfp-output", folder)
	cmd.Dir = folder
	return cmd
}

// CopyFilesToNewBase copies a list of files to a new base directory, preserving their original directory structure
func (ss *ScanService) CopyFilesToNewBase(files []string, basePath string, newBasePath string) error {
	for _, file := range files {
		srcPath := filepath.Join(basePath, file)
		destPath := filepath.Join(newBasePath, file)

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories for %s: %w", destPath, err)
		}

		if err := ss.copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %w", srcPath, destPath, err)
		}
	}

	return nil
}

func (ss *ScanService) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
