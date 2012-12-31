package stomp

import (
	"launchpad.net/gocheck"
	"testing"
)

// Runs all gocheck tests in this package.
// See other *_test.go files for gocheck tests.
func TestStomp(t *testing.T) {
	gocheck.TestingT(t)
}

type StompSuite struct{}

var _ = gocheck.Suite(&StompSuite{})
