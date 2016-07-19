package gomuti

// Matcher is a method parameter matcher for an Allowed call.
type Matcher interface {
	Match(actual interface{}) (success bool, err error)
}
