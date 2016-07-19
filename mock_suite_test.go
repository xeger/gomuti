package gomuti_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
)

type good struct {
	Mock Mock
}

type bad struct {
	mock Mock
}

var _ = Describe("Allow", func() {
	It("has an RSpec-like DSL", func() {
		double := good{Mock: Mock{}}
		Allow(double).ToReceive("Foo", 2)
	})

	It("has a gomega-like DSL", func() {
		double := good{Mock: Mock{}}
		Â(double).Call("Foo", 1)
	})

	It("accepts Mock", func() {
		double := Mock{}
		Â(double).Call("Foo", 1)
	})

	It("accepts struct with a Mock field", func() {
		double := good{Mock: Mock{}}
		Â(double).Call("Foo", 2)
	})

	It("accepts pointer-to-struct with a Mock field", func() {
		double := &good{}
		Allow(double).ToReceive("Foo", 2)
		Expect(double.Mock).NotTo(BeNil())
	})

	It("panics with pointer-to-Mock", func() {
		Expect(func() {
			var double *Mock
			Â(double).Call("Foo", 1)
			Expect(double).NotTo(BeNil())
		}).To(Panic())
	})
})

func TestMock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock Suite")
}
