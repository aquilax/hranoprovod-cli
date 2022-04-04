package print

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_newPrintCommand(t *testing.T) {
	mockError := errors.New("Mock error")
	tests := []struct {
		name        string
		args        []string
		lintError   error
		wantContent string
		wantError   error
	}{
		{
			"runs as expected",
			[]string{"print"},
			nil,
			"dummy",
			nil,
		},
		{
			"returns an error if the print command returns an error",
			[]string{"print"},
			mockError,
			"dummy",
			mockError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callbackExecuted := 0
			mockCu := testutils.GetMockCmdUtils([]string{tt.wantContent}, options.New())
			mockPrint := func(logStream io.Reader, pc PrintConfig) error {
				callbackExecuted++
				content, _ := io.ReadAll(logStream)
				assert.Equal(t, string(content), tt.wantContent)
				return tt.lintError
			}
			a := testutils.GetMockApp(NewPrintCommand(mockCu, mockPrint))

			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, 1, callbackExecuted)
		})
	}
}
