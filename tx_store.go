package stomp

import (
	"container/list"
)

type txStore struct {
	transactions map[string]*list.List
}

// Initializes a new store or clears out an existing store
func (txs *txStore) Init() {
	txs.transactions = make(map[string]*list.List)
}

func (txs *txStore) Begin(tx string) error {
	if txs.transactions == nil {
		txs.transactions = make(map[string]*list.List)
	}

	if _, ok := txs.transactions[tx]; ok {
		return txAlreadyInProgress
	}

	txs.transactions[tx] = list.New()
	return nil
}

func (txs *txStore) Abort(tx string) error {
	if list, ok := txs.transactions[tx]; ok {
		list.Init()
		delete(txs.transactions, tx)
		return nil
	}
	return txUnknown
}

// Commit causes all requests that have been queued for the transaction
// to be sent to the request channel for processing. Calls the commit
// function (commitFunc) in order for each request that is part of the
// transaction.
func (txs *txStore) Commit(tx string, commitFunc func(r request)) error {
	if list, ok := txs.transactions[tx]; ok {
		for element := list.Front(); element != nil; element = list.Front() {
			commitFunc(list.Remove(element).(request))
		}
		delete(txs.transactions, tx)
		return nil
	}
	return txUnknown
}

func (txs *txStore) Add(tx string, request request) error {
	if list, ok := txs.transactions[tx]; ok {
		list.PushBack(request)
		return nil
	}
	return txUnknown
}
