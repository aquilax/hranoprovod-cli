package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newBalanceCommand(t *testing.T) {
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
			mockCu := getMockCmdUtils([]string{tt.dbContent, tt.logContent}, New())
			callbackExecutedTimes := 0
			mockBalance := func(logStream, dbStream io.Reader, bc BalanceConfig) error {
				callbackExecutedTimes++
				return nil
			}

			a := getMockApp(newBalanceCommand(mockCu, mockBalance))

			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.executedTimes, callbackExecutedTimes)
		})
	}
}
