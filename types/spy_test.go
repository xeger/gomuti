package types_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
	"github.com/xeger/gomuti/types"
)

// This exercises code in HaveCallMatcher in order to test functionality of Spy.
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

	var s types.Spy
	BeforeEach(func() {
		s = types.Spy{}
	})

	It("complains when nil is called", func() {
		s = nil

		checker := func() {
			r := recover()

			str, _ := r.(string)
			Expect(str).NotTo(BeNil())

			Expect(str).To(MatchRegexp("^gomuti:"))
			panic(r)
		}

		Expect(func() {
			defer checker()
			s.Observe("Foo")
		}).To(Panic())

		Expect(func() {
			defer checker()
			s.Count("Foo")
		}).To(Panic())
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

	It("equality-matches parameters of basic types", func() {
		s.Observe("Bar", "shenanigans", true)
		s.Observe("Bar", "puppies", 7)
		s.Observe("Bar", "puppies", 7.0)

		Expect(s).To(HaveCall("Bar").Times(2))
		Expect(s).To(HaveCall("Bar").With("shenanigans", true).Times(1))
		Expect(s).To(HaveCall("Bar").With("hedgehogs", true).Never())
		Expect(s).To(HaveCall("Bar").With("puppies", 7).Times(1))
		Expect(s).To(HaveCall("Bar").With("puppies", 7.0).Times(1))
		Expect(s).To(HaveCall("Bar").With("puppies", 7.1).Never())
	})

	It("equivalence-matches parameters of non-basic types", func() {
		t0 := time.Now()
		t1 := t0
		t2 := t0

		s.Observe("Bar", &t1)
		s.Observe("Bar", &t2)
		Expect(s).To(HaveCall("Bar").Times(2))
		Expect(s).To(HaveCall("Bar").With(&t0).Times(2))
	})

	PIt("matches partial parameter lists")

	PIt("has a useful failure message")
})
