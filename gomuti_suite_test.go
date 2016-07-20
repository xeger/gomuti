package gomuti_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGomuti(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomuti Suite")
}
