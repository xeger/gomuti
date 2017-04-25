package types

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
	if s == nil {
		panic("gomuti: must initialize Spy before calling Observe")
	}

	events := s[method]
	events = append(events, called{Params: params})
	s[method] = events
}

// Count returns the number of times a method was called that matched the given
// criteria.
func (s Spy) Count(method string, criteria ...Matcher) int {
	if s == nil {
		panic("gomuti: must initialize Spy before calling Count")
	}

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

// ClosestMatch returns the parameters of the recorded call that most closely
// matches the given criteria, or nil if the method was never called at all.
func (s Spy) ClosestMatch(method string, criteria ...Matcher) []interface{} {
	if s == nil {
		panic("gomuti: must initialize Spy before calling ClosestMatch")
	}

	var best []interface{}
	var bestCount int

	for _, call := range s[method] {
		count := 0
		for i, crit := range criteria {
			if len(call.Params) > i {
				if ok, _ := crit.Match(call.Params[i]); ok {
					count++
				}
			}
		}
		if count > bestCount {
			best = call.Params
			bestCount = count
		}
	}

	return best
}

func isSpy(t reflect.Type) bool {
	return t.String() == "types.Spy" && strings.Index(t.PkgPath(), "gomuti") > 0
}

// FindSpy uses reflection to find the spy-controller associated with a given
// value. Its behavior varies depending on the type of the value:
//
// 1) Instance of Spy: return the value itself
// 2) Pointer to Spy: return the pointed-to value
// 3) Struct that contains a Spy field:
//      3a) if the field is nil, panic (user must initialize the field)
//      3b) return the field's value
// 4) Pointer to struct that contains a Spy field:
//      4a) if the field is nil, initialize it to an empty Spy
//      4b) return the field's value
// 6) Anything else: panic (don't know how to spy on ...)
func FindSpy(v reflect.Value) Spy {
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
	panic(fmt.Sprintf("gomuti: don't know how to spy on %s", v.Type().Name()))
}
