package gomuti

import (
	"github.com/onsi/gomega/types"
	"github.com/xeger/gomuti/matchers"
)

// BeAnything matches _any_ value including nil and zero values.
func BeAnything() types.GomegaMatcher {
	return &matchers.BeAnythingMatcher{}
}

// Anything is an alias for BeAnything. It promotes readable mock call
// expectations. Example:
//
//    Allow(myMock).ToReceive("Foo").With(Anything())
//
// Which is equivalent to, but more readable, than:
//
//    Allow(myMock).ToReceive("Foo").With(BeAnything())
func Anything() types.GomegaMatcher {
	return BeAnything()
}

// HaveType matches any value whose type matches the specified name. Example:
//
//     Expect(4).To(HaveType("int"))
//
// Named types must be prefixed with the name of the package in which they
// are DEFINED (i.e. with the name that appears in the package statement of
// the source file where they are defined) and not with the import name that
// is used to REFER to them. Example:
//
//     import banana "time"
//		 Expect(banana.Now()).To(HaveType("time.Time"))
func HaveType(name string) types.GomegaMatcher {
	return &matchers.HaveTypeMatcher{Expected: name}
}

// AnythingOfType is an alias for HaveType. It promotes readable mock call
// expectations. Example:
//
//    Allow(myMock).ToReceive("Foo").With(AnythingOfType("mypkg.Widget"))
//
// Which is equivalent to, but more readable, than:
//
//    Allow(myMock).ToReceive("Foo").With(HaveType("mypkg.Widget"))
func AnythingOfType(name string) types.GomegaMatcher {
	return HaveType(name)
}

// HaveCall verifies that a method call was recored by a spy.
// You can add additional constraints to the verification by calling With() or
// Returning() on the matcher returned by this method.
func HaveCall(method string) *matchers.HaveCallMatcher {
	return &matchers.HaveCallMatcher{Method: method, Count: 1}
}

// HaveReceived is an alias for HaveCall().
func HaveReceived(method string) *matchers.HaveCallMatcher {
	return HaveCall(method)
}
