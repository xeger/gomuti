package gomuti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
)

var _ = Describe("Allowed", func() {
	var Receiver Mock

	BeforeEach(func() {
		Receiver = Mock{}
	})

	Context("given a nonsensical call chain", func() {
		Receiver := Mock{}

		It("panics", func() {
			Expect(func() {
				Â(Receiver).With(true)
			}).To(Panic())
			Expect(func() {
				Â(Receiver).Return(false)
			}).To(Panic())
		})
	})

	Context("Call", func() {
		It("begins call chains", func() {
			Â(Receiver).Call("Foo").Return(true)
			Â(Receiver).Call("Bar").Panic(true)
			Â(Receiver).Call("Baz").With(42).Return(true)
			Â(Receiver).Call("Baz").With(Not(Equal(42))).Panic("not the answer")
		})
	})

	Context("With", func() {
		Context("given basic types", func() {
			PIt("matches equivalency")
		})
		Context("given struct types", func() {
			PIt("matches equivalency")
		})
		Context("given matchers", func() {
			PIt("tests satisfaction")
		})
	})

	Context("Return", func() {
		PIt("programs return values")
	})

	Context("Panic", func() {
		PIt("causes a panic")
	})
})
