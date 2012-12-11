package stomp

import (
	. "launchpad.net/gocheck"
	"github.com/jjeffery/stomp/message"
)

type TxStoreSuite struct{}

var _ = Suite(&TxStoreSuite{})

func (s *TxStoreSuite) TestDoubleBegin(c *C) {
	txs := txStore{}

	err := txs.Begin("tx1")
	c.Assert(err, IsNil)

	err = txs.Begin("tx1")
	c.Assert(err, Equals, txAlreadyInProgress)
}

func (s *TxStoreSuite) TestSuccessfulTx(c *C) {
	txs := txStore{}

	err := txs.Begin("tx1")
	c.Check(err, IsNil)

	err = txs.Begin("tx2")
	c.Assert(err, IsNil)

	f1 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/1")

	f2 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/2")

	f3 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/3")

	f4 := message.NewFrame(message.MESSAGE,
		message.Destination, "/queue/4")

	r1 := request{op: frameOp, frame: f1}
	r2 := request{op: frameOp, frame: f2}
	r3 := request{op: frameOp, frame: f3}
	r4 := request{op: frameOp, frame: f4}

	err = txs.Add("tx1", r1)
	c.Assert(err, IsNil)
	err = txs.Add("tx1", r2)
	c.Assert(err, IsNil)
	err = txs.Add("tx1", r3)
	c.Assert(err, IsNil)
	err = txs.Add("tx2", r4)

	var tx1Requests []request

	txs.Commit("tx1", func(r request) {
		tx1Requests = append(tx1Requests, r)
	})
	c.Check(err, IsNil)

	var tx2Requests []request

	err = txs.Commit("tx2", func(r request) {
		tx2Requests = append(tx2Requests, r)
	})
	c.Check(err, IsNil)
	
	c.Check(len(tx1Requests), Equals, 3)
	c.Check(tx1Requests[0].frame, Equals, f1)
	c.Check(tx1Requests[1].frame, Equals, f2)
	c.Check(tx1Requests[2].frame, Equals, f3)

	c.Check(len(tx2Requests), Equals, 1)
	c.Check(tx2Requests[0].frame, Equals, f4)
	
	// already committed, so should cause an error
	err = txs.Commit("tx1", func(r request) {
		c.Fatal("should not be called")
	})
	c.Check(err, Equals, txUnknown)
}
