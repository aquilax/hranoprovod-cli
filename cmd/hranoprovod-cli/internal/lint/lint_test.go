package lint

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_newLintCommand(t *testing.T) {
	mockError := errors.New("Mock error")
	tests := []struct {
		name        string
		args        []string
		lintError   error
		wantContent string
		wantSilent  bool
		wantError   error
	}{
		{
			"runs as expected",
			[]string{"lint", "mock.yaml"},
			nil,
			"dummy",
			false,
			nil,
		},
		{
			"runs silently",
			[]string{"lint", "--silent", "mock.yaml"},
			nil,
			"dummy",
			true,
			nil,
		},
		{
			"returns an error if the linter returns an error",
			[]string{"lint", "mock.yaml"},
			mockError,
			"dummy",
			false,
			mockError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callbackExecuted := 0
			mockCu := testutils.GetMockCmdUtils([]string{tt.wantContent}, options.New())
			mockLint := func(stream io.Reader, lc LintConfig) error {
				callbackExecuted++
				content, _ := io.ReadAll(stream)
				assert.Equal(t, string(content), tt.wantContent)
				assert.Equal(t, lc.Silent, tt.wantSilent)
				return tt.lintError
			}
			a := testutils.GetMockApp(newLintCommand(mockCu, mockLint))

			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, 1, callbackExecuted)
		})
	}
}
