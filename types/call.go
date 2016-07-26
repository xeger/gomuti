package types

import "github.com/onsi/gomega/matchers"

// Internal representation of a call that has been programmed on a Mock via
// calls to an Allowed.
type call struct {
	Params  []Matcher
	Do      CallFunc
	Panic   interface{}
	Results []interface{}
}

// Determine whether this call's Params match the given params, and if so,
// how well. Returns 0 without match, or a positive integer with match.
func (c *call) score(params []interface{}) int {
	if c.Params == nil {
		// a call with no params matches any params, but just barely...
		return 1
	} else if len(c.Params) == len(params) {
		// compute the score by considering all matchy matchers
		score := 0
		for i, p := range params {
			success, err := c.Params[i].Match(widen(p))
			if err != nil {
				panic(err.Error())
			} else if !success {
				return 0
			}
			switch c.Params[i].(type) {
			case *matchers.EqualMatcher:
				score += 4
			case *matchers.BeEquivalentToMatcher:
				score += 3
			default:
				score += 2
			}
		}
		return score
	}
	return 0
}
