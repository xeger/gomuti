# Gomuti

[![Build Status](https://travis-ci.org/xeger/gomuti.png)](https://travis-ci.org/xeger/gomuti)

Gomuti is DSL for [mocking](https://en.wikipedia.org/wiki/Mock_object) Golang interfaces, inspired by  [Gomega](https://github.com/onsi/gomega) and drawing upon Gomega matchers to dispatch mock method calls.
With a [Ginkgo](https://github.com/onsi/ginkgo)-like DSL for programming mock behavior, Gomuti makes it
easy to write beautiful, well-isolated unit tests.

Mocks can also be [spies](https://robots.thoughtbot.com/spy-vs-spy#what-do-you-mean-spy) and [stubs](),
enabling [behavior-driven development](https://en.wikipedia.org/wiki/Behavior-driven_development) and
terse, easy-to-maintain mock setup.

## How to use

Imagine you have an interface that you want to mock.

```go
type Adder interface {
  Add(l, r int64) int64
}
```

To properly mock this interface, you need to create a struct type that has the
same methods. The struct also holds two Gomuti-related fields that keep state
about the programmed behavior of the mock and the observed method calls:

```go
  import gtypes "github.com/xeger/gomuti/types"

  type MockAdder struct {
    Mock gtypes.Mock
    Spy gtypes.Spy
  }

  func(m *MockAdder) Add(l, r int64) int64 {
    m.Spy.Observe(l, r)
    r := m.Mock.Call("Add", l, r)
    return r[0].(int64)
  }
```

In reality you would use [Mongoose](https://github.com/xeger/mongoose) to
generate a mock type and methods for every interface in your package,
but a hand-coded mock is fine for example purposes.

To program behavior into your mock, use the DSL methods in the `gomuti`
package. `Allow()` instructs your mock to expect a method call and
tells it what to return. After calling the code under test, you can
use the `HaveReceived()` Gomega matcher to verify the number of calls
actually made to your mock as well as the specific parameter values.

Imagine you're writing unit tests for the `Multiplier` type
and you want to isolate yourself from bugs in `Adder`.

```go
  import (
    . "github.com/onsi/ginkgo"
    . "github.com/xeger/gomuti"
  )

  Describe("multiplier", func() {
    var subject *multiplier
    var adder Adder
    BeforeEach(func() {
      adder = &MockAdder{}
      m := &multiplier{Adder:a}
    })

    It("computes the product of two integers", func() {
      Allow(adder).ToReceive("Add").With(5,5).AndReturn(10)
      Allow(adder).ToReceive("Add").With(10,5).AndReturn(15)
      result := subject.Multiply(3,5))
      Expect(adder).To(HaveReceived("Add").Times(2))
      Expect(result).To(Equal(15))
    })
  })
```

The `Allow()` DSL can use any Gomega matcher for method parameters and Gomuti
provides a few matchers of its own; together, these allow you to mock
sophisticated behavior. Imagine your adder has a new `AddStuff()` feature
that adds arbitrarily-typed values.

```go
  Allow(adder).ToReceive("AddStuff").With(AnythingOfType("bool"), Anything()).AndReturn(True)
```

### Terse DSL

Gomuti's long-form DSL is inspired by RSpec and its "plain English" approach.
There is also a short-form DSL built around the method `gomuti.Â()`. To produce
the Â character, type `Alt+0194` on Windows keyboards or `Shift+Option+M` on Mac keyboards
(as a mnemonic, think "**Â** allows my **M**ock the **o**ption of being called.") 

Short-form equivalents are provided for `ToReceive()` and other chained methods, and
a super-terse form of Â delivers maximum brevity. If we also use Gomega's Ω method, our
tests get very terse indeed. (Some would say "unreadable," but beauty is in the eye of
the beholder.)

```go
  // Super terse DSL
  Â(adder,"Add",5,5).Return(10)

  // Moderately terse DSL with complex matcher.
  big := BeNumerically(">",2**32-1)
  Â(adder).Call("Add").With(big,Anything()).Panic("integer overflow")

  Ω(subject.Multiply(2,5)).Should(Equal(10))
  Ω(subject).Should(HaveCall("Add").Times(2))

  Ω(func() {
    subject.Multiply(2**32-1,1)
  }).Should(Panic())
```

Long and short method calls are interchangeable; even when using the
long-form `Allow()`, we recommended using `Call()` instead of `ToReceive()`
because the word "receive" is usually associated with the channel-receive
operation.

## How to get help

Check the [frequently-asked questions](FAQ.md) to see if your problem is common.

If you think Gomuti is missing a feature, check the [roadmap](TODO.md) to see
if a similar feature is already planned.

If you still need help, [open an Issue](https://github.com/xeger/mongoose/issues/new).
Clearly explain your problem, steps to reproduce, and your ideal solution (if known).

## How to contribute

Fork the `xeger/gomuti` [repository](https://github.com/xeger/gomuti) on GitHub; make your changes; open a pull request.
