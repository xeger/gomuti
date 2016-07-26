package gomuti

import (
	"fmt"
	"reflect"

	"github.com/xeger/gomuti/types"
)

// Allow accepts an instance of Mock, or any struct contains a Mock, and returns
// an object that can be used to program an expected method call to the mock.
func Allow(double interface{}) *types.Allowed {
	m := types.FindMock(reflect.ValueOf(double))
	return m.Allow()
}

// Â is an alias for Allow. Use Shift+Option+M to type this symbol on Mac; Alt+0194 on Windows.
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
