stomp
=====

Go language implementation of a STOMP client library.

[![Build Status](https://travis-ci.org/go-stomp/stomp.svg?branch=master)](https://travis-ci.org/go-stomp/stomp)

Features:

* Supports STOMP Specifications Versions 1.0, 1.1, 1.2 (https://stomp.github.io/)
* Protocol negotiation to select the latest mutually supported protocol
* Heart beating for testing the underlying network connection
* Tested against RabbitMQ v3.0.1

For API documentation, see http://godoc.org/github.com/go-stomp/stomp

Usage Instructions
==================

```
go get github.com/go-stomp/stomp
```

Previous Versions
=================

An earlier version of this package made use of Gustavo Niemeyer's gopkg.in facility
for versioning Go libraries. This earlier version of the library is still available

```
go get gopkg.in/stomp.v1
```

API documentation for this earlier version can be found at http://gopkg.in/stomp.v1




