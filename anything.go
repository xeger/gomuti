package gomuti

// Anything matches any value.
func Anything() Matcher {
	return &matchAnything{}
}

type matchAnything struct{}

func (ma *matchAnything) Match(actual interface{}) (bool, error) {
	return true, nil
}
