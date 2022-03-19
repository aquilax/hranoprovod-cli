package reporter

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/accumulator"
	"github.com/aquilax/hranoprovod-cli/v2/shared"
)

type reportElement struct {
	shared.Element
	Ingredients shared.Elements
}

type total struct {
	Name     string
	Positive float64
	Negative float64
	Sum      float64
}

type reportItem struct {
	Time     time.Time
	Elements *[]reportElement
	Totals   *[]total
}

func getReportItem(ln *shared.LogNode, db shared.DBNodeMap, config Config) reportItem {
	var acc accumulator.Accumulator
	if config.Totals {
		acc = accumulator.NewAccumulator()
	}
	re := make([]reportElement, len(ln.Elements))
	//sort.Sort(ln.Elements)
	for i := range ln.Elements {
		re[i].Name = ln.Elements[i].Name
		re[i].Value = ln.Elements[i].Value
		if repl, found := db[re[i].Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Value * re[i].Value
				re[i].Ingredients.Add(repl.Name, res)
				if config.Totals {
					acc.Add(repl.Name, res)
				}
			}
		} else {
			re[i].Ingredients.Add(re[i].Name, re[i].Value)
			if config.Totals {
				acc.Add(ln.Elements[i].Name, ln.Elements[i].Value)
			}
		}
	}
	var totals *[]total
	if config.TotalsOnly {
		re = nil
	}
	if config.Totals {
		totals = newTotalFromAccumulator(acc)
	}
	return reportItem{ln.Time, &re, totals}
}
