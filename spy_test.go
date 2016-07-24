package gomuti_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
)

// This exercises code in SpyMatcher in order to test functionality of Spy.
// Not strictly a unit test, but it's more idiomatic this way. If it causes
// trouble, just rewrite these tests to verify the behavior of Spy by hand.
var _ = Describe("Spy", func() {
	now := time.Now()
	isOne := Equal(1)
	isTrue := BeTrue()
	isFalse := BeFalse()
	puppies := BeEquivalentTo("puppies")
	shenanigans := BeEquivalentTo("shenanigans")
	answer := BeEquivalentTo(42.0)
	then := BeEquivalentTo(now)

	var s Spy
	BeforeEach(func() {
		s = Spy{}
	})

	It("counts method calls", func() {
		s.Observe("Foo", true)
		s.Observe("Foo", true)
		s.Observe("Foo", false)

		Expect(s).To(HaveCall("Foo").Times(3))
		Expect(s).To(HaveCall("Foo").With(isTrue).Times(2))
		Expect(s).To(HaveCall("Foo").With(isFalse).Times(1))
	})

	It("applies matchers to parameters", func() {
		s.Observe("Foo", 1, true, "puppies", 42.0, now)

		Expect(s).To(HaveCall("Foo").With(isOne, isTrue, puppies, answer, then).Once())
		Expect(s).To(HaveCall("Foo").With(isOne, isTrue, shenanigans, answer, then).Never())
	})

	It("fuzzy-matches parameters of basic types", func() {
		s.Observe("Bar", "shenanigans", true)
		s.Observe("Bar", "puppies", 7)

		Expect(s).To(HaveCall("Bar").Times(2))
		Expect(s).To(HaveCall("Bar").With("shenanigans", true).Times(1))
		Expect(s).To(HaveCall("Bar").With("hedgehogs", true).Never())
		Expect(s).To(HaveCall("Bar").With("puppies", 7).Times(1))
		Expect(s).To(HaveCall("Bar").With("puppies", 7.0).Times(1))
		Expect(s).To(HaveCall("Bar").With("puppies", 7.1).Never())
	})

	PIt("matches partial parameter lists")

	PIt("has a useful failure message")
})
