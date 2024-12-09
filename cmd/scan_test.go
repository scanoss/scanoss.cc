package cmd_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/service/mocks"
	"github.com/scanoss/scanoss.lui/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScanCommand(t *testing.T) {
	t.Run("successful scan with all flags", func(t *testing.T) {
		mockService := mocks.NewMockScanService(t)

		mockService.EXPECT().CheckDependencies().Return(nil)
		mockService.EXPECT().Scan("/test/path", mock.MatchedBy(func(args []string) bool {
			expectedArgs := []string{
				"--output", "results.json",
				"--no-wfp-output",
				"--threads", "10",
				"--format", "json",
			}
			// Sort both slices for comparison
			sort.Strings(args)
			sort.Strings(expectedArgs)
			return reflect.DeepEqual(args, expectedArgs)
		})).Return(nil)

		cmd := cmd.NewScanCmd(mockService)
		cmd.SetArgs([]string{
			"/test/path",
			"--output=results.json",
			"--no-wfp-output",
			"--threads=10",
			"--format=json",
		})

		err := cmd.Execute()

		assert.NoError(t, err)
	})
}
