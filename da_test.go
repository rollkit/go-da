package da_test

import (
	"github.com/rollkit/go-da/test"
	"testing"
)

func TestDummyDA(t *testing.T) {
	dummy := test.NewDummyDA()
	test.RunDATestSuite(t, dummy)
}
