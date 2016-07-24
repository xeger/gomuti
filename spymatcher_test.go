package gomuti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
)

// Type that contains a spy.
type infiltrated struct {
	Espion Spy
}

// Type that embeds a spy.
type infiltratedCovertly struct {
	Spy
}

var _ = Describe("SpyMatcher", func() {
	Context("Match", func() {
		var o infiltrated
		var po *infiltrated
		var e infiltratedCovertly
		var pe *infiltratedCovertly
		var sm *SpyMatcher

		BeforeEach(func() {
			o = infiltrated{Espion: Spy{}}
			po = &infiltrated{Espion: Spy{}}
			e = infiltratedCovertly{Spy{}}
			pe = &infiltratedCovertly{Spy{}}
			o.Espion.Observe("SecretMeeting")
			po.Espion.Observe("SecretMeeting")
			e.Observe("SecretMeeting")
			pe.Observe("SecretMeeting")
			sm = &SpyMatcher{Method: "SecretMeeting"}
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
