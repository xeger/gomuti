package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xeger/gomuti/matchers"
	"github.com/xeger/gomuti/types"
)

// Type that contains a spy.
type infiltrated struct {
	Espion types.Spy
}

// Type that embeds a spy.
type infiltratedCovertly struct {
	types.Spy
}

var _ = Describe("HaveCallMatcher", func() {
	Context("Match", func() {
		var o infiltrated
		var po *infiltrated
		var e infiltratedCovertly
		var pe *infiltratedCovertly
		var sm *matchers.HaveCallMatcher

		BeforeEach(func() {
			o = infiltrated{Espion: types.Spy{}}
			po = &infiltrated{Espion: types.Spy{}}
			e = infiltratedCovertly{types.Spy{}}
			pe = &infiltratedCovertly{types.Spy{}}
			o.Espion.Observe("SecretMeeting")
			po.Espion.Observe("SecretMeeting")
			e.Observe("SecretMeeting")
			pe.Observe("SecretMeeting")
			sm = &matchers.HaveCallMatcher{Method: "SecretMeeting"}
		})

		It("detects spies in the field", func() {
			m, err := sm.Match(o)
			Expect(m).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())

			m, err = sm.Match(po)
			Expect(m).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})

		It("detects embedded spies", func() {
			m, err := sm.Match(e)
			Expect(m).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())

			m, err = sm.Match(pe)
			Expect(m).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
