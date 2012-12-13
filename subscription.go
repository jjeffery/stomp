package stomp

import (
	"container/list"
)

// Represents a subscription.
type subscription struct {
	id          string
	destination string
	conn        *conn
	ack         bool // does the subscription acknowledge MESSAGE frames
	pending     bool // is there a pending frame acknowledgement
}

// Manages subscriptions
type subscriptionManager struct {
	destinations  map[string]*list.List              // maps queue name to a list of subscriptions
	subscriptions map[*conn]map[string]*subscription // maps connection and id to a subscription
}

// Create a subscription for the connection to the destination. The subscription
// is identified by id. If ack is true, all messages sent to this subscription
// will require acknowledgement.
func (s *subscriptionManager) Subscribe(c *conn, id, destination string, ack bool) error {
	sub := new(subscription)
	sub.id = id
	sub.destination = destination
	sub.conn = c
	sub.ack = ack

	if s.subscriptions == nil {
		s.subscriptions = make(map[*conn]map[string]*subscription)
	}

	connMap, ok := s.subscriptions[c]
	if !ok {
		connMap = make(map[string]*subscription)
		s.subscriptions[c] = connMap
	}

	if _, ok = connMap[id]; ok {
		return subscriptionInUse
	}

	connMap[id] = sub

	if s.destinations == nil {
		s.destinations = make(map[string]*list.List)
	}

	subList, ok := s.destinations[destination]
	if !ok {
		subList = list.New()
		s.destinations[destination] = subList
	}

	subList.PushFront(sub)
	return nil
}

// Unsubscribe an existing subscription.
func (s *subscriptionManager) Unsubscribe(c *conn, id string) error {
	sub := s.subscriptions[c][id]
	if sub == nil {
		return subscriptionNotFound
	}

	subList := s.destinations[sub.destination]
	for e := subList.Front(); e != nil; e = e.Next() {
		otherSub := e.Value.(*subscription)
		if otherSub == sub {
			subList.Remove(e)
			break
		}
	}

	return nil
}

// Find a subscription matching the destination. If multiple subscriptions
// are available, choose one on a round-robin basis.
func (s *subscriptionManager) Find(destination string) *subscription {
	subList := s.destinations[destination]
	if subList != nil {
		for element := subList.Front(); element != nil; element = element.Next() {
			sub := element.Value.(*subscription)
			if !sub.pending {
				// Move to back of list so that if there are multiple
				// subscriptions for the same destination, they round-robin.
				subList.MoveToBack(element)
				return sub
			}
		}
	}

	return nil
}
