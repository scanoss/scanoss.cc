package internal_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
	"github.com/stretchr/testify/mock"
)

func InitValidatorForTests() {
	v := validator.New()
	v.RegisterValidation("valid-purl", utils.ValidatePurl)
	utils.SetValidator(v)
}

func InitializeTestEnvironment(t *testing.T) func() {
	t.Helper()

	InitValidatorForTests()

	cfg := config.GetInstance()
	cfg.ScanRoot = t.TempDir()

	return func() {
		cfg = nil
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
