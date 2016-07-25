package types

import (
	"reflect"

	"github.com/onsi/gomega"
)

// MatchParams returns a sequence of gomuti Matchers that will match the specified
// method-parameter sequence. For parameters that are not already a Matcher,
// it uses a heuristic to create an equality, equivalency or be-nil matcher. For
// parameters that are already a matcher, it returns the matcher verbatim.
func MatchParams(params []interface{}) []Matcher {
	matchers := make([]Matcher, len(params))
	for i, p := range params {
		m, ok := p.(Matcher)
		if ok {
			matchers[i] = m
		} else if p == nil {
			matchers[i] = gomega.BeNil()
		} else {
			switch reflect.TypeOf(p).Kind() {
			case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
				matchers[i] = gomega.BeEquivalentTo(p)
			default:
				matchers[i] = gomega.Equal(p)
			}
		}
	}
	return matchers
}
