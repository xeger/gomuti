package matchers

import (
	"reflect"

	"github.com/onsi/gomega/format"
)

// HaveTypeMatcher matches any value whose type matches the specified name.
// Named types must be prefixed with the name of the package in which they
// are defined (i.e. with the name that appears in the package statement of
// the source file where they are defined).
type HaveTypeMatcher struct {
	Expected string
}

// Match returns true if the type name of actual is similar to the expected
// name. The type name of actual is determined by calling reflect.TypeOf().
func (m *HaveTypeMatcher) Match(actual interface{}) (bool, error) {
	tn := reflect.TypeOf(actual).String()
	return (tn == m.Expected), nil
}

// FailureMessage returns a description of why the matcher did not match.
func (m *HaveTypeMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to have type", m.Expected)
}

// NegatedFailureMessage returns a description of why the matcher matched.
func (m *HaveTypeMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to have type", m.Expected)
}
