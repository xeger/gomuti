package types

import (
	"fmt"
	"reflect"
)

// Allowed is a DSL object that lets you specify parameters, return values
// and other behaviors for a mocked call. Calls to gomuti.Allow() return
// this type; each method that you call, refines the mock behavior that you
// are defining. The behavior is "completed" when you specify an outcome
// for the mock behavior by calling Return, Panic or Do.
//
// When your mock receives a method call, it is compared to each call that
// you have allowed for that method name. Gomuti matches the actual parameters
// against the allowed call's parameter matchers and computes a score for that
// allowed call, then picks the allowed call with the highest score. If two or
// more calls have an identical score, the most recently-allowed call wins.
//
// Gomuti uses the following rules to score method calls:
//
//
// 1) If the number of actual parameters varies from the allowed, score = 0
//    (disqualify the allowed call as a match).
//
// 2) If the allowed call did not specify any parameters, score = 1
//    (the allowed call matches any number of actual parameters, but just barely).
//
// 3) For each parameter that matches a gomega.BeEqual matcher, score += 4.
//
// 4) For each parameter that matches a gomega.BeEquivalentTo matcher, score += 3
//
// 5) For each parameter that matches another matcher, score += 2
//
// The scoring algorithm sounds complicated, but it results in a very natural-
// feeling matching behavior. Imagine that we are mocking a method
// Add(a, b, c interface{}) which adds its parameters:
//
//     Â(double, "Add").Panic("so confused")
//     Â(double, "Add").With(1,2,3).Return(6)
//     Â(double, "Add").With(1,2,AnythingOfType("int")).Panic("so number. very count. wow.")
//     Â(double, "Add").With(BeNumerically(">",0),EquivalentTo(2.0),Anything()).Return(0)
//
//     double.Add(7,2,"hi") // returns 0
//     double.Add(1,2,99)   // panics with a doge message
//     double.Add(1,2,3)    // returns 6
//     double.Add(8,8,8)    // panics with a confused message
type Allowed struct {
	mock Mock
	last string
}

// Call allows the mock to receive a method call with matching parameters and
// return a specific set of values.
//
// There is typically no need to call this method directly since calls to
// gomuti.Allow() have already set the method name on the Allowed that they
// return to you.
//
// If you call this method twice on the same Allowed, gomuti panics.
func (a *Allowed) Call(method string, params ...interface{}) *Allowed {
	if a.last != "" {
		panic("gomuti: cannot use Call() twice on the same Allowed")
	}

	calls := a.mock[method]
	calls = append(calls, Call{})
	a.mock[method] = calls
	a.last = method
	if len(params) > 0 {
		a.With(params...)
	}
	return a
}

// With allows you to match method parameters of a mock call. You provide
// a Matcher or a literal value for each position; when someone calls the mock,
// the call is considered to be matched to this behavior if all parameters
// match the values you have provided.
//
// If you call this method twice on the same Allowed, gomuti panics.
//
// MATCHING LITERAL PARAMETER VALUES
//
// Literal values of basic type (int, bool, string, etc) are converted to
// equality matchers; literal pointers and interface types (slice, map, etc)
// are converted to equivalency matchers. Therefore, calling this:
//
//     Â(double, "Foo").With("hello", 42, time.Now())
//
// is the same as calling this:
//
//     Â(double, "Foo").With(Equals("hello"), Equals(42), BeEquivalentTo(time.Now()))
//
// MATCHING RANGES OF VALUES
//
// Any Gomega or Gomuti matcher can be used to match a parameter, enabling
// very sophisticated behavior; for instance:
//
//     Â(double, "Foo").With(
//       MatchRegexp("hello|goodbye"),
//       BeNumerically("<=", 42), BeTrue()
//     )
//
// MATCHING VARIADIC PARAMETERS
//
// Variadic parameters are converted to a slice of values and matched against
// the final matcher provided to With. For instance:
//
//   Â(double, "Foo").With(true, []int{1,2,3})
//   double.Foo(true, 1, 2, 3)
//
// If you call this method twice on the same Allowed, gomuti panics.
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

// Do allows you to provide a function that the mock will call in order to
// determine the correct behavior when a call is matched. You can use it to
// cause side effects or record parameters of your mock calls.
//
// You can pass any function signature to Do; however, if the signature
// does not match the signature of the method being mocked, gomuti will cause
// a panic when you call the mock method. Use this method with care!
//
// If you call this method twice on the same Allowed, gomuti panics.
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
