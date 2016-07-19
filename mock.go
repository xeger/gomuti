package gomuti

// Internal representation of a call that has been programmed via Allowed.
type allowed struct {
	Params  []Matcher
	Do      DoFunc
	Panic   interface{}
	Results []interface{}
}

// Mock is a state container for mocked behavior. Rather than instantiating it
// directly, you should include a field of this type in your mock structs and
// define methods that delegate their behavior to the mock's Delegate() method.
//
// If all of this sounds like too much work, then you should really check out
// https://github.com/xeger/mongoose
type Mock map[string][]allowed

// Allow returns an object that can be used to program an expected method call.
// Rather than calling this directly, you probably want to call gomuti.Allow()
// on some struct that contains a Mock.
func (m Mock) Allow() *Allowed {
	return &Allowed{mock: m}
}

// Delegate informs the mock that a call has been made; if the call matches
// a call that was programmed with Allow(), it returns non-nil. Methods
// that return nothing, still return an empty slice if the call was matched.
//
// In contrast, if this method returns nil then the method call was NOT
// matched and the caller should behave accordingly, i.e. panic unless some
// stubbed default behavior is appropriate.
func (m Mock) Delegate(method string, params ...interface{}) []interface{} {
	if m == nil {
		return nil
	}

	c := m.bestMatch(method, params...)
	if c != nil {
		if c.Do != nil {
			return c.Do(params...)
		} else if c.Panic != nil {
			panic(c.Panic)
		}
		return c.Results
	}
	return nil
}

// Finds the closest matching call
// TODO: actual find closest match, rather than finding exact match
func (m Mock) bestMatch(method string, params ...interface{}) *allowed {
	calls := m[method]

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
				return &c
			}
		} else if c.Params == nil {
			return &c
		}
	}
	return nil
}
