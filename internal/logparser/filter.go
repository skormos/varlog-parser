package logparser

import "strings"

type (
	// Filterer interface defines a contract to determine if a provided string passes a specific criteria.
	Filterer interface {
		Filter(input string) bool
	}

	// FiltererFn is a convenience function wrapper to implement Filterer.
	FiltererFn func(input string) bool
)

// Filter delegates to the function method receiver, passing the input parameter.
func (fn FiltererFn) Filter(input string) bool {
	return fn(input)
}

// FilterNone returns a Filterer instance that ignores the input and always returns true. Useful for when you want to
// Filter nothing and accept all lines.
func FilterNone() Filterer {
	return FiltererFn(func(_ string) bool {
		return true
	})
}

// FilterOnSubstring calls strings.Contains with the provided substr parameter to check if a line should be accepted.
func FilterOnSubstring(substr string) Filterer {
	return FiltererFn(func(input string) bool {
		return strings.Contains(input, substr)
	})
}
