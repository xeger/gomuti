package types_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"

	"github.com/xeger/gomuti/types"
)

var _ = Describe("MatchParams", func() {
	It("matches basic-type values using equality", func() {
		out := types.MatchParams([]interface{}{1, true, 2.0, "hi"})
		for _, m := range out {
			Expect(m).To(HaveType("*matchers.EqualMatcher"))
		}
	})

	It("matches non-basic-type values using equivalence", func() {
		out := types.MatchParams([]interface{}{
			time.Now(),
			[]int{1, 2, 3},
			map[string]bool{"Joe": true},
		})
		for _, m := range out {
			Expect(m).To(HaveType("*matchers.BeEquivalentToMatcher"))
		}
	})

	It("matches nil values using BeNil", func() {
		out := types.MatchParams([]interface{}{
			nil,
		})
		for _, m := range out {
			Expect(m).To(HaveType("*matchers.BeNilMatcher"))
		}
	})
})
