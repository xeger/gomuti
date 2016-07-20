package gomuti

import "fmt"

// SpyMatcher consults the spy of a test double in order to verify that
// method calls were received with specified parameters.
type SpyMatcher struct {
}

// HaveCall returns a SpyMatcher that verifies the named method was called.
// You can add additional constraints to the verification by calling With()
// or Returning().
func HaveCall(method string) *SpyMatcher {
	return &SpyMatcher{}
}

// HaveReceived is an alias for HaveCall().
func HaveReceived(method string) *SpyMatcher {
	return HaveCall(method)
}

// Match verifies that a method was called on a mock
// TODO actually implement this method
func (sm *SpyMatcher) Match(actual interface{}) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

// FailureMessage returns an explanation of the method call that was expected.
func (sm *SpyMatcher) FailureMessage(actual interface{}) (message string) {
	panic("not implemented")
}

// NegatedFailureMessage returns an explanation of the method call that was unexpected.
func (sm *SpyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic("not implemented")
}

// With adds an expectation about method parameters.
func (sm *SpyMatcher) With(params ...interface{}) *SpyMatcher {
	return sm
}

// Times adds an expectation about the number of times a method was called.
func (sm *SpyMatcher) Times(number int) *SpyMatcher {
	return sm
}

// Once is a shortcut for Times(1).
func (sm *SpyMatcher) Once() *SpyMatcher {
	return sm.Times(1)
}

// Twice is a shortcut for Times(2).
func (sm *SpyMatcher) Twice() *SpyMatcher {
	return sm.Times(2)
}
