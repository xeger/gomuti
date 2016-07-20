package gomuti

import (
	"fmt"
	"reflect"
	"strings"
)

type called struct {
	Params []interface{}
}

// Spy is a state container for recording information about calls made to a
// test double.
type Spy map[string][]called

// Observe records a method call.
func (s Spy) Observe(method string, params ...interface{}) {
	events := s[method]
	events = append(events, called{Params: params})
	s[method] = events
}

// Count returns the number of times a method was called that matched the given
// criteria.
func (s Spy) Count(method string, criteria ...Matcher) int {
	events := s[method]
	res := 0

	for _, ev := range events {
		if len(ev.Params) < len(criteria) {
			continue
		}
		failed := false
		for i, c := range criteria {
			succ, err := c.Match(ev.Params[i])
			if err != nil {
				panic(err)
			}
			if !succ {
				failed = true
				break
			}
		}
		if !failed {
			res++
		}
	}

	return res
}

func isSpy(t reflect.Type) bool {
	return t.String() == "gomuti.Spy" && strings.Index(t.PkgPath(), "gomuti") > 0
}

// Find the Spy associated with an arbitrary value and initialize it if
// necessary; panic if no Mock is found or a new Mock cannot be initialized.
func findSpy(v reflect.Value) Spy {
	t := v.Type()
	ptr := (t.Kind() == reflect.Ptr)
	if ptr {
		t = t.Elem()
	}

	if isSpy(t) {
		// The real McCoy! (Or a pointer to it.)
		if ptr {
			if v.IsNil() {
				panic(fmt.Sprintf("gomuti.Allow: must initialize %s before calling", v.Type().String()))
			}
			return reflect.Indirect(v).Interface().(Spy)
		}
		return v.Interface().(Spy)
	} else if t.Kind() == reflect.Struct {
		// A struct type (or pointer-to-struct); search its fields for a Mock.
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if isSpy(sf.Type) {
				// Found a field. Initialize if necessary (and possible) and return
				// the Mock interface value of the field.
				var spy Spy
				if ptr {
					if v.IsNil() {
						panic(fmt.Sprintf("gomuti: must initialize *%s before calling", t.Name()))
					}
					v = reflect.Indirect(v)
					if !v.IsValid() {
						panic(fmt.Sprintf("gomuti: must initialize %s.%s before calling", t.Name(), sf.Name))
					}
					f := v.Field(i)
					if !f.CanInterface() {
						panic(fmt.Sprintf("gomuti: cannot work with unexported field %s of %s; change it to %s", sf.Name, t.String(), strings.Title(sf.Name)))
					}
				}
				spy = v.Field(i).Interface().(Spy)
				if spy == nil {
					if ptr {
						spy = Spy{}
						reflect.Indirect(v).Field(i).Set(reflect.ValueOf(spy))
					} else {
						panic(fmt.Sprintf("gomuti: must pass a pointer to %s or initialize its .Spy before calling", t.String()))
					}
				}
				return spy
			}
		}
	}
	panic(fmt.Sprintf("gomuti: don't know how to spy on %s", t.String()))
}
