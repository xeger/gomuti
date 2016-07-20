package gomuti

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/onsi/gomega"
)

// Matcher is a method parameter matcher for an Allowed call. It is a subset
// of Gomega's Matcher interface, but does not contain any error-reporting
// functions since a failure to match method parameters is not by itself an
// error.
type Matcher interface {
	Match(actual interface{}) (success bool, err error)
}

// Returns a multi-line string describing the position and nature of
// each matcher in a list. Indents each line the specified number of spaces.
func formatMatcherInfo(b *bytes.Buffer, indent int, params []Matcher) string {
	spacer := strings.Repeat(" ", indent)
	for i, p := range params {
		b.WriteString(fmt.Sprintf("%s%2d: %s\n", spacer, i, matcherString(p)))
	}
	return b.String()
}

// Noise words that tend to appear in matcher typenames.
var noise = regexp.MustCompile("^[a-z*]+[.]|Matcher")

// Returns a human-readable description of a matcher and its expected value.
// Uses reflection to grab expected values from any matcher, and removes noise
// words and package prefixes from the matcher type name.
func matcherString(m Matcher) string {
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

// Converts all non-matcher parameters to an equivalency matcher.
func paramsToMatchers(params []interface{}) []Matcher {
	matchers := make([]Matcher, len(params))
	for i, p := range params {
		m, ok := p.(Matcher)
		if ok {
			matchers[i] = m
		} else {
			matchers[i] = gomega.BeEquivalentTo(p)
		}
	}
	return matchers
}
