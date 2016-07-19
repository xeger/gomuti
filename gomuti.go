package gomuti

import (
	"fmt"
	"reflect"
	"strings"
)

// Allow returns an object that can be used to program an expected method call to
// the mock.
func Allow(double interface{}) *Allowed {
	m := find(reflect.ValueOf(double))
	return &Allowed{mock: m}
}

// Â is an alias for Allow.
func Â(double interface{}) *Allowed {
	return Allow(double)
}

// Ø is used to delegate behavior to instances of Mock. Not meant to be called
// directly. If it returns non-nil, then the method call was matched; methods
// that return nothing still return an empty slice.
//
// In contrast, if this method returns nil then the method call was NOT
// matched and the caller should behave accordingly, i.e. panic unless some
// stubbed default behavior is appropriate.
func Ø(mock Mock, method string, params ...interface{}) []interface{} {
	if mock == nil {
		return nil
	}
	calls := mock[method]

	for _, c := range calls {
		if len(c.Params) == len(params) {
			matched := true
			for i, p := range params {
				success, err := c.Params[i].Match(p)
				if err != nil {
					panic(err.Error())
				}
				matched = matched && success
			}
			if matched {
				if c.Panic != nil {
					panic(c.Panic)
				}
				return c.Results
			}
		} else if c.Params == nil {
			return c.Results
		}
	}
	return nil
}

func isMock(t reflect.Type) bool {
	return t.String() == "gomuti.Mock" && strings.Index(t.PkgPath(), "gomuti") > 0
}

// Find the Mock associated with an arbitrary value and initialize it if
// necessary; panic if no Mock is found or a new Mock cannot be initialized.
func find(v reflect.Value) Mock {
	t := v.Type()
	ptr := (t.Kind() == reflect.Ptr)
	if ptr {
		t = t.Elem()
	}

	if isMock(t) {
		// The real McCoy! (Or a pointer to it.)
		if ptr {
			if v.IsNil() {
				panic(fmt.Sprintf("mock.Allow: must initialize %s before calling", v.Type().String()))
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
					v = reflect.Indirect(v)
					f := reflect.Indirect(v).Field(i)
					if !f.CanInterface() {
						panic(fmt.Sprintf("mock.Allow: cannot work with unexported field %s of %s; change it to %s", sf.Name, t.String(), strings.Title(sf.Name)))
					}
				}
				mock = v.Field(i).Interface().(Mock)
				if mock == nil {
					if ptr {
						mock = Mock{}
						reflect.Indirect(v).Field(i).Set(reflect.ValueOf(mock))
					} else {
						panic(fmt.Sprintf("mock.Allow: must pass a pointer to %s or initialize its .Mock before calling", t.String()))
					}
				}
				return mock
			}
		}
	}
	panic(fmt.Sprintf("mock: don't know how to program behaviors for %s", t.String()))
}
