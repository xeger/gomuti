package types_test

import (
	"net/url"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeger/gomuti"
	"github.com/xeger/gomuti/types"
)

var _ = Describe("Allowed", func() {
	var Receiver types.Mock

	BeforeEach(func() {
		Receiver = types.Mock{}
	})

	Context("given a nonsensical call chain", func() {
		Receiver := types.Mock{}

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
		It("accepts valid DSL", func() {
			Â(Receiver).Call("Foo").Return(true)
			Â(Receiver).Call("Bar").Panic(true)
			Â(Receiver).Call("Baz").With(42).Return(true)
			Â(Receiver).Call("Baz").With(Not(Equal(42))).Panic("not the answer")
			Â(Receiver).Call("Quux").With(42).Do(func(...interface{}) []interface{} {
				return nil
			})
		})
		It("panics with invalid DSL", func() {
			Expect(func() {
				Â(Receiver).Call("Foo").With(1).With(2)
			}).To(Panic())
			Expect(func() {
				a := Â(Receiver).Call("Foo")
				a.Return(true)
				a.Return(false)
			}).To(Panic())
			Expect(func() {
				a := Â(Receiver).Call("Foo")
				a.Panic("oh no")
				a.Return(false)
			}).To(Panic())
		})
	})

	Context("With", func() {
		Context("given basic types", func() {
			It("matches equivalency", func() {
				Â(Receiver).Call("Foo").With(1, 1.0, true, "stringy", 'X')

				r1 := Receiver.Call("Foo", 1, 1.0, true, "stringy", 'X')
				Expect(r1).NotTo(BeNil())

				r2 := Receiver.Call("Foo", 0, 0.0, false, "cheesy", 'Y')
				Expect(r2).To(BeNil())
			})
		})

		Context("given struct types", func() {
			It("matches equivalency", func() {
				u1, _ := url.Parse("https://github.com/foo")
				u2, _ := url.Parse("https://github.com/foo")
				u3, _ := url.Parse("https://github.com/bar")
				Â(Receiver).Call("Foo").With(*u1).Return(true)

				Expect(Receiver.Call("Foo", *u1)).NotTo(BeNil())
				Expect(Receiver.Call("Foo", *u2)).NotTo(BeNil())

				Expect(Receiver.Call("Foo", *u3)).To(BeNil())
			})
		})

		Context("given matchers", func() {
			It("tests satisfaction", func() {
				Â(Receiver).Call("Foo").With(BeNumerically(">", 0)).Return(true)
				Expect(Receiver.Call("Foo", 0)).To(BeNil())
				Expect(Receiver.Call("Foo", 1)).NotTo(BeNil())
			})
		})
	})

	Context("Return", func() {
		It("returns nothing when not called", func() {
			Â(Receiver).Call("Bar")
			Expect(Receiver.Call("Bar")).To(BeEmpty())
		})

		It("returns results when called", func() {
			Â(Receiver).Call("Foo").Return(1, 2, 3, 4)
			Expect(Receiver.Call("Foo")).To(BeEquivalentTo([]interface{}{1, 2, 3, 4}))
		})
	})

	Context("Panic", func() {
		It("causes a panic", func() {
			Â(Receiver).Call("Foo").Panic("howdy")
			Expect(func() {
				Receiver.Call("Foo")
			}).To(Panic())
		})
	})

	Context("Do", func() {
		It("accepts CallFunc", func() {
			Â(Receiver).Call("Foo").Do(func(params ...interface{}) []interface{} {
				r := len(params) >= 1 && reflect.ValueOf(params[0]).Bool()
				return []interface{}{r}
			})

			r := Receiver.Call("Foo", true, 42, "answer")
			Expect(r).NotTo(BeEmpty())
			Expect(r[0]).To(BeTrue())
			r = Receiver.Call("Foo", false, "question", nil)
			Expect(r).NotTo(BeEmpty())
			Expect(r[0]).NotTo(BeTrue())
		})

		It("accepts funcs with arbitrary signatures", func() {
			Â(Receiver).Call("Foo").Do(func(likable bool) bool {
				return likable
			})

			r := Receiver.Call("Foo", true)
			Expect(r).NotTo(BeEmpty())
			Expect(r[0]).To(BeTrue())
			r = Receiver.Call("Foo", false)
			Expect(r).NotTo(BeEmpty())
			Expect(r[0]).NotTo(BeTrue())
		})

		It("panics on parameter mismatches", func() {
			Â(Receiver).Call("Foo").Do(func(likable bool) bool {
				return likable
			})

			Expect(func() {
				Receiver.Call("Foo", 174)
			}).To(Panic())
		})
	})
})
