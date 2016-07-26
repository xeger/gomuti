package types

import (
	"fmt"
	"reflect"
)

// Allowed is a DSL object that lets you specify parameters, return values
// and other behaviors for a mocked call. For details on usage, see the
// documentation for Allow() and Â().
type Allowed struct {
	mock Mock
	last string
}

// Call allows the mock to receive a method call with matching parameters and
// return a specific set of values.
func (a *Allowed) Call(method string, params ...interface{}) *Allowed {
	calls := a.mock[method]
	calls = append(calls, call{})
	a.mock[method] = calls
	a.last = method
	return a
}

// With defines a Matcher for each method parameter; the method call
// is not matched unless all matchers are satisfied.
func (a *Allowed) With(params ...interface{}) *Allowed {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("gomuti: must use Call() before specifying With()")
	}
	call := &calls[len(calls)-1]
	if call.Params != nil {
		panic("gomuti: cannot specify With() twice")
	}
	call.Params = MatchParams(params)
	return a
}

func do(fn interface{}) CallFunc {
	cf, ok := fn.(CallFunc)
	if ok {
		return cf
	}

	fp, ok := fn.(func(...interface{}) []interface{})
	if ok {
		return CallFunc(fp)
	}

	v := reflect.ValueOf(fn)
	switch v.Kind() {
	case reflect.Func:
		return func(params ...interface{}) []interface{} {
			in := make([]reflect.Value, 0, len(params))
			for _, p := range params {
				in = append(in, reflect.ValueOf(p))
			}
			out := v.Call(in)
			result := make([]interface{}, 0, len(out))
			for _, o := range out {
				if o.CanInterface() {
					result = append(result, o.Interface())
				} else {
					panic(fmt.Sprintf("gomuti: CallFunc adapter cannot handle %s (CanInterface==false)", o.String()))
				}
			}
			return result
		}
	default:
		panic(fmt.Sprintf("gomuti: cannot convert %T into CallFunc", fn))
	}
}

// Do provides a function that the mock will call in order to determine the
// correct behavior when a call is matched.
func (a Allowed) Do(doer interface{}) {
	df := do(doer)
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("gomuti: must use Call() before specifying Do()")
	}
	call := &calls[len(calls)-1]
	a.behave(df, call.Panic, call.Results)
	if call.Do != nil {
		panic("gomuti: cannot specify Do() twice")
	}
	call.Do = df
}

// Return specifies what the mock should return when a method call is matched.
// It must be called after Call/ToReceive.
func (a Allowed) Return(results ...interface{}) {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("gomuti: must use Call() before specifying Return()")
	}
	call := &calls[len(calls)-1]
	a.behave(call.Do, call.Panic, results)
	if call.Results != nil {
		panic("gomuti: cannot specify Return() twice")
	}
	call.Results = results
	return
}

// Panic specifies that the mock should panic with the given reason when
// a method call is matched. It must be called after Call/ToReceive.
func (a Allowed) Panic(reason interface{}) {
	calls := a.mock[a.last]
	if calls == nil || len(calls) < 1 {
		panic("mock: must use Call() before specifying Panic()")
	}
	call := &calls[len(calls)-1]
	a.behave(call.Do, reason, call.Results)
	if call.Panic != nil {
		panic("gomuti: cannot specify Panic() twice")
	}
	call.Panic = reason
	return
}

// ToReceive is an alias for Call()
func (a *Allowed) ToReceive(method string, params ...interface{}) *Allowed {
	return a.Call(method, params...)
}

// AndReturn is an alias for Return()
func (a *Allowed) AndReturn(results ...interface{}) {
	a.Return(results...)
}

// AndPanic is an alias for Panic()
func (a *Allowed) AndPanic(reason interface{}) {
	a.Panic(reason)
}

// Ensure that the user only specifies ONE behavior: Do, Panic or Return.
func (a Allowed) behave(d CallFunc, p interface{}, r []interface{}) {
	if d != nil && p != nil {
		panic("gomuti: cannot simultaneously Do() and Panic(); choose one")
	} else if d != nil && r != nil {
		panic("gomuti: cannot simultaneously Do() and Return(); choose one")
	} else if p != nil && r != nil {
		panic("gomuti: cannot simultaneously Panic() and Return(); choose one")
	}
}
