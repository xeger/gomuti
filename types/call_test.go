package types

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("call", func() {
	It("widens numeric values before scoring", func() {
		wide := []interface{}{uint64(0), int64(1), float64(2.0)}
		narrow := []interface{}{uint8(0), int8(1), float32(2.0)}
		c := call{Params: MatchParams(wide)}
		Expect(c.score(narrow)).NotTo(BeZero())
	})

	It("scores 1 given no params", func() {
		params := []interface{}{}
		params2 := []interface{}{1, 2, 3, 4, 5}
		c := call{}
		Expect(c.score(params)).To(Equal(1))
		Expect(c.score(params2)).To(Equal(1))
	})

	It("scores equality higher than equivalence", func() {
		c := call{
			Params: []Matcher{Equal(42), Equal(true)},
		}
		c2 := call{
			Params: []Matcher{BeEquivalentTo(42.0), BeEquivalentTo(true)},
		}
		high := c.score([]interface{}{42, true})
		low := c2.score([]interface{}{42, true})
		Expect(low).To(BeNumerically(">", 0))
		Expect(high).To(BeNumerically(">", low))
	})

	It("scores equivalence higher than other matches", func() {
		c := call{
			Params: []Matcher{BeEquivalentTo(42.0), BeEquivalentTo(true)},
		}
		c2 := call{
			Params: []Matcher{BeNumerically(">", 12.0), BeTrue()},
		}
		high := c.score([]interface{}{42, true})
		low := c2.score([]interface{}{42, true})
		Expect(low).To(BeNumerically(">", 0))
		Expect(high).To(BeNumerically(">", low))
	})
})
