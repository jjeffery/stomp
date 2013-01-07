stomp
=====

Go language implementation of a STOMP client library.

Features:

* Supports STOMP Specifications Versions 1.0, 1.1, 1.2
* Protocol negotiation to select the latest mutually supported protocol
* Heart beating for testing the underlying network connection

For more information see http://godoc.org/github.com/jjeffery/stomp

Also contains a package for implementing a simple STOMP server.
Supports in-memory queues and topics, but can be easily extended to
support other queue storage mechanisms.

For more information on the STOMP server package see
http://godoc.org/github.com/jjeffery/stomp/server 

