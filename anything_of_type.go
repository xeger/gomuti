package gomuti

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

// AnythingOfType matches any value whose type matches the specified name.
// If name contains a dot then the package and type name must match, else
// only the type name is matched.
func AnythingOfType(name string) Matcher {
	return &matchType{name: name}
}

type matchType struct {
	name string
}

func (ma *matchType) Match(actual interface{}) (bool, error) {
	maDot := strings.Index(ma.name, ".")

	t := reflect.TypeOf(actual)
	acTypename := t.Name()
	acPkgname := filepath.Base(t.PkgPath())
	acName := fmt.Sprintf("%s.%s", acPkgname, acTypename)
	acDot := strings.Index(acName, ".")

	if maDot > 0 {
		return (ma.name == acName), nil
	}

	// User wants to match against a type name (any package).
	maLocal := ma.name[maDot+1:]
	acLocal := acName
	if acDot > 0 {
		acLocal = acLocal[acDot+1:]
	}
	return (maLocal == acLocal), nil
}
