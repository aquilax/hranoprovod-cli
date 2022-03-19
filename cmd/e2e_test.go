package cmd

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/stretchr/testify/assert"
)

func Test_E2E(t *testing.T) {
	dbContent, err := ioutil.ReadFile("../examples/food.yaml")
	assert.Equal(t, nil, err)
	logContent, err := ioutil.ReadFile("../examples/log.yaml")
	assert.Equal(t, nil, err)

	tests := []struct {
		name                string
		args                []string
		wantError           error
		dbContent           string
		logContent          string
		wantContentFileName string
	}{
		{
			"balance works without extra options",
			[]string{"bal"},
			nil,
			string(dbContent),
			string(logContent),
			`testAssets/balance-no-extra-options.txt`,
		},
		{
			"balance works collapse-last",
			[]string{"bal", "--collapse-last"},
			nil,
			string(dbContent),
			string(logContent),
			`testAssets/balance-collapse-last.txt`,
		},
		{
			"balance works with collapse",
			[]string{"bal", "--collapse"},
			nil,
			string(dbContent),
			string(logContent),
			`testAssets/balance-collapse.txt`,
		},
		{
			"balance works with single-element",
			[]string{"bal", "--single-element", "protein"},
			nil,
			string(dbContent),
			string(logContent),
			`testAssets/balance-single-element.txt`,
		},
		{
			"balance works with begin date",
			[]string{"bal", "-b", "2021/01/25"},
			nil,
			string(dbContent),
			string(logContent),
			`testAssets/balance-begin-date.txt`,
		},
	}

	updateSnapshots := os.Getenv("UPDATE_SNAPSHOTS") == "1"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			mockCu := getMockCmdUtilsRealOptions([]string{tt.dbContent, tt.logContent}, w)
			a := getMockApp(newBalanceCommand(mockCu, app.Balance))
			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			w.Flush()
			gotContent := buf.String()
			if updateSnapshots {
				err := ioutil.WriteFile(tt.wantContentFileName, buf.Bytes(), 0644)
				assert.Equal(t, nil, err)
			}
			wantContent, err := ioutil.ReadFile(tt.wantContentFileName)
			assert.Equal(t, nil, err)
			assert.Equal(t, string(wantContent), gotContent)
		})
	}
}
