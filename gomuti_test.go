package gomuti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
	"github.com/xeger/gomuti/types"
)

type good struct {
	Mock types.Mock
}

type bad struct {
	mock types.Mock
}

var _ = Describe("Allow", func() {
	It("has an RSpec-like DSL", func() {
		double := good{Mock: types.Mock{}}
		Allow(double).ToReceive("Foo", 2)
	})

	It("has a gomega-like DSL", func() {
		double := good{Mock: types.Mock{}}
		Â(double).Call("Foo", 1)
	})

	It("accepts Mock", func() {
		double := types.Mock{}
		Â(double).Call("Foo", 1)
	})

	It("accepts struct with an initialized Mock field", func() {
		double := good{Mock: types.Mock{}}
		Â(double).Call("Foo", 2)
	})

	It("accepts pointer-to-struct with a nil Mock field", func() {
		double := &good{}
		Allow(double).ToReceive("Foo", 2)
		Expect(double.Mock).NotTo(BeNil())
	})

	It("panics when it cannot initialize a Mock", func() {
		Expect(func() {
			var double *types.Mock
			Â(double).Call("Foo", 1)
		}).To(Panic())

		Expect(func() {
			double := good{}
			Â(double).Call("Foo", 1)
			Expect(double).NotTo(BeNil())
		}).To(Panic())

		Expect(func() {
			var double types.Mock
			Â(double).Call("Foo", 1)
		}).To(Panic())
	})
})
