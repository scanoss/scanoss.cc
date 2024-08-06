package git

//
//import (
//	"fmt"
//	"log"
//	"os/exec"
//	"strings"
//)
//
//type GitStagingDirectoryRepositoryCLI struct {
//	rootPath string
//}
//
//func NewGitStagingDirectoryRepositoryCLI(rootPath string) (*GitStagingDirectoryRepositoryCLI, error) {
//	//TODO: Validate I can execute the git command
//	return &GitStagingDirectoryRepositoryCLI{
//		rootPath: rootPath,
//	}, nil
//}
//
//func (r *GitStagingDirectoryRepositoryCLI) Read() (*StagingDirectory, error) {
//	gsd := NewGitStagingDirectory()
//
//	cmd := exec.Command("git", "status", "--short")
//	output, err := cmd.Output()
//
//	if err != nil {
//		log.Fatalf("Error executing git command: %v", err)
//	}
//
//	return r.readFromBufferGitCommand(output)
//}
//
//func (r *GitStagingDirectoryRepositoryCLI) readFromBufferGitCommand(b []byte) (*StagingDirectory, error) {
//
//	outputStr := strings.TrimSpace(string(b))
//	lines := strings.Split(outputStr, "\n")
//
//	for _, line := range lines {
//		status := line[:2]
//		file := strings.TrimSpace(line[3:])
//		gsd.AddFile()
//		fmt.Printf("Status: %s, File: %s\n", status, file)
//	}
//
//}
