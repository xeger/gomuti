package gomuti

import (
	"github.com/onsi/gomega/types"
	"github.com/xeger/gomuti/matchers"
)

// BeAnything creates a matcher that is always satisfied. It is useful when
// you want to mock a method call but don't care about the parameter in a given
// position.
//
// Example:
//
//     Allow(boat).Call("Sail").With("west", BeAnything(), "mi").Panic("Please use kilometers")
func BeAnything() types.GomegaMatcher {
	return &matchers.BeAnythingMatcher{}
}

// Anything is an alias for BeAnything, designed to be more readable in the
// context of a mocked method call.
//
// Example:
//
//     Allow(boat).Call("Sail").With("west", Anything(), "mi").Panic("Please use kilometers")
//
// Which makes more sense than "with be-anything."
func Anything() types.GomegaMatcher {
	return BeAnything()
}

// HaveType creates a matcher that is satisfied by any value whose type matches
// the specified name. Example:
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

// AnythingOfType is an alias for HaveType, designed to be more readable in the
// context of a mocked method call.
//
// Example:
//
//     Allow(myMock).ToReceive("Foo").With(AnythingOfType("mypkg.Widget"))
//
// Which makes more sense than "with have-type."
func AnythingOfType(name string) types.GomegaMatcher {
	return HaveType(name)
}

// HaveCall is a spy method. It returns a matcher to verify that a method call
// was recorded by a spy. You can add more verifications (of parameter values,
// call count, etc) by calling methods on the returned matcher.
//
// Example:
//     Expect(double).To(HaveCall("Bar").With(true, 42).Twice())
func HaveCall(method string) *matchers.HaveCallMatcher {
	return &matchers.HaveCallMatcher{Method: method, Count: 1}
}

// HaveReceived is an alias for HaveCall().
func HaveReceived(method string) *matchers.HaveCallMatcher {
	return HaveCall(method)
}
