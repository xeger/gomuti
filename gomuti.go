package gomuti

import "reflect"

// Allow accepts an instance of Mock, or any struct contains a Mock, and returns
// an object that can be used to program an expected method call to the mock.
func Allow(double interface{}) *Allowed {
	m := findMock(reflect.ValueOf(double))
	return m.Allow()
}

// Â is an alias for Allow.
func Â(double interface{}) *Allowed {
	return Allow(double)
}
