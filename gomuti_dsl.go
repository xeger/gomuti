// Package gomuti provides a DSL that makes it easy to create test doubles
// (mocks, spies and stubs) for Golang interfaces. The DSL consists of three
// components:
//
// 1) Allow: a mocking method that records behavior for test doubles.
//
// 2) HaveCall: a spying method that verifies test doubles were called in an expected way.
//
// 3) Anything, AnythingOfType: parameter matchers used with mocking and spying methods.
// Gomega matchers can be used as Gomuti parameter matchers: BeNumerically, HaveOccurred, etc.
//
// All of these methods rely on the Mock and Spy types exported by package
// gomuti/types; test doubles are generally struct types that contain exported
// fields of type Mock and Spy. The DSL operates on pointers to these structs
// and uses reflection to access their fields.
//
// The DSL methods accept struct values as well as pointer-to-struct; the
// benefit of passing pointers is that the nested Mock or Spy will be allocated
// as needed with no intervention by the caller.
//
// Stubbing is provided by the mongoose package (https://github.com/xeger/mongoose),
// which also generates Gomuti-compatible mock code for any interface. Stubbed methods
// are called whenever no mock expectations match a method call; the return value(s)
// from a stubbed method call are always zero values. Stubbing must be enabled on
// a per-object basis by setting the Stub field to true.
package gomuti

import (
	"fmt"
	"reflect"

	"github.com/xeger/gomuti/types"
)

// Allow is a mocking method. It accepts a test double and returns a DSL-context
// object whose methods allow you to specify the test double's behavior when
// its methods are called.
//
// For information about parameter matching and return values, see types.Allowed.
func Allow(double interface{}) *types.Allowed {
	m := types.FindMock(reflect.ValueOf(double))
	return m.Allow()
}

// Â is a mocking method that is an alias for Allow. Use Shift+Option+M to
// type this symbol on Mac; Alt+0194 on Windows.
//
// As an additional shortcut, Â accepts the mocked method name and parameters
// as variadic parameters. The caller is still responsible for completing the
// behavior by calling Return or Panic() on the returned object.
//
// Examples:
//     Â(double).Call("Foo").With(1,1).Return(2)  // no shortcuts
//     Â(double, "Foo").With(2,1).Return(3)       // shortcut call
//     Â(double, "Foo",3,1).Return(4)             // shortcut params
func Â(double interface{}, methodAndParams ...interface{}) *types.Allowed {
	if len(methodAndParams) == 0 {
		return Allow(double)
	}

	m, ok := methodAndParams[0].(string)
	if !ok {
		panic(fmt.Sprintf("gomuti.Â: expected string as method name; got %T", methodAndParams[0]))
	}

	p := methodAndParams[1:]
	if len(p) > 0 {
		return Allow(double).Call(m).With(p)
	}
	return Allow(double).Call(m)
}
