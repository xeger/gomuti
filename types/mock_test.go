package types_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
	"github.com/xeger/gomuti/types"
)

var _ = Describe("Mock", func() {
	var m types.Mock
	BeforeEach(func() {
		m = types.Mock{}
		m.Allow().Call("Foo", 0, Anything()).Panic("I hate zero")
		m.Allow().Call("Foo", Anything(), 0).Panic("I also hate zero")
		m.Allow().Call("Foo", Anything(), Anything()).Return(false)
		m.Allow().Call("Foo", 42, 42).Return(42)
	})

	It("feigns panic", func() {
		Expect(func() {
			m.Call("Foo", 0)
		}).To(Panic())
	})

	Context("given ChooseCall is nil", func() {
		It("picks the most recently-allowed call", func() {
			r1 := m.Call("Foo", 1, 4)
			Expect(r1[0]).To(BeFalse())
			r2 := m.Call("Foo", 42, 42)
			Expect(r2[0]).To(Equal(42))
		})
	})
})
