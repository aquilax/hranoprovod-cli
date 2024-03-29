package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/balance"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/csv"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/print"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/register"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/report"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/summary"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func Test_E2E(t *testing.T) {
	dbContent, err := testutils.ReadAsset("testAssets/food.yaml")
	assert.Equal(t, nil, err)
	logContent, err := testutils.ReadAsset("testAssets/log.yaml")
	assert.Equal(t, nil, err)

	balanceApp := func(w io.Writer) cli.App {
		mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent), string(logContent)}, w)
		return testutils.GetMockApp(balance.NewBalanceCommand(mockCu, balance.Balance))
	}

	tests := []struct {
		name                string
		args                []string
		appGetter           func(w io.Writer) cli.App
		wantError           error
		wantContentFileName string
	}{
		{
			"balance works without extra options",
			[]string{"bal"},
			balanceApp,
			nil,
			"testAssets/balance-no-extra-options.txt",
		},
		{
			"balance works collapse-last",
			[]string{"bal", "--collapse-last"},
			balanceApp,
			nil,
			"testAssets/balance-collapse-last.txt",
		},
		{
			"balance works with collapse",
			[]string{"bal", "--collapse"},
			balanceApp,
			nil,
			"testAssets/balance-collapse.txt",
		},
		{
			"balance works with single-element",
			[]string{"bal", "--single-element", "protein"},
			balanceApp,
			nil,
			"testAssets/balance-single-element.txt",
		},
		{
			"balance works with begin date",
			[]string{"bal", "-b", "2021/01/25"},
			balanceApp,
			nil,
			"testAssets/balance-begin-date.txt",
		},
		{
			"csv log works as expected",
			[]string{"log"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(logContent)}, w)
				return testutils.GetMockApp(csv.NewCSVLogCommand(mockCu, csv.CSVLog))
			},
			nil,
			"testAssets/csv-log.csv",
		},
		{
			"csv database works as expected",
			[]string{"database"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent)}, w)
				return testutils.GetMockApp(csv.NewCSVDatabaseCommand(mockCu, csv.CSVDatabase))
			},
			nil,
			"testAssets/csv-database.csv",
		},
		{
			"csv database-resolved works as expected",
			[]string{"database-resolved"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent)}, w)
				return testutils.GetMockApp(csv.NewCSVDatabaseResolvedCommand(mockCu, csv.CSVDatabaseResolved))
			},
			nil,
			"testAssets/csv-database-resolved.csv",
		},
		{
			"print works as expected with log file",
			[]string{"print"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(logContent)}, w)
				return testutils.GetMockApp(print.NewPrintCommand(mockCu, print.Print))
			},
			nil,
			"testAssets/print-log.yaml",
		},
		{
			"summary works as expected",
			[]string{"summary", "2021/01/24"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent), string(logContent)}, w)
				return testutils.GetMockApp(summary.NewSummaryCommand(mockCu, summary.Summary))
			},
			nil,
			"testAssets/summary.txt",
		},
		{
			"register works as expected",
			[]string{"register"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent), string(logContent)}, w)
				return testutils.GetMockApp(register.NewRegisterCommand(mockCu, register.Register))
			},
			nil,
			"testAssets/register.txt",
		},
		{
			"report totals works as expected",
			[]string{"totals"},
			func(w io.Writer) cli.App {
				mockCu := testutils.GetMockCmdUtilsRealOptions([]string{string(dbContent), string(logContent)}, w)
				return testutils.GetMockApp(report.NewReportTotalsCommand(mockCu, report.ReportTotals))
			},
			nil,
			"testAssets/report-total.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			a := tt.appGetter(w)
			err := a.Run(append(os.Args[:1], tt.args...))
			assert.Equal(t, tt.wantError, err)
			w.Flush()
			gotContent := buf.String()
			wantContent, err := testutils.ReadAsset(tt.wantContentFileName)
			assert.Equal(t, nil, err)
			if string(wantContent) != gotContent {
				println(tt.wantContentFileName)
				println(gotContent)
			}
			assert.Equal(t, string(wantContent), gotContent)
		})
	}
}
