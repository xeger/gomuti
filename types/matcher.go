package types

// Matcher is a method parameter matcher for mock calls. It is a subset
// of Gomega's Matcher interface, but does not contain any failure-message
// functions since a failure to match method parameters is not by itself a
// failure.
type Matcher interface {
	Match(actual interface{}) (success bool, err error)
}
