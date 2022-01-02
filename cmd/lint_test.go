package cmd

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
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
			mockCu := getMockCu(tt.wantContent)

			mockLint := func(stream io.Reader, silent bool, pc parser.Config, rpc reporter.Config) error {
				callbackExecuted++
				content, _ := io.ReadAll(stream)
				assert.Equal(t, string(content), tt.wantContent)
				assert.Equal(t, silent, tt.wantSilent)
				return tt.lintError
			}
			a := getMockApp(newLintCommand(mockCu, mockLint))

			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, 1, callbackExecuted)
		})
	}
}
