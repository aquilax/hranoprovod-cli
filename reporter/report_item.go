package reporter

import (
	"time"

	"github.com/aquilax/hranoprovod-cli/accumulator"
	"github.com/aquilax/hranoprovod-cli/shared"
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

func getReportItem(ln *shared.LogNode, db shared.DBNodeList, options *Options) reportItem {
	var acc accumulator.Accumulator
	if options.Totals {
		acc = accumulator.NewAccumulator()
	}
	re := make([]reportElement, len(ln.Elements))
	//sort.Sort(ln.Elements)
	for i := range ln.Elements {
		re[i].Name = ln.Elements[i].Name
		re[i].Val = ln.Elements[i].Val
		if repl, found := db[re[i].Name]; found {
			for _, repl := range repl.Elements {
				res := repl.Val * re[i].Val
				re[i].Ingredients.Add(repl.Name, res)
				if options.Totals {
					acc.Add(repl.Name, res)
				}
			}
		} else {
			re[i].Ingredients.Add(re[i].Name, re[i].Val)
			if options.Totals {
				acc.Add(ln.Elements[i].Name, ln.Elements[i].Val)
			}
		}
	}
	var totals *[]total
	if options.TotalsOnly {
		re = nil
	}
	if options.Totals {
		totals = newTotalFromAccumulator(acc)
	}
	return reportItem{ln.Time, &re, totals}
}
