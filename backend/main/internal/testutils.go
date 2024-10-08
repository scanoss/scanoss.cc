package internal_test

import (
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/stretchr/testify/mock"
)

func InitializeTestEnvironment(t *testing.T) func() {
	t.Helper()

	cfgPath := t.TempDir() + "/config.json"
	configModule := config.NewConfigModule(cfgPath)

	err := configModule.Init()
	if err != nil {
		t.Fatalf("Error initializing config: %s", err.Error())
	}

	err = configModule.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	configModule.Config.ScanRoot = t.TempDir()

	return func() {
		configModule = nil
	}
}

type MockUtils struct {
	mock.Mock
}

func NewMockUtils() *MockUtils { return &MockUtils{} }

func (m *MockUtils) ReadFile(filePath string) ([]byte, error) {
	args := m.Called(filePath)
	return args.Get(0).([]byte), args.Error(1)
}
