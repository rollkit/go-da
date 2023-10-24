package da_test

import (
	"testing"

	"github.com/rollkit/go-da/test"
)

func TestDummyDA(t *testing.T) {
	dummy := test.NewDummyDA()
	test.RunDATestSuite(t, dummy)
}
