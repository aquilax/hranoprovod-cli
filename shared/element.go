package shared

import (
	"sort"
)

// Element contains single element data
type Element struct {
	Name string
	Val  float32
}

// Elements contains array of elements
type Elements []*Element

// NewElement creates new element
func NewElement(name string, val float32) *Element {
	return &Element{name, val}
}

// NewElements creates new element list
func NewElements() *Elements {
	return &Elements{}
}

// Add adds new element to the list
func (el *Elements) Add(name string, val float32) {
	*el = append(*el, NewElement(name, val))
}

// Index returns the index of the element by name
func (el *Elements) Index(name string) (int, bool) {
	for n, e := range *el {
		if e.Name == name {
			return n, true
		}
	}
	return 0, false
}

// SumMerge merges the element list with new one multiplied by coef
func (el *Elements) SumMerge(left *Elements, coef float32) {
	for _, v := range *left {
		if ndx, exists := (*el).Index(v.Name); exists {
			(*el)[ndx].Val += v.Val * coef
		} else {
			(*el).Add(v.Name, v.Val*coef)
		}
	}
}

// Len returns the length of the element list
func (el Elements) Len() int {
	return len(el)
}

// Less compares two elements
func (el Elements) Less(i, j int) bool {
	return el[i].Name < el[j].Name
}

// Swap swaps two elements
func (el Elements) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

// Sort sorts the element list
func (el Elements) Sort() {
	sort.Sort(el)
}
