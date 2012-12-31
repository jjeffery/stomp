package stomp

import (
	. "launchpad.net/gocheck"
)

func (s *StompSuite) TestHeaderGetSetAddDel(c *C) {
	h := Header{}
	c.Assert(h.Get("xxx"), Equals, "")
	h.Add("xxx", "yyy")
	c.Assert(h.Get("xxx"), Equals, "yyy")
	h.Add("xxx", "zzz")
	c.Assert(h["xxx"], DeepEquals, []string{"yyy", "zzz"})
	h.Set("xxx", "111")
	c.Assert(h.Get("xxx"), Equals, "111")
	h.Del("xxx")
	c.Assert(h.Get("xxx"), Equals, "")
}

func (s *StompSuite) TestHeaderClone(c *C) {
	h := Header{}
	h.Set("xxx", "yyy")
	h.Set("yyy", "zzz")

	hc := h.Clone()
	h.Del("xxx")
	h.Del("yyy")
	c.Assert(hc.Get("xxx"), Equals, "yyy")
	c.Assert(hc.Get("yyy"), Equals, "zzz")
}

func (s *StompSuite) TestLit(c *C) {
	_ = Frame{
		Command: "CONNECT",
		Header: Header{
			"login":    []string{"xxx"},
			"passcode": []string{"yyy"},
		},
		Body: []byte{1, 2, 3, 4},
	}
}
