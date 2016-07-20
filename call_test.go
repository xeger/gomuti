package gomuti

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("call", func() {
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
