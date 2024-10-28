package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/config"
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	"github.com/scanoss/scanoss.lui/backend/main/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFileRepositoryImpl(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	t.Run("ReadLocalFile", func(t *testing.T) {
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
	})

	t.Run("ReadRemoteFileByMD5", func(t *testing.T) {
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

		assert.NoError(t, err)
		assert.Equal(t, testContent, file.GetContent())

		repo.AssertExpectations(t)
	})
}