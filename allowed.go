package gomuti

import "github.com/onsi/gomega"

// Allowed is a DSL object that lets you specify parameters, return values
// and other behaviors for a mocked call. For details on usage, see the
// documentation for Allow() and Â().
type Allowed struct {
	mock Mock
	last string
}

func (a *Allowed) params(params ...interface{}) []Matcher {
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

// Call allows the mock to receive a method call with matching parameters and
// return a specific set of values.
func (a *Allowed) Call(method string, params ...interface{}) *Allowed {
	calls := a.mock[method]
	calls = append(calls, allowed{})
	a.mock[method] = calls
	a.last = method
	return a
}

// With defines a Matcher for each method parameter; the method call
// is not matched unless all matchers are satisfied.
func (a *Allowed) With(params ...interface{}) *Allowed {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("mock: must use Call() before specifying With()")
	}
	call := &calls[len(calls)-1]
	if call.Params != nil {
		panic("mock: cannot specify With() twice")
	}
	call.Params = a.params(params...)
	return a
}

// Return specifies what the mock should return when a method call is matched.
// It must be called after Call/ToReceive.
func (a Allowed) Return(results ...interface{}) *Allowed {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("mock: must use Call() before specifying Return()")
	}
	call := &calls[len(calls)-1]
	if call.Results != nil {
		panic("mock: cannot specify Return() twice")
	}
	call.Results = results
	return &a
}

// Panic specifies that the mock should panic with the given reason when
// a method call is matched. It must be called after Call/ToReceive.
func (a Allowed) Panic(reason interface{}) *Allowed {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("mock: must use Call() before specifying Panic()")
	}
	call := &calls[len(calls)-1]
	if call.Panic != nil {
		panic("mock: cannot specify Panic() twice")
	}
	call.Panic = reason
	return &a
}

// ToReceive is an alias for Call()
func (a *Allowed) ToReceive(method string, params ...interface{}) *Allowed {
	return a.Call(method, params...)
}

// AndReturn is an alias for Return()
func (a *Allowed) AndReturn(results ...interface{}) *Allowed {
	return a.Return(results...)
}

// AndPanic is an alias for Panic()
func (a *Allowed) AndPanic(reason interface{}) *Allowed {
	return a.Panic(reason)
}
