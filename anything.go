package gomuti

import (
	"reflect"
	"strings"
)

type matchAnything struct{}

func (ma *matchAnything) Match(actual interface{}) (bool, error) {
	return true, nil
}

type matchType struct {
	name string
}

func (ma *matchType) Match(actual interface{}) (bool, error) {
	maDot := strings.Index(ma.name, ".")

	t := reflect.TypeOf(actual)
	acName := t.Name()
	acDot := strings.Index(acName, ".")

	if maDot > 0 {
		return (ma.name == acName), nil
	}

	maLocal := ma.name[maDot+1:]
	acLocal := acName
	if acDot > 0 {
		acLocal = acLocal[acDot+1:]
	}

	return (maLocal == acLocal), nil
}

// Anything matches any value.
func Anything() Matcher {
	return &matchAnything{}
}

// AnythingOfType matches any value whose type matches the specified name.
// If name contains a dot then the package and type name must match, else
// only the type name is matched.
func AnythingOfType(name string) Matcher {
	return &matchType{name: name}
}
