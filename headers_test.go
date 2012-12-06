package stomp

import (
	. "launchpad.net/gocheck"
)

type HeadersSuite struct{}

var _ = Suite(&HeadersSuite{})

func (s *HeadersSuite) TestAppend(c *C) {
	h := Headers{}
	h.Append("k1", "v1")
	h.Append("k2", "v2")
	h.Append("k1", "v1a")

	c.Check(len(h.headers), Equals, 6)
	v, ok := h.Contains("k1")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v1")
	v, ok = h.Contains("zz")
	c.Check(ok, Equals, false)
	c.Check(v, Equals, "")
}

func (s *HeadersSuite) TestContains(c *C) {
	h := Headers{}
	h.Append("k1", "v1")
	h.Append("k2", "v2")
	h.Append("k3", "v3")
	h.Append("k1", "obsolete")

	// tests that the first value is returned when
	// multiple keys exist
	v, ok := h.Contains("k1")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v1")

	// tests that values are not treated as keywords
	v, ok = h.Contains("v1")
	c.Check(ok, Equals, false)
	c.Check(v, Equals, "")

	v, ok = h.Contains("k2")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v2")

	v, ok = h.Contains("k3")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v3")

	v, ok = h.Contains("k4")
	c.Check(ok, Equals, false)
	c.Check(v, Equals, "")
}

func (s *HeadersSuite) TestSet(c *C) {
	h := Headers{}
	h.Append("k1", "xx")
	h.Append("k2", "v2")
	h.Set("k1", "v1")

	c.Check(len(h.headers), Equals, 4)
	v, ok := h.Contains("k1")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v1")
}

func (s *HeadersSuite) TestRemove(c *C) {
	h := Headers{}
	h.Append("k1", "v1")
	h.Append("k2", "v2")
	h.Remove("k1")

	c.Check(len(h.headers), Equals, 2)
	v, ok := h.Contains("k1")
	c.Check(ok, Equals, false)
	c.Check(v, Equals, "")
}

func (s *HeadersSuite) TestClone(c *C) {
	h1 := Headers{}
	h1.Append("k1", "v1")
	h1.Append("k2", "v2")

	h2 := h1.Clone()

	// after cloning, modify the original	
	h1.Set("k1", "xx")
	h1.Remove("k2")

	c.Check(len(h2.headers), Equals, 4)
	c.Check(h2.headers[0], Equals, "k1")
	c.Check(h2.headers[1], Equals, "v1")
	c.Check(h2.headers[2], Equals, "k2")
	c.Check(h2.headers[3], Equals, "v2")

	v, ok := h2.Contains("k1")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v1")
	v, ok = h2.Contains("k2")
	c.Check(ok, Equals, true)
	c.Check(v, Equals, "v2")

}
