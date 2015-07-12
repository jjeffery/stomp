# Breaking Changes

This document provides a list of breaking changes since the V1 release
of the stomp client library.

## 1. No longer using gopkg.in

Version 1 of the library used Gustavo Niemeyer's `gopkg.in` facility for versioning Go libraries.
For a number of reasons, the `stomp` library no longer uses this facility. For this reason the
import path has changed.

Version 1:
```go
import (
    "gopkg.in/stomp.v1"
)
```

Version 2:
```go
import (
    "github.com/go-stomp/stomp"
)
```

## 2. Frame types moved to frame package

Version 1 of the library included a number of types to do with STOMP frames in the `stomp`
package, and the `frame` package consisted of just a few constant definitions.

It was decided to move the following types out of the `stomp` package and into the `frame` package:

* `stomp.Frame` -> `frame.Frame`
* `stomp.Header` -> `frame.Header`
* `stomp.Reader` -> `frame.Reader`
* `stomp.Writer` -> `frame.Writer`

This change was considered worthwhile for the following reasons:

* This change reduces the surface area of the `stomp` package and makes it easier to learn.
* Ideally, users of the `stomp` package do not need to directly reference the items in the `frame`
package, and the types moved are not needed in normal usage of the `stomp` package.



