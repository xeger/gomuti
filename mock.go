package gomuti

type allowed struct {
	Params  []Matcher
	Panic   interface{}
	Results []interface{}
}

// Mock is a state container for mocked behavior.
type Mock map[string][]allowed
