package matchers

import "github.com/onsi/gomega/format"

// BeAnythingMatcher matches _any_ value including nil and zero values.
type BeAnythingMatcher struct{}

// Match always returns true.
func (m *BeAnythingMatcher) Match(actual interface{}) (bool, error) {
	return true, nil
}

// FailureMessage returns a description of why the matcher did not match.
func (m *BeAnythingMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to be any value")
}

// NegatedFailureMessage returns a description of why the matcher matched.
func (m *BeAnythingMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to be a value")
}
