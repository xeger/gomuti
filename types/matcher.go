package types

// Matcher is a method parameter matcher for an Allowed call. It is a subset
// of Gomega's Matcher interface, but does not contain any error-reporting
// functions since a failure to match method parameters is not by itself an
// error.
type Matcher interface {
	Match(actual interface{}) (success bool, err error)
}
