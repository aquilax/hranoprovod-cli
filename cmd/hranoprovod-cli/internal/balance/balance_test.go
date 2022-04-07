package balance

import (
	"io"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/options"
	"github.com/aquilax/hranoprovod-cli/v3/cmd/hranoprovod-cli/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_NewBalanceCommand(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantError     error
		dbContent     string
		logContent    string
		executedTimes int
		wantConfig    BalanceConfig
	}{
		{
			"works with empty input",
			[]string{"bal"},
			nil,
			"",
			"",
			1,
			BalanceConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCu := testutils.GetMockCmdUtils([]string{tt.dbContent, tt.logContent}, options.New())
			callbackExecutedTimes := 0
			mockBalance := func(logStream, dbStream io.Reader, bc BalanceConfig) error {
				callbackExecutedTimes++
				return nil
			}

			a := testutils.GetMockApp(NewBalanceCommand(mockCu, mockBalance))

			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.executedTimes, callbackExecutedTimes)
		})
	}
}
