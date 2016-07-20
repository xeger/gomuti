package gomuti

import "fmt"

// Mock is a state container for mocked behavior. Rather than instantiating it
// directly, you should include a field of this type in your mock structs and
// define methods that delegate their behavior to the mock's Delegate() method.
//
// If all of this sounds like too much work, then you should really check out
// https://github.com/xeger/mongoose to let the computer generate your mocks
// for you!
type Mock map[string][]call

var defaultReturn = make([]interface{}, 0)

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
		} else if c.Results != nil {
			return c.Results
		}
		// Lazy user didn't tell us to do, panic or return; assume he meant to
		// return nothing
		return defaultReturn
	}
	return nil
}

// Finds the closest matching call for the specified method, or nil if no
// calls match. Panics if two or more calls are an equally good match.
func (m Mock) bestMatch(method string, params ...interface{}) *call {
	calls := m[method]

	matches := make([]*call, 0, 3)
	bestScore := 0

	for _, c := range calls {
		score := c.score(params)
		if score > 0 && score >= bestScore {
			matches = append(matches, &c)
			bestScore = score
		}
	}

	switch len(matches) {
	case 0:
		return nil
	case 1:
		return matches[0]
	default:
		panic(fmt.Sprintf("gomuti: matched %d mocked calls to %s; don't know which to behave like", len(matches), method))
	}
}
