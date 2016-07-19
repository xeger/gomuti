package gomuti

import (
	"fmt"
	"reflect"
	"strings"
)

// Allow accepts an instance of Mock, or any struct contains a Mock, and returns
// an object that can be used to program an expected method call to the mock.
func Allow(double interface{}) *Allowed {
	m := find(reflect.ValueOf(double))
	return m.Allow()
}

// Â is an alias for Allow.
func Â(double interface{}) *Allowed {
	return Allow(double)
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
