package gomuti

import (
	"fmt"
	"reflect"
	"strings"
)

// Mock is a state container for mocked behavior. Rather than instantiating it
// directly, you should include a field of this type in your mock structs and
// define methods that delegate their behavior to the mock's Delegate() method.
//
// If all of this sounds like too much work, then you should really check out
// https://github.com/xeger/mongoose to let the computer generate your mocks
// for you!
type Mock map[string][]call

var defaultReturn = make([]interface{}, 0)

// Allow returns an object that can be used to program an expected method call.
// Rather than calling this directly, you probably want to call gomuti.Allow()
// on some struct that contains a Mock.
func (m Mock) Allow() *Allowed {
	return &Allowed{mock: m}
}

// Delegate informs the mock that a call has been made; if the call matches
// a call that was programmed with Allow(), it returns non-nil. Methods
// that return nothing, still return an empty slice if the call was matched.
//
// In contrast, if this method returns nil then the method call was NOT
// matched and the caller should behave accordingly, i.e. panic unless some
// stubbed default behavior is appropriate.
func (m Mock) Delegate(method string, params ...interface{}) []interface{} {
	if m == nil {
		return nil
	}

	c := m.bestMatch(method, params...)
	if c != nil {
		if c.Do != nil {
			return c.Do(params...)
		} else if c.Panic != nil {
			panic(c.Panic)
		} else if c.Results != nil {
			return c.Results
		}
		// Lazy user didn't tell us to do, panic or return; assume he meant to
		// return nothing
		return defaultReturn
	}
	return nil
}

// Finds the closest matching call for the specified method, or nil if no
// calls match. Panics if two or more calls are an equally good match.
func (m Mock) bestMatch(method string, params ...interface{}) *call {
	calls := m[method]

	matches := make([]call, 0, 3)
	bestScore := 0

	for _, c := range calls {
		score := c.score(params)
		if score > 0 && score >= bestScore {
			matches = append(matches, c)
			bestScore = score
		}
	}

	switch len(matches) {
	case 0:
		return nil
	case 1:
		return &matches[0]
	default:
		panic(fmt.Sprintf("gomuti: matched %d mocked calls to %s; don't know which to behave like", len(matches), method))
	}
}

func isMock(t reflect.Type) bool {
	return t.String() == "gomuti.Mock" && strings.Index(t.PkgPath(), "gomuti") > 0
}

// Find the Mock associated with an arbitrary value and initialize it if
// necessary; panic if no Mock is found or a new Mock cannot be initialized.
func findMock(v reflect.Value) Mock {
	t := v.Type()
	ptr := (t.Kind() == reflect.Ptr)
	if ptr {
		t = t.Elem()
	}

	if isMock(t) {
		// The real McCoy! (Or a pointer to it.)
		if ptr {
			if v.IsNil() {
				panic(fmt.Sprintf("gomuti.Allow: must initialize %s before calling", v.Type().String()))
			}
			return reflect.Indirect(v).Interface().(Mock)
		}
		return v.Interface().(Mock)
	} else if t.Kind() == reflect.Struct {
		// A struct type (or pointer-to-struct); search its fields for a Mock.
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if isMock(sf.Type) {
				// Found a field. Initialize if necessary (and possible) and return
				// the Mock interface value of the field.
				var mock Mock
				if ptr {
					if v.IsNil() {
						panic(fmt.Sprintf("gomuti.Allow: must initialize *%s before calling", t.Name()))
					}
					v = reflect.Indirect(v)
					if !v.IsValid() {
						panic(fmt.Sprintf("gomuti.Allow: must initialize %s.%s before calling", t.Name(), sf.Name))
					}
					f := v.Field(i)
					if !f.CanInterface() {
						panic(fmt.Sprintf("gomuti.Allow: cannot work with unexported field %s of %s; change it to %s", sf.Name, t.String(), strings.Title(sf.Name)))
					}
				}
				mock = v.Field(i).Interface().(Mock)
				if mock == nil {
					if ptr {
						mock = Mock{}
						reflect.Indirect(v).Field(i).Set(reflect.ValueOf(mock))
					} else {
						panic(fmt.Sprintf("gomuti.Allow: must pass a pointer to %s or initialize its .Mock before calling", t.String()))
					}
				}
				return mock
			}
		}
	}
	panic(fmt.Sprintf("gomuti: don't know how to program behaviors for %s", t.String()))
}
