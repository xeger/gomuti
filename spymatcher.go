package gomuti

import (
	"bytes"
	"fmt"
	"reflect"
)

// SpyMatcher consults the spy of a test double in order to verify that
// method calls were received with specified parameters.
type SpyMatcher struct {
	Method string
	Params []Matcher
	Count  int
}

// HaveCall returns a SpyMatcher that verifies the named method was called.
// You can add additional constraints to the verification by calling With()
// or Returning().
func HaveCall(method string) *SpyMatcher {
	return &SpyMatcher{Method: method, Count: 1}
}

// HaveReceived is an alias for HaveCall().
func HaveReceived(method string) *SpyMatcher {
	return HaveCall(method)
}

// Match verifies that a method was called on a mock
func (sm *SpyMatcher) Match(actual interface{}) (bool, error) {
	spy := findSpy(reflect.ValueOf(actual))
	if spy == nil {
		return false, fmt.Errorf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)
	return matched >= sm.Count, nil
}

// FailureMessage returns an explanation of the method call that was expected.
func (sm *SpyMatcher) FailureMessage(actual interface{}) (message string) {
	spy := findSpy(reflect.ValueOf(actual))
	if spy == nil {
		return fmt.Sprintf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)
	return sm.describe("Expected", matched)
}

// NegatedFailureMessage returns an explanation of the method call that was unexpected.
func (sm *SpyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	spy := findSpy(reflect.ValueOf(actual))
	if spy == nil {
		return fmt.Sprintf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)
	return sm.describe("Did not expect", matched)
}

// With adds an expectation about method parameters.
func (sm *SpyMatcher) With(params ...interface{}) *SpyMatcher {
	sm.Params = paramsToMatchers(params)
	return sm
}

// Times adds an expectation about the number of times a method was called.
func (sm *SpyMatcher) Times(number int) *SpyMatcher {
	sm.Count = number
	return sm
}

// Never is a shortcut for Times(0)
func (sm *SpyMatcher) Never() *SpyMatcher {
	return sm.Times(0)
}

// Once is a shortcut for Times(1).
func (sm *SpyMatcher) Once() *SpyMatcher {
	return sm.Times(1)
}

// Twice is a shortcut for Times(2).
func (sm *SpyMatcher) Twice() *SpyMatcher {
	return sm.Times(2)
}

func (sm *SpyMatcher) describe(lede string, got int) string {
	var ecalls, gcalls string
	if sm.Count == 1 {
		ecalls = "call"
	} else {
		ecalls = "calls"
	}

	if got == 1 {
		gcalls = "call"
	} else {
		gcalls = "calls"
	}

	b := bytes.NewBufferString(fmt.Sprintf("%s %d %s to %s", lede, sm.Count, ecalls, sm.Method))
	if len(sm.Params) > 0 {
		b.WriteString(" with:\n")
		formatMatcherInfo(b, 2, sm.Params)
	} else {
		b.WriteString(" ")
	}
	b.WriteString(fmt.Sprintf("but got %d %s", got, gcalls))
	return b.String()
}
