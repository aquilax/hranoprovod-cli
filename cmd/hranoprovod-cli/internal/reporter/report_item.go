package reporter

import (
	"sort"
	"time"

	shared "github.com/aquilax/hranoprovod-cli/v3"
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

func newTotalFromAccumulator(acc shared.Accumulator) *[]total {
	var result = make([]total, len(acc))
	var ss = make(sort.StringSlice, len(acc))
	i := 0
	for name := range acc {
		ss[i] = name
		i++
	}
	sort.Sort(ss)
	for i, name := range ss {
		result[i] = total{name, acc[name][shared.Positive], acc[name][shared.Negative], acc[name][shared.Positive] + acc[name][shared.Negative]}
	}
	return &result
}

func GetReportItem(ln *shared.LogNode, db shared.DBNodeMap, config Config) reportItem {
	var acc shared.Accumulator
	if config.Totals {
		acc = shared.NewAccumulator()
	}

	re := make([]reportElement, len(ln.Elements))
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
