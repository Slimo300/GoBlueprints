package meander

import (
	"errors"
	"strings"
)

// Cost is an enumerator type
type Cost int8

const (
	_ Cost = iota
	// Cost1 = 1
	Cost1
	// Cost2 = 2
	Cost2
	// Cost3 = 3
	Cost3
	// Cost4 = 4
	Cost4
	// Cost5 = 5
	Cost5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

// String is a method to implement Stringer interface
func (l Cost) String() string {
	for s, v := range costStrings {
		if l == v {
			return s
		}
	}
	return "invalid"
}

// ParseCost func returns Cost values for its string representations
func ParseCost(s string) Cost {
	return costStrings[s]
}

// CostRange is a type for understanding ranges of Cost
type CostRange struct {
	From Cost
	To   Cost
}

// String is a method to implement Stringer interface
func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

// ParseCostRange is a function that takes string representation of
// cost range (e.g.: "$...$$" is a range 1-2) and creates CostRange object
// from it
func ParseCostRange(s string) (CostRange, error) {
	var r CostRange
	segs := strings.Split(s, "...")

	if len(segs) != 2 {
		return r, errors.New("invalid cost range")
	}

	r.From = ParseCost(segs[0])
	r.To = ParseCost(segs[1])

	return r, nil
}
