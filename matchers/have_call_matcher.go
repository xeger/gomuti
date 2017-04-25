package matchers

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/xeger/gomuti/types"
)

// Noise words that tend to appear in matcher typenames:
//   - initial package name i.e. "matchers.Xyz"
//   - the word "Matcher" as a suffix i.e. "BeAnythingMatcher"
//
// By stripping noise words from a type name, we generally recover the name of
// the DSL function that is used to instantiate a matcher of that type. This
// is due to the naming convention that gomega/matchers and gomuti/matchers
// rigidly follow.
var noise = regexp.MustCompile("^[a-z*]+[.]|Matcher")

// Returns a human-readable description of a matcher and its expected value.
// Uses reflection to grab expected values from any matcher, and removes noise
// words and package prefixes from the matcher type name.
func matcherString(m types.Matcher) string {
	v := reflect.ValueOf(m)

	if v.Kind() == reflect.Ptr {
		// Common case: pointer to matcher struct. Describe it like a method call
		// (which it probably initially was).
		t := v.Elem().Type()
		_, ok := t.FieldByName("Expected")
		nam := noise.ReplaceAllString(t.Name(), "")
		if ok {
			exp := v.Elem().FieldByName("Expected")
			return fmt.Sprintf("%s(%#v)", nam, exp)
		}
		return nam
	}

	// Oddball case: a matcher that is a wrapped interface type. Describe it
	// as best we can...
	return fmt.Sprintf("%#v", m)
}

// Returns a multi-line string describing the position and nature of
// each matcher in a list. Indents each line the specified number of spaces.
func formatMatcherInfo(b *bytes.Buffer, indent int, params []types.Matcher) string {
	spacer := strings.Repeat(" ", indent)
	for i, p := range params {
		b.WriteString(fmt.Sprintf("%s%2d: %s\n", spacer, i, matcherString(p)))
	}
	return b.String()
}

// Returns a multi-line string describing the position and nature of
// each recorded actual parameter in a list. Indents each line the
// specified number of spaces.
func formatParamInfo(b *bytes.Buffer, indent int, params []interface{}) string {
	spacer := strings.Repeat(" ", indent)
	for i, p := range params {
		v := fmt.Sprintf("%#v", p)
		if len(v) > 60 {
			v = fmt.Sprintf("%59.59s…", v)
		}
		b.WriteString(fmt.Sprintf("%s%2d: %s\n", spacer, i, v))
	}
	return b.String()
}

// HaveCallMatcher consults the spy of a test double in order to verify that
// method calls were received with specified parameters.
type HaveCallMatcher struct {
	Method string
	Params []types.Matcher
	Count  int
}

// Match verifies that a method was called on a mock
func (sm *HaveCallMatcher) Match(actual interface{}) (bool, error) {
	spy := types.FindSpy(reflect.ValueOf(actual))
	if spy == nil {
		return false, fmt.Errorf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)
	return matched >= sm.Count, nil
}

// FailureMessage returns an explanation of the method call that was expected.
func (sm *HaveCallMatcher) FailureMessage(actual interface{}) (message string) {
	spy := types.FindSpy(reflect.ValueOf(actual))
	if spy == nil {
		return fmt.Sprintf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)

	var closest []interface{}
	if spy.Count(sm.Method) > 0 {
		closest = spy.ClosestMatch(sm.Method, sm.Params...)
	}

	return sm.describe("Expected", matched, closest)
}

// NegatedFailureMessage returns an explanation of the method call that was unexpected.
func (sm *HaveCallMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	spy := types.FindSpy(reflect.ValueOf(actual))
	if spy == nil {
		return fmt.Sprintf("Cannot spy on %T", actual)
	}
	matched := spy.Count(sm.Method, sm.Params...)
	return sm.describe("Did not expect", matched, nil)
}

// With adds an expectation about method parameters.
func (sm *HaveCallMatcher) With(params ...interface{}) *HaveCallMatcher {
	sm.Params = types.MatchParams(params)
	return sm
}

// Times adds an expectation about the number of times a method was called.
func (sm *HaveCallMatcher) Times(number int) *HaveCallMatcher {
	sm.Count = number
	return sm
}

// Never is a shortcut for Times(0)
func (sm *HaveCallMatcher) Never() *HaveCallMatcher {
	return sm.Times(0)
}

// Once is a shortcut for Times(1).
func (sm *HaveCallMatcher) Once() *HaveCallMatcher {
	return sm.Times(1)
}

// Twice is a shortcut for Times(2).
func (sm *HaveCallMatcher) Twice() *HaveCallMatcher {
	return sm.Times(2)
}

func (sm *HaveCallMatcher) describe(lede string, got int, closest []interface{}) string {
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

	if got == 0 && closest != nil {
		b.WriteString("but no call matched exactly. Closest match:\n")
		formatParamInfo(b, 2, closest)
	} else {
		b.WriteString(fmt.Sprintf("but observed %d %s", got, gcalls))
	}
	return b.String()
}
