/*
Package stomp provides operations that allow communication with a message broker that supports the STOMP protocol. 
STOMP is the Streaming Text-Oriented Messaging Protocol. See http://stomp.github.com/ for more details.

This package provides support for all STOMP protocol features in 1.0, 1.1 and 1.2 including protocol negotiation,
heart-beating, and value encoding.

Connecting to a STOMP server is achieved using the stomp.Dial function, or the stomp.Connect function. See 
the examples section for a summary of how to use these functions. Both functions return a stomp.Conn object
for subsequent interaction with the STOMP server.

Once a connection (stomp.Conn) is created, it can be used to send messages to the STOMP server, or create
subscriptions for receiving messages from the STOMP server. Transactions can be created to send multiple
messages and/ or acknowledge multiple received messages from the server in one, atomic transaction. The 
examples section has examples of using subscriptions and transactions.

This package also exports types that represent a STOMP frame, and operate on STOMP frames. These types
include stomp.Frame, stomp.Reader, stomp.Writer and stomp.Validator. While a program can
use this package to communicate with a STOMP server without using these types directly, they could be
useful implementing a STOMP server in go.

The server subpackage provides a simple implementation of a STOMP server
that is useful for testing and could be useful for applications that require a simple,
STOMP-compliant message broker.
*/
package stomp
