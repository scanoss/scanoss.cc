package repository

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestReadLocalFile(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	testFilePath := "test.js"
	testContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
	currentPath := t.TempDir()
	absolutePath := filepath.Join(currentPath, testFilePath)
	err := os.WriteFile(absolutePath, testContent, 0644)
	assert.NoError(t, err)

	config.Get().ScanRoot = currentPath

	repo := NewFileRepositoryImpl()

	file, err := repo.ReadLocalFile(testFilePath)
	assert.NoError(t, err)
	assert.Equal(t, testContent, file.GetContent())
}

func TestReadLocalFile_FileNotFound(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	testFilePath := "nonexistent.js"
	currentPath := t.TempDir()

	config.Get().ScanRoot = currentPath

	repo := NewFileRepositoryImpl()

	_, err := repo.ReadLocalFile(testFilePath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, entities.ErrReadingFile))
}

func TestReadRemoteFileByMD5(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	testFilePath := "test.js"
	testContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
	md5 := "test-md5"

	repo := mocks.NewMockFileRepository(t)

	repo.EXPECT().ReadRemoteFileByMD5(testFilePath, md5).Return(*entities.NewFile(
		"",
		testFilePath,
		testContent,
	), nil)

	file, err := repo.ReadRemoteFileByMD5(testFilePath, md5)

	repo.AssertCalled(t, "ReadRemoteFileByMD5", testFilePath, md5)

	assert.NoError(t, err)
	assert.Equal(t, []byte(testContent), file.GetContent())
}
